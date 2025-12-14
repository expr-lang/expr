package main

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/builtin"
)

var (
	dict       []string
	predicates []string
	builtins   []string
	operators  = []string{
		"or",
		"||",
		"and",
		"&&",
		"==",
		"!=",
		"<",
		">",
		">=",
		"<=",
		"..",
		"??",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"^",
		"in",
		"matches",
		"contains",
		"startsWith",
		"endsWith",
		"not in",
		"not matches",
		"not contains",
		"not startsWith",
		"not endsWith",
	}
)

func init() {
	for name, x := range Env {
		dict = append(dict, name)
		v := reflect.ValueOf(x)
		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				dict = append(dict, v.Type().Field(i).Name)
			}
			for i := 0; i < v.NumMethod(); i++ {
				dict = append(dict, v.Type().Method(i).Name)
			}
		}
		if v.Kind() == reflect.Map {
			for _, key := range v.MapKeys() {
				dict = append(dict, fmt.Sprintf("%v", key.Interface()))
			}
		}
	}
	for _, b := range builtin.Builtins {
		if b.Predicate {
			predicates = append(predicates, b.Name)
		} else {
			builtins = append(builtins, b.Name)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var corpus = make(map[string]struct{})
	var corpusMutex sync.Mutex

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Worker %d recovered from panic: %v\n", workerID, r)
					debug.PrintStack()
				}
			}()

			defer wg.Done()
			for {
				var code string

				code = node(oneOf(list[int]{
					{3, 100},
					{4, 40},
					{5, 50},
					{6, 30},
					{7, 20},
					{8, 10},
					{9, 5},
					{10, 5},
				}))

				program, err := expr.Compile(code, expr.Env(Env))
				if err != nil {
					continue
				}
				_, err = expr.Run(program, Env)
				if err != nil {
					continue
				}

				corpusMutex.Lock()
				if _, exists := corpus[code]; exists {
					corpusMutex.Unlock()
					continue
				}
				corpus[code] = struct{}{}
				corpusMutex.Unlock()

				fmt.Println(code)
			}
		}(i)
	}

	wg.Wait()
}

type fn func(depth int) string

func node(depth int) string {
	if depth <= 0 {
		return oneOf(list[fn]{
			{nilNode, 1},
			{envNode, 1},
			{floatNode, 1},
			{integerNode, 1},
			{stringNode, 1},
			{booleanNode, 1},
			{identifierNode, 10},
			{pointerNode, 10},
		})(depth - 1)
	}
	return oneOf(list[fn]{
		{sequenceNode, 1},
		{variableNode, 1},
		{arrayNode, 10},
		{mapNode, 10},
		{identifierNode, 1000},
		{memberNode, 1500},
		{unaryNode, 100},
		{binaryNode, 2000},
		{callNode, 2000},
		{pipeNode, 1000},
		{builtinNode, 500},
		{predicateNode, 1000},
		{pointerNode, 500},
		{sliceNode, 100},
		{conditionalNode, 100},
	})(depth - 1)
}

func nilNode(_ int) string {
	return "nil"
}

func envNode(_ int) string {
	return "$env"
}

func floatNode(_ int) string {
	return "1.0"
}

func integerNode(_ int) string {
	return oneOf(list[string]{
		{"1", 1},
		{"0", 1},
	})
}

func stringNode(_ int) string {
	return "foo"
}

func booleanNode(_ int) string {
	return random([]string{"true", "false"})
}

func identifierNode(_ int) string {
	if maybe() {
		return "foobar"
	}
	return random(dict)
}

func memberNode(depth int) string {
	dot := "."
	if maybe() {
		dot = "?."
	}
	prop := oneOf(list[fn]{
		{func(_ int) string { return random(dict) }, 5},
		{node, 1},
	})(depth - 1)
	if maybe() {
		return fmt.Sprintf("%v%v%v", node(depth-1), dot, prop)
	}
	return fmt.Sprintf("%v%v[%v]", node(depth-1), dot, prop)
}

func unaryNode(depth int) string {
	op := random([]string{"-", "!", "not"})
	// Use a simple formatting to ensure valid unary expression syntax
	if op == "not" {
		return fmt.Sprintf("not %v", node(depth-1))
	}
	return fmt.Sprintf("%s%v", op, node(depth-1))
}

func binaryNode(depth int) string {
	return fmt.Sprintf("%v %v %v", node(depth-1), random(operators), node(depth-1))
}

func methodNode(depth int) string {
	dot := "."
	if maybe() {
		dot = "?."
	}
	method := random(dict)
	if maybe() {
		return fmt.Sprintf("%v%v%v", node(depth-1), dot, method)
	}
	return fmt.Sprintf("%v%v[%v]", node(depth-1), dot, method)
}

func funcNode(_ int) string {
	return random(dict)
}

