package main

import (
	"fmt"
	"math/rand"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/test/playground"
)

func main() {
	var code string
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("==========================\n%s\n==========================\n", code)
			panic(r)
		}
	}()

	corpus := map[string]struct{}{}
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

		program, err := expr.Compile(code, expr.Env(playground.Blog{}))
		if err != nil {
			continue
		}
		_, err = expr.Run(program, playground.ExampleData())
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

func nilNode(depth int) ast.Node {
	return &ast.NilNode{}
}

func floatNode(depth int) ast.Node {
	cases := []float64{
		0.0,
		0.5,
	}
	return &ast.FloatNode{
		Value: cases[rand.Intn(len(cases))],
	}
}

func integerNode(depth int) ast.Node {
	return &ast.IntegerNode{
		Value: rand.Intn(3),
	}
}

func stringNode(depth int) ast.Node {
	corpus := []string{
		"Go", "JavaScript", " ",
	}
	return &ast.StringNode{
		Value: corpus[rand.Intn(len(corpus))],
	}
}

func booleanNode(depth int) ast.Node {
	return &ast.BoolNode{
		Value: maybe(),
	}
}

func identifierNode(depth int) ast.Node {
	cases := []string{
		"Posts",
		"Authors",
		"TotalViews",
	}

	return &ast.IdentifierNode{
		Value: cases[rand.Intn(len(cases))],
	}
}

func memberNode(depth int) ast.Node {
	cases := []string{
		"Birthday",
		"Biography",
		"Website",
		"ID",
		"FirstName",
		"LastName",
		"Email",
		"Profile",
		"ID",
		"Title",
		"Content",
		"PublishDate",
		"Author",
		"Comments",
		"Tags",
		"Likes",
		"ID",
		"AuthorName",
		"Content",
		"CommentDate",
		"Upvotes",
	}

	return &ast.MemberNode{
		Node: node(depth - 1),
		Property: weightedRandom([]fnWeight{
			{func(_ int) ast.Node { return &ast.StringNode{Value: cases[rand.Intn(len(cases))]} }, 5},
			{node, 1},
		})(depth - 1),
		Optional: maybe(),
	}
}

func unaryNode(depth int) ast.Node {
	cases := []string{
		"-",
		"!",
		"not",
	}
	return &ast.UnaryNode{
		Operator: cases[rand.Intn(len(cases))],
		Node:     node(depth - 1),
	}
}

func binaryNode(depth int) ast.Node {
	cases := []string{
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
		"in",
		"matches",
		"contains",
		"startsWith",
		"endsWith",
		"..",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"^",
	}
	return &ast.BinaryNode{
		Operator: cases[rand.Intn(len(cases))],
		Left:     node(depth - 1),
		Right:    node(depth - 1),
	}
}

func methodNode(depth int) ast.Node {
	cases := []string{
		"Age",
		"FullName",
		"IsAdmin",
		"Published",
		"After",
		"Before",
		"Compare",
		"Equal",
		"IsZero",
		"Upvoted",
		"AuthorEmail",
	}

	return &ast.MemberNode{
		Node:     node(depth - 1),
		Property: &ast.StringNode{Value: cases[rand.Intn(len(cases))]},
		Optional: maybe(),
	}
}

func funcNode(depth int) ast.Node {
	cases := []string{
		"RecentPosts",
		"PopularPosts",
		"TotalUpvotes",
		"TotalComments",
		"Add",
		"Sub",
		"Title",
		"HasTag",
		"IsAdmin",
		"IsZero",
		"WithID",
	}

	return &ast.IdentifierNode{
		Value: cases[rand.Intn(len(cases))],
	}
}

func callNode(depth int) ast.Node {
	var args []ast.Node
	max := weightedRandomInt([]intWeight{
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
		{3, 25},
		{4, 10},
		{5, 5},
	})
	for i := 0; i < max; i++ {
		args = append(args, node(depth-1))
	}
	return &ast.BuiltinNode{
		Name:      builtin.Names[rand.Intn(len(builtin.Names))],
		Arguments: args,
	}
}

func predicateNode(depth int) ast.Node {
	cases := []string{
		"all",
		"none",
		"any",
		"one",
		"filter",
		"map",
		"count",
	}
	return &ast.BuiltinNode{
		Name:      cases[rand.Intn(len(cases))],
		Arguments: []ast.Node{node(depth - 1), node(depth - 1)},
	}
}

func pointerNode(depth int) ast.Node {
	return &ast.PointerNode{}
}

func arrayNode(depth int) ast.Node {
	var items []ast.Node
	max := weightedRandomInt([]intWeight{
		{1, 100},
		{2, 50},
		{3, 25},
		{4, 10},
		{5, 5},
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
		{4, 10},
		{5, 5},
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
