package fuzz

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
)

func FuzzExpr(f *testing.F) {
	b, err := os.ReadFile("./fuzz_corpus.txt")
	if err != nil {
		panic(err)
	}
	corpus := strings.Split(strings.TrimSpace(string(b)), "\n")
	for _, s := range corpus {
		f.Add(s)
	}

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
