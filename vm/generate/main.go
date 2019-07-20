package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var data string
	echo := func(s string, xs ...interface{}) {
		data += fmt.Sprintf(s, xs...) + "\n"
	}

	echo(`package vm`)
	echo(`import (`)
	echo(`"fmt"`)
	echo(`"reflect"`)
	echo(`)`)

	types := []string{
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
	}

	helpers := []struct {
		name, op        string
		noFloat, string bool
	}{
		{
			name:   "equal",
			op:     "==",
			string: true,
		},
		{
			name:   "less",
			op:     "<",
			string: true,
		},
		{
			name:   "more",
			op:     ">",
			string: true,
		},
		{
			name:   "lessOrEqual",
			op:     "<=",
			string: true,
		},
		{
			name:   "moreOrEqual",
			op:     ">=",
			string: true,
		},
		{
			name:   "add",
			op:     "+",
			string: true,
		},
		{
			name: "subtract",
			op:   "-",
		},
		{
			name: "multiply",
			op:   "*",
		},
		{
			name: "divide",
			op:   "/",
		},
		{
			name:    "modulo",
			op:      "%",
			noFloat: true,
		},
	}

	for _, helper := range helpers {
		name := helper.name
		op := helper.op
		echo(`func %v(a, b interface{}) interface{} {`, name)
		echo(`switch x := a.(type) {`)
		for i, a := range types {
			if helper.noFloat && strings.HasPrefix(a, "float") {
				continue
			}
			echo(`case %v:`, a)
			echo(`switch y := b.(type) {`)
			for j, b := range types {
				if helper.noFloat && strings.HasPrefix(b, "float") {
					continue
				}
				echo(`case %v:`, b)
				if i == j {
					echo(`return x %v y`, op)
				}
				if i < j {
					echo(`return %v(x) %v y`, b, op)
				}
				if i > j {
					echo(`return x %v %v(y)`, op, a)
				}
			}
			echo(`}`)
		}
		if helper.string {
			echo(`case string:`)
			echo(`switch y := b.(type) {`)
			echo(`case %v:`, "string")
			echo(`return x %v y`, op)
			echo(`}`)
		}
		echo(`}`)
		if name == "equal" {
			echo(`// Two nil values should be considered as equal.`)
			echo(`if isNil(a) && isNil(b) { return true }`)
			echo(`return reflect.DeepEqual(a, b)`)
		} else {
			echo(`panic(fmt.Sprintf("invalid operation: %%T %%v %%T", a, "%v", b))`, op)
		}
		echo(`}`)
		echo(``)
	}

	// isNil func
	echo(`func isNil(v interface{}) bool {`)
	echo(`if v == nil {`)
	echo(`return true`)
	echo(`}`)
	echo(`r := reflect.ValueOf(v)`)
	echo(`switch r.Kind() {`)
	echo(`case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:`)
	echo(`return r.IsNil()`)
	echo(`default:`)
	echo(`return false`)
	echo(`}`)
	echo(`}`)

	b, err := format.Source([]byte(data))
	check(err)
	err = ioutil.WriteFile("helpers.go", b, 0644)
	check(err)
}
