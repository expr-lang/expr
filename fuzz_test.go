package expr_test

import (
	"regexp"
	"testing"

	"github.com/antonmedv/expr"
)

func FuzzExpr(f *testing.F) {
	env := map[string]interface{}{
		"i":   1,
		"j":   2,
		"b":   true,
		"a":   []int{1, 2, 3},
		"m":   map[string]interface{}{"a": 1, "b": 2, "m": map[string]int{"a": 1}},
		"s":   "abc",
		"add": func(a, b int) int { return a + b },
		"foo": Foo{A: 1, B: 2, Bar: Bar{A: 1, B: 2}},
	}

	head := expr.Function(
		"head",
		func(params ...interface{}) (interface{}, error) {
			return params[0], nil
		},
		new(func(int) int),
	)

	corpus := []string{
		`.5 + .5`,
		`i + j`,
		`i - j`,
		`i * j`,
		`i / j`,
		`i % j`,
		`true || false`,
		`true && false`,
		`i == j`,
		`i != j`,
		`i > j`,
		`i >= j`,
		`i < j`,
		`i <= j`,
		`i in a`,
		`i not in a`,
		`s in m`,
		`m.a`,
		`m.m.a`,
		`a[0]`,
		`a[i]`,
		`a[i:j]`,
		`a[i:]`,
		`a[:j]`,
		`a[:]`,
		`a[1:-1]`,
		`len(a)`,
		`type(a)`,
		`abs(-1)`,
		`int(0.5)`,
		`float(42)`,
		`string(i)`,
		`trim(" a ")`,
		`trim("_a_", "_")`,
		`trimPrefix("  a", " ")`,
		`trimSuffix("a  ")`,
		`upper("a")`,
		`lower("A")`,
		`split("a,b,c", ",")`,
		`replace("a,b,c", ",", "_")`,
		`repeat("a", 3)`,
		`join(["a", "b", "c"], ",")`,
		`indexOf("abc", "b")`,
		`max(1,2,3)`,
		`min(1,2,3)`,
		`toJSON(a)`,
		`fromJSON("[1,2,3]")`,
		`now()`,
		`duration("1s")`,
		`first(a)`,
		`last(a)`,
		`get(m, "a")`,
		`1..9 | filter(i > 5) | map(i * 2)`,
		`s startsWith "a"`,
		`s endsWith "c"`,
		`s contains "a"`,
		`s matches "a"`,
		`s matches "a+"`,
		`true ? 1 : 2`,
		`false ? 1 : 2`,
		`b ?? true`,
		`head(1)`,
		`{a: 1, b: 2}`,
		`[1, 2, 3]`,
		`type(1)`,
		`type("a")`,
		`type([1, 2, 3])`,
		`type({a: 1, b: 2})`,
		`type(head)`,
		`keys(m)`,
		`values(m)`,
		`foo.bar.a`,
		`foo.bar.b`,
	}

	for _, s := range corpus {
		f.Add(s)
	}

	skip := []*regexp.Regexp{
		regexp.MustCompile(`cannot fetch .* from .*`),
		regexp.MustCompile(`cannot get .* from .*`),
		regexp.MustCompile(`cannot slice`),
		regexp.MustCompile(`slice index out of range`),
		regexp.MustCompile(`error parsing regexp`),
		regexp.MustCompile(`integer divide by zero`),
		regexp.MustCompile(`interface conversion`),
		regexp.MustCompile(`invalid argument for .*`),
		regexp.MustCompile(`invalid character`),
		regexp.MustCompile(`invalid operation`),
		regexp.MustCompile(`invalid duration`),
		regexp.MustCompile(`time: missing unit in duration`),
		regexp.MustCompile(`time: unknown unit .* in duration`),
		regexp.MustCompile(`json: unsupported value`),
		regexp.MustCompile(`unexpected end of JSON input`),
		regexp.MustCompile(`memory budget exceeded`),
		regexp.MustCompile(`using interface \{} as type .*`),
		regexp.MustCompile(`reflect.Value.MapIndex: value of type .* is not assignable to type .*`),
		regexp.MustCompile(`reflect: Call using .* as type .*`),
		regexp.MustCompile(`reflect: call of reflect.Value.Call on .* Value`),
		regexp.MustCompile(`reflect: call of reflect.Value.Index on map Value`),
		regexp.MustCompile(`reflect: call of reflect.Value.Len on .* Value`),
		regexp.MustCompile(`strings: negative Repeat count`),
		regexp.MustCompile(`strings: illegal bytes to escape`),
		regexp.MustCompile(`operator "in" not defined on int`),
	}

	f.Fuzz(func(t *testing.T, code string) {
		program, err := expr.Compile(code, expr.Env(env), head, expr.ExperimentalPipes())
		if err != nil {
			t.Skipf("compile error: %s", err)
		}

		_, err = expr.Run(program, env)
		if err != nil {
			for _, r := range skip {
				if r.MatchString(err.Error()) {
					t.Skipf("skip error: %s", err)
					return
				}
			}
			t.Errorf("code: %s\nerr: %s", code, err)
		}
	})
}

type Foo struct {
	A   int `expr:"a"`
	B   int `expr:"b"`
	Bar Bar `expr:"bar"`
}

type Bar struct {
	A int `expr:"a"`
	B int `expr:"b"`
}
