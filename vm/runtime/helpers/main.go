package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"
)

func main() {
	var b bytes.Buffer
	err := template.Must(
		template.New("helpers").
			Funcs(template.FuncMap{
				"cases":          func(op string) string { return cases(op, uints, ints, floats) },
				"cases_int_only": func(op string) string { return cases(op, uints, ints) },
				"cases_with_duration": func(op string) string {
					return cases(op, uints, ints, floats, []string{"time.Duration"})
				},
				"array_equal_cases": func() string { return arrayEqualCases([]string{"string"}, uints, ints, floats) },
			}).
			Parse(helpers),
	).Execute(&b, nil)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Print(string(formatted))
}

var ints = []string{
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
}

var uints = []string{
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
}

var floats = []string{
	"float32",
	"float64",
}

func cases(op string, xs ...[]string) string {
	var types []string
	for _, x := range xs {
		types = append(types, x...)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Generating %s cases for %v\n", op, types)

	var out string
	echo := func(s string, xs ...any) {
		out += fmt.Sprintf(s, xs...) + "\n"
	}
	for _, a := range types {
		echo(`case %v:`, a)
		echo(`switch y := b.(type) {`)
		for _, b := range types {
			echo(`case %v:`, b)
			if op == "/" {
				echo(`return float64(x) / float64(y)`)
			} else {
				t := castType(a, b, op)
				if t == "safe_uint64" {
					// Special cross-sign comparison for uint64 vs signed
					// Use direct append to avoid fmt.Sprintf interpreting % as format verb
					out += safeUint64Op(a, b, op) + "\n"
				} else {
					echo(`return %v(x) %v %v(y)`, t, op, t)
				}
			}
		}
		echo(`}`)
	}
	return strings.TrimRight(out, "\n")
}

func castType(a, b, op string) string {
	// Float takes priority over duration (matching original behavior):
	// duration * float → float64(x) * float64(y) → float64 result
	if isFloat(a) || isFloat(b) {
		return "float64"
	}
	if isDuration(a) || isDuration(b) {
		return "time.Duration"
	}
	// For uint64 mixed with signed integers, we need safe comparison
	if isUint64(a) && isSigned(b) {
		return "safe_uint64"
	}
	if isSigned(a) && isUint64(b) {
		return "safe_uint64"
	}
	// For uint64 vs uint64 or uint64 vs other unsigned types, use uint64
	if isUint64(a) || isUint64(b) {
		return "uint64"
	}
	return "int"
}

func isUint64(t string) bool {
	return t == "uint64"
}

func isSigned(t string) bool {
	return strings.HasPrefix(t, "int")
}

func safeUint64Op(a, b, op string) string {
	// a is the x type (outer switch), b is the y type (inner switch)
	if isUint64(a) && isSigned(b) {
		// x is uint64, y is signed
		switch op {
		case "==":
			return "if y < 0 { return false }\nreturn x == uint64(y)"
		case "<":
			return "if y < 0 { return false }\nreturn x < uint64(y)"
		case ">":
			return "if y < 0 { return true }\nreturn x > uint64(y)"
		case "<=":
			return "if y < 0 { return false }\nreturn x <= uint64(y)"
		case ">=":
			return "if y < 0 { return true }\nreturn x >= uint64(y)"
		case "+", "-", "*":
			return fmt.Sprintf("return uint64(x) %s uint64(y)", op)
		case "%":
			return fmt.Sprintf("return uint64(x) %s uint64(y)", op)
		}
	}
	if isSigned(a) && isUint64(b) {
		// x is signed, y is uint64
		switch op {
		case "==":
			return "if x < 0 { return false }\nreturn uint64(x) == y"
		case "<":
			return "if x < 0 { return true }\nreturn uint64(x) < y"
		case ">":
			return "if x < 0 { return false }\nreturn uint64(x) > y"
		case "<=":
			return "if x < 0 { return true }\nreturn uint64(x) <= y"
		case ">=":
			return "if x < 0 { return false }\nreturn uint64(x) >= y"
		case "+", "-", "*":
			return fmt.Sprintf("return uint64(x) %s uint64(y)", op)
		case "%":
			return fmt.Sprintf("return uint64(x) %s uint64(y)", op)
		}
	}
	return fmt.Sprintf("return int(x) %s int(y)", op)
}

func arrayEqualCases(xs ...[]string) string {
	var types []string
	for _, x := range xs {
		types = append(types, x...)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Generating array equal cases for %v\n", types)

	var out string
	echo := func(s string, xs ...any) {
		out += fmt.Sprintf(s, xs...) + "\n"
	}
	echo(`case []any:`)
	echo(`switch y := b.(type) {`)
	for _, a := range append(types, "any") {
		echo(`case []%v:`, a)
		echo(`if len(x) != len(y) { return false }`)
		echo(`for i := range x {`)
		echo(`if !Equal(x[i], y[i]) { return false }`)
		echo(`}`)
		echo("return true")
	}
	echo(`}`)
	for _, a := range types {
		echo(`case []%v:`, a)
		echo(`switch y := b.(type) {`)
		echo(`case []any:`)
		echo(`return Equal(y, x)`)
		echo(`case []%v:`, a)
		echo(`if len(x) != len(y) { return false }`)
		echo(`for i := range x {`)
		echo(`if x[i] != y[i] { return false }`)
		echo(`}`)
		echo("return true")
		echo(`}`)
	}
	return strings.TrimRight(out, "\n")
}

func isFloat(t string) bool {
	return strings.HasPrefix(t, "float")
}

func isDuration(t string) bool {
	return t == "time.Duration"
}

const helpers = `// Code generated by vm/runtime/helpers/main.go. DO NOT EDIT.

package runtime

import (
	"fmt"
	"reflect"
	"time"
)

func Equal(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "==" }}
	{{ array_equal_cases }}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x == y
		}
	case bool:
		switch y := b.(type) {
		case bool:
			return x == y
		}
	}
	if IsNil(a) && IsNil(b) {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func Less(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "<" }}
	case string:
		switch y := b.(type) {
		case string:
			return x < y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x < y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T < %T", a, b))
}

func More(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases ">" }}
	case string:
		switch y := b.(type) {
		case string:
			return x > y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x > y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T > %T", a, b))
}

func LessOrEqual(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "<=" }}
	case string:
		switch y := b.(type) {
		case string:
			return x <= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x <= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T <= %T", a, b))
}

func MoreOrEqual(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases ">=" }}
	case string:
		switch y := b.(type) {
		case string:
			return x >= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x >= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T >= %T", a, b))
}

func Add(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases "+" }}
	case string:
		switch y := b.(type) {
		case string:
			return x + y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Duration:
			return x.Add(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Time:
			return y.Add(x)
		case time.Duration:
			return x + y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T + %T", a, b))
}

func Subtract(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases "-" }}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Sub(y)
		case time.Duration:
			return x.Add(-y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x - y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T - %T", a, b))
}

func Multiply(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases_with_duration "*" }}
	}
	panic(fmt.Sprintf("invalid operation: %T * %T", a, b))
}

func Divide(a, b interface{}) float64 {
	switch x := a.(type) {
	{{ cases "/" }}
	}
	panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
}

func Modulo(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases_int_only "%" }}
	}
	panic(fmt.Sprintf("invalid operation: %T %% %T", a, b))
}
`
