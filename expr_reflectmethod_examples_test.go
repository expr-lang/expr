//go:build !expr_noreflectmethod

package expr_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/test/mock"
)

func ExampleEval() {
	output, err := expr.Eval("greet + name", map[string]any{
		"greet": "Hello, ",
		"name":  "world!",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: Hello, world!
}

func ExampleEval_runtime_error() {
	_, err := expr.Eval(`map(1..3, {1 % (# - 3)})`, nil)
	fmt.Print(err)

	// Output: runtime error: integer divide by zero (1:14)
	//  | map(1..3, {1 % (# - 3)})
	//  | .............^
}

func ExampleEval_bytes_literal() {
	// Bytes literal returns []byte.
	output, err := expr.Eval(`b"abc"`, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: [97 98 99]
}

func ExampleEnv_tagged_field_names() {
	env := struct {
		FirstWord  string
		Separator  string `expr:"Space"`
		SecondWord string `expr:"second_word"`
	}{
		FirstWord:  "Hello",
		Separator:  " ",
		SecondWord: "World",
	}

	output, err := expr.Eval(`FirstWord + Space + second_word`, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: Hello World
}

func ExampleEnv_hidden_tagged_field_names() {
	type Internal struct {
		Visible string
		Hidden  string `expr:"-"`
	}
	type environment struct {
		Visible         string
		Hidden          string   `expr:"-"`
		HiddenInternal  Internal `expr:"-"`
		VisibleInternal Internal
	}

	env := environment{
		Hidden: "First level secret",
		HiddenInternal: Internal{
			Visible: "Second level secret",
			Hidden:  "Also hidden",
		},
		VisibleInternal: Internal{
			Visible: "Not a secret",
			Hidden:  "Hidden too",
		},
	}

	hiddenValues := []string{
		`Hidden`,
		`HiddenInternal`,
		`HiddenInternal.Visible`,
		`HiddenInternal.Hidden`,
		`VisibleInternal["Hidden"]`,
	}
	for _, expression := range hiddenValues {
		output, err := expr.Eval(expression, env)
		if err == nil || !strings.Contains(err.Error(), "cannot fetch") {
			fmt.Printf("unexpected output: %v; err: %v\n", output, err)
			return
		}
		fmt.Printf("%q is hidden as expected\n", expression)
	}

	visibleValues := []string{
		`Visible`,
		`VisibleInternal`,
		`VisibleInternal["Visible"]`,
	}
	for _, expression := range visibleValues {
		_, err := expr.Eval(expression, env)
		if err != nil {
			fmt.Printf("unexpected error: %v\n", err)
			return
		}
		fmt.Printf("%q is visible as expected\n", expression)
	}

	testWithIn := []string{
		`not ("Hidden" in $env)`,
		`"Visible" in $env`,
		`not ("Hidden" in VisibleInternal)`,
		`"Visible" in VisibleInternal`,
	}
	for _, expression := range testWithIn {
		val, err := expr.Eval(expression, env)
		shouldBeTrue, ok := val.(bool)
		if err != nil || !ok || !shouldBeTrue {
			fmt.Printf("unexpected result; value: %v; error: %v\n", val, err)
			return
		}
	}

	// Output: "Hidden" is hidden as expected
	// "HiddenInternal" is hidden as expected
	// "HiddenInternal.Visible" is hidden as expected
	// "HiddenInternal.Hidden" is hidden as expected
	// "VisibleInternal[\"Hidden\"]" is hidden as expected
	// "Visible" is visible as expected
	// "VisibleInternal" is visible as expected
	// "VisibleInternal[\"Visible\"]" is visible as expected
}

func ExampleOperator() {
	code := `
		Now() > CreatedAt &&
		(Now() - CreatedAt).Hours() > 24
	`

	type Env struct {
		CreatedAt time.Time
		Now       func() time.Time
		Sub       func(a, b time.Time) time.Duration
		After     func(a, b time.Time) bool
	}

	options := []expr.Option{
		expr.Env(Env{}),
		expr.Operator(">", "After"),
		expr.Operator("-", "Sub"),
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	env := Env{
		CreatedAt: time.Date(2018, 7, 14, 0, 0, 0, 0, time.UTC),
		Now:       func() time.Time { return time.Now() },
		Sub:       func(a, b time.Time) time.Duration { return a.Sub(b) },
		After:     func(a, b time.Time) bool { return a.After(b) },
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)

	// Output: true
}

func ExampleAllowUndefinedVariables_zero_value_functions() {
	code := `words == "" ? Split("foo,bar", ",") : Split(words, ",")`

	// Env is map[string]string type on which methods are defined.
	env := mock.MapStringStringEnv{}

	options := []expr.Option{
		expr.Env(env),
		expr.AllowUndefinedVariables(), // Allow to use undefined variables.
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%v", output)

	// Output: [foo bar]
}

func ExampleTimezone() {
	program, err := expr.Compile(`now().Location().String()`, expr.Timezone("Asia/Kamchatka"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v", output)
	// Output: Asia/Kamchatka
}
