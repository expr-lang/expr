package expr_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/test/playground"
)

func FuzzExpr(f *testing.F) {
	env := playground.ExampleData()

	b, err := os.ReadFile("testdata/corpus.txt")
	if err != nil {
		f.Fatal(err)
	}
	for _, testcase := range strings.Split(string(b), "\n") {
		f.Add(testcase)
	}

	f.Add(`Posts[0].Comments[0].AuthorEmail()`)
	f.Add(`Posts[0].Comments[0].Upvoted()`)
	f.Add(`Posts[0].Comments[0].CommentDate.Add(now())`)
	f.Add(`Posts[0].Comments[0].CommentDate.AddDate(0, 0, 1)`)
	f.Add(`Authors[Posts[0].Author.ID].Profile.Biography`)
	f.Add(`Authors[2].Profile.Age()`)
	f.Add(`Authors[2].Profile.Website`)

	okCases := []*regexp.Regexp{
		regexp.MustCompile(`cannot slice`),
		regexp.MustCompile(`integer divide by zero`),
		regexp.MustCompile(`invalid operation`),
		regexp.MustCompile(`interface conversion`),
		regexp.MustCompile(`memory budget exceeded`),
		regexp.MustCompile(`slice index out of range`),
		regexp.MustCompile(`cannot fetch .* from .*`),
		regexp.MustCompile(`cannot get .* from .*`),
		regexp.MustCompile(`invalid argument for .*`),
		regexp.MustCompile(`json: unsupported value`),
		regexp.MustCompile(`error parsing regexp`),
		regexp.MustCompile(`time: missing unit in duration`),
		regexp.MustCompile(`using interface \{} as type .*`),
		regexp.MustCompile(`reflect.Value.MapIndex: value of type .* is not assignable to type .*`),
		regexp.MustCompile(`reflect: call of reflect.Value.Call on zero Value`),
		regexp.MustCompile(`reflect: call of reflect.Value.Len on bool Value`),
		regexp.MustCompile(`reflect: Call using .* as type .*`),
		regexp.MustCompile(`reflect: call of reflect.Value.Index on map Value`),
	}

	skipCode := []string{
		`??`,
	}

	f.Fuzz(func(t *testing.T, code string) {
		for _, skipCase := range skipCode {
			if strings.Contains(code, skipCase) {
				t.Skip()
				return
			}
		}

		program, err := expr.Compile(code, expr.Env(playground.Blog{}))
		if err != nil {
			t.Skip()
		}

		_, err = expr.Run(program, env)
		if err != nil {
			for _, okCase := range okCases {
				if okCase.MatchString(err.Error()) {
					t.Skip()
					return
				}
			}
			t.Errorf("code: %s, err: %s", code, err)
		}
	})
}
