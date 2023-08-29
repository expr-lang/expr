package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime/debug"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/builtin"
)

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

type Foo struct {
	Bar string
}

func (f Foo) String() string {
	return "foo"
}

func (f Foo) Qux(s string) string {
	return f.Bar + s
}

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
	for name, x := range env {
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
	var code string
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("==========================\n%s\n==========================\n%s\n==========================\n", code, r)
			debug.PrintStack()
		}
	}()

	var corpus = map[string]struct{}{}

	for {
		code = node(weightedRandomInt([]intWeight{
			{3, 100},
			{4, 40},
			{5, 50},
			{6, 30},
			{7, 20},
			{8, 10},
			{9, 5},
			{10, 5},
		})).String()

		program, err := expr.Compile(code, expr.Env(env))
		if err != nil {
			continue
		}
		_, err = expr.Run(program, env)
		if err != nil {
			continue
		}

		if _, ok := corpus[code]; ok {
			continue
		}
		corpus[code] = struct{}{}
		fmt.Println(code)
	}
}

func node(depth int) ast.Node {
	if depth <= 0 {
		return weightedRandom([]fnWeight{
			{nilNode, 1},
			{floatNode, 1},
			{integerNode, 1},
			{stringNode, 1},
			{booleanNode, 1},
			{identifierNode, 10},
			{pointerNode, 10},
		})(depth - 1)
	}
	return weightedRandom([]fnWeight{
		{arrayNode, 1},
		{mapNode, 1},
		{identifierNode, 1000},
		{memberNode, 1500},
		{unaryNode, 100},
		{binaryNode, 2000},
		{callNode, 2000},
		{builtinNode, 500},
		{predicateNode, 1000},
		{pointerNode, 500},
		{sliceNode, 100},
		{conditionalNode, 100},
	})(depth - 1)
}

func nilNode(_ int) ast.Node {
	return &ast.NilNode{}
}

func floatNode(_ int) ast.Node {
	return &ast.FloatNode{
		Value: .5,
	}
}

func integerNode(_ int) ast.Node {
	return &ast.IntegerNode{
		Value: 1,
	}
}

func stringNode(_ int) ast.Node {
	words := []string{
		"foo",
		"bar",
	}
	return &ast.StringNode{
		Value: words[rand.Intn(len(words))],
	}
}

func booleanNode(_ int) ast.Node {
	return &ast.BoolNode{
		Value: maybe(),
	}
}

func identifierNode(_ int) ast.Node {
	return &ast.IdentifierNode{
		Value: dict[rand.Intn(len(dict))],
	}
}

func memberNode(depth int) ast.Node {
	return &ast.MemberNode{
		Node: node(depth - 1),
		Property: weightedRandom([]fnWeight{
			{func(_ int) ast.Node { return &ast.StringNode{Value: dict[rand.Intn(len(dict))]} }, 5},
			{node, 1},
		})(depth - 1),
		Optional: maybe(),
	}
}

func unaryNode(depth int) ast.Node {
	cases := []string{"-", "!", "not"}
	return &ast.UnaryNode{
		Operator: cases[rand.Intn(len(cases))],
		Node:     node(depth - 1),
	}
}

func binaryNode(depth int) ast.Node {
	return &ast.BinaryNode{
		Operator: operators[rand.Intn(len(operators))],
		Left:     node(depth - 1),
		Right:    node(depth - 1),
	}
}

func methodNode(depth int) ast.Node {
	return &ast.MemberNode{
		Node:     node(depth - 1),
		Property: &ast.StringNode{Value: dict[rand.Intn(len(dict))]},
		Optional: maybe(),
	}
}

func funcNode(_ int) ast.Node {
	return &ast.IdentifierNode{
		Value: dict[rand.Intn(len(dict))],
	}
}

func callNode(depth int) ast.Node {
	var args []ast.Node
	max := weightedRandomInt([]intWeight{
		{0, 100},
		{1, 100},
		{2, 50},
		{3, 25},
		{4, 10},
		{5, 5},
	})
	for i := 0; i < max; i++ {
		args = append(args, node(depth-1))
	}
	return &ast.CallNode{
		Callee: weightedRandom([]fnWeight{
			{methodNode, 2},
			{funcNode, 2},
		})(depth - 1),
		Arguments: args,
	}
}

func builtinNode(depth int) ast.Node {
	var args []ast.Node
	max := weightedRandomInt([]intWeight{
		{1, 100},
		{2, 50},
		{3, 50},
		{4, 10},
	})
	for i := 0; i < max; i++ {
		args = append(args, node(depth-1))
	}
	return &ast.BuiltinNode{
		Name:      builtins[rand.Intn(len(builtins))],
		Arguments: args,
	}
}

func predicateNode(depth int) ast.Node {
	return &ast.BuiltinNode{
		Name: predicates[rand.Intn(len(predicates))],
		Arguments: []ast.Node{
			node(depth - 1),
			node(depth - 1),
		},
	}
}

func pointerNode(_ int) ast.Node {
	return &ast.PointerNode{}
}

func arrayNode(depth int) ast.Node {
	var items []ast.Node
	max := weightedRandomInt([]intWeight{
		{1, 100},
		{2, 50},
		{3, 25},
	})
	for i := 0; i < max; i++ {
		items = append(items, node(depth-1))
	}
	return &ast.ArrayNode{
		Nodes: items,
	}
}

func mapNode(depth int) ast.Node {
	var items []ast.Node
	max := weightedRandomInt([]intWeight{
		{1, 100},
		{2, 50},
		{3, 25},
	})
	for i := 0; i < max; i++ {
		items = append(items, &ast.PairNode{
			Key:   stringNode(depth - 1),
			Value: node(depth - 1),
		})
	}
	return &ast.MapNode{
		Pairs: items,
	}
}

func sliceNode(depth int) ast.Node {
	return &ast.SliceNode{
		Node: node(depth - 1),
		From: node(depth - 1),
		To:   node(depth - 1),
	}
}

func conditionalNode(depth int) ast.Node {
	return &ast.ConditionalNode{
		Cond: node(depth - 1),
		Exp1: node(depth - 1),
		Exp2: node(depth - 1),
	}
}