func callNode(depth int) string {
	var args []string
	for i := 0; i < oneOf(list[int]{
		{0, 100},
		{1, 100},
		{2, 50},
		{3, 25},
		{4, 10},
		{5, 5},
	}); i++ {
		args = append(args, node(depth-1))
	}

	fn := oneOf(list[fn]{
		{methodNode, 2},
		{funcNode, 2},
	})(depth - 1)

	return fmt.Sprintf("%v(%v)", fn, strings.Join(args, ", "))
}

func pipeNode(depth int) string {
	a := node(depth - 1)
	b := oneOf(list[fn]{
		{callNode, 2},
		{builtinNode, 5},
		{predicateNode, 10},
	})(depth - 1)

	return fmt.Sprintf("%v | %v", a, b)
}

func builtinNode(depth int) string {
	var args []string
	for i := 0; i < oneOf(list[int]{
		{1, 100},
		{2, 50},
		{3, 50},
		{4, 10},
	}); i++ {
		args = append(args, node(depth-1))
	}
	return fmt.Sprintf("%v(%v)", random(builtins), strings.Join(args, ", "))
}

func predicateNode(depth int) string {
	var args []string
	for i := 0; i < oneOf(list[int]{
		{1, 100},
		{2, 50},
		{3, 50},
	}); i++ {
		args = append(args, node(depth-1))
	}
	return fmt.Sprintf("%v(%v)", random(predicates), strings.Join(args, ", "))
}

func pointerNode(_ int) string {
	return oneOf(list[string]{
		{"#", 100},
		{"#." + random(dict), 100},
		{"." + random(dict), 100},
		{"#acc", 10},
		{"#index", 10},
	})
}

func arrayNode(depth int) string {
	var items []string
	for i := 0; i < oneOf(list[int]{
		{1, 100},
		{2, 50},
		{3, 25},
	}); i++ {
		items = append(items, node(depth-1))
	}
	return fmt.Sprintf("[%v]", strings.Join(items, ", "))
}

func mapNode(depth int) string {
	var items []string
	for i := 0; i < oneOf(list[int]{
		{1, 100},
		{2, 50},
		{3, 25},
	}); i++ {
		items = append(items, fmt.Sprintf("%v: %v", stringNode(depth-1), node(depth-1)))
	}
	return fmt.Sprintf("{%v}", strings.Join(items, ", "))
}

func sliceNode(depth int) string {
	return oneOf(list[string]{
		{fmt.Sprintf("%v[%v:%v]", node(depth-1), node(depth-1), node(depth-1)), 100},
		{fmt.Sprintf("%v[%v:]", node(depth-1), node(depth-1)), 100},
		{fmt.Sprintf("%v[:%v]", node(depth-1), node(depth-1)), 100},
		{fmt.Sprintf("%v[:]", node(depth-1)), 1},
	})
}

func conditionalNode(depth int) string {
	return oneOf(list[string]{
		{fmt.Sprintf("if %v { %v } else { %v }", node(depth-1), node(depth-1), node(depth-1)), 100},
		{fmt.Sprintf("%v ? %v : %v", node(depth-1), node(depth-1), node(depth-1)), 100},
		{fmt.Sprintf("%v ?: %v", node(depth-1), node(depth-1)), 20},
	})
}

func sequenceNode(depth int) string {
	var items []string
	for i := 0; i < oneOf(list[int]{
		{2, 50},
		{3, 25},
	}); i++ {
		items = append(items, node(depth-1))
	}
	if maybe() {
		return strings.Join(items, "; ")
	}
	return fmt.Sprintf("(%v)", strings.Join(items, ", "))
}

func variableNode(depth int) string {
	// Generate 1-3 variable declarations with diverse names and then make sure
	// at least one of them is used in the final expression.
	namesPool := []string{"x", "y", "z", "foo", "bar", "foobar", "tmp"}
	nDecls := oneOf(list[int]{
		{1, 60},
		{2, 30},
		{3, 10},
	})

	var decls []string
	var chosen []string
	for i := 0; i < nDecls; i++ {
		name := random(namesPool)
		chosen = append(chosen, name)
		decls = append(decls, fmt.Sprintf("let %s = %v", name, node(depth-1)))
	}

	// Build a usage expression that references declared vars to guarantee coverage.
	use := oneOf(list[string]{
		{chosen[0], 50},
		{fmt.Sprintf("%s + %v", chosen[0], node(depth-1)), 30},
		{fmt.Sprintf("%v + %s", node(depth-1), chosen[0]), 30},
		{fmt.Sprintf("if %v { %s } else { %v }", node(depth-1), chosen[0], node(depth-1)), 20},
		{fmt.Sprintf("%s ? %v : %v", chosen[0], node(depth-1), node(depth-1)), 10},
	})

	return fmt.Sprintf("%s; %s", strings.Join(decls, "; "), use)
}
