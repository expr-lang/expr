package fuzz

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
)

//go:embed fuzz_corpus.txt
var fuzzCorpus string

func FuzzExpr(f *testing.F) {
	corpus := strings.Split(strings.TrimSpace(fuzzCorpus), "\n")
	for _, s := range corpus {
		f.Add(s)
	}

	var env = map[string]any{
		"ok":    true,
		"f64":   .5,
		"f32":   float32(.5),
		"i":     1,
		"i64":   int64(1),
		"i32":   int32(1),
		"array": []int{1, 2, 3, 4, 5},
		"list":  []Foo{{"bar"}, {"baz"}},
		"foo":   Foo{"bar"},
		"add":   func(a, b int) int { return a + b },
		"div":   func(a, b int) int { return a / b },
		"half":  func(a float64) float64 { return a / 2 },
		"score": func(a int, x ...int) int {
			s := a
			for _, n := range x {
				s += n
			}
			return s
		},
		"greet": func(name string) string { return "Hello, " + name },
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
		regexp.MustCompile(`invalid date .*`),
		regexp.MustCompile(`cannot parse .* as .*`),
	}

	f.Fuzz(func(t *testing.T, code string) {
		program, err := expr.Compile(code, expr.Env(env))
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
			t.Errorf("%s", err)
		}
	})
}

type Foo struct {
	Bar string
}

func (f Foo) String() string {
	return "foo"
}

func (f Foo) Qux(s string) string {
	return f.Bar + s
}
