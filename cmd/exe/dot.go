package main

import (
	"fmt"
	. "gopkg.in/antonmedv/expr.v2/ast"
	"os"
)

const format = `digraph {
	ranksep=.3;

	node [shape=box, fixedsize=false, fontsize=12, fontname="Helvetica-bold", fontcolor="#1259FF", width=.25, height=.25, color="black", fillcolor="white", style="filled, solid, bold"];
	edge [arrowsize=.5, color="black", style="bold"]

%v
%v}
`

func dotAst(node Node) {
	v := &visitor{}
	Walk(&node, v)
	dot := fmt.Sprintf(format, v.nodes, v.links)
	_, _ = fmt.Fprintf(os.Stdout, dot)
}

type visitor struct {
	BaseVisitor
	nodes string
	links string
	index int
	stack []int
}

func (v *visitor) push(label string) {
	v.index++
	v.nodes += fmt.Sprintf("	n%v [label=%q];\n", v.index, label)
	v.stack = append(v.stack, v.index)
}

func (v *visitor) pop() int {
	node := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]
	return node
}

func (v *visitor) link(node int) {
	v.links += fmt.Sprintf("	n%v -> n%v\n", v.index, node)
}

func (v *visitor) NilNode(node *NilNode) {
	v.push("nil")
}

func (v *visitor) IdentifierNode(node *IdentifierNode) {
	v.push(node.Value)
}

func (v *visitor) IntegerNode(node *IntegerNode) {
	v.push(fmt.Sprintf("%v", node.Value))
}

func (v *visitor) FloatNode(node *FloatNode) {
	v.push(fmt.Sprintf("%v", node.Value))
}

func (v *visitor) BoolNode(node *BoolNode) {
	v.push(fmt.Sprintf("%v", node.Value))
}

func (v *visitor) StringNode(node *StringNode) {
	v.push(fmt.Sprintf("%q", node.Value))
}

func (v *visitor) UnaryNode(node *UnaryNode) {
	n := v.pop()
	v.push("-")
	v.link(n)
}

func (v *visitor) BinaryNode(node *BinaryNode) {
	b := v.pop()
	a := v.pop()
	v.push(node.Operator)
	v.link(a)
	v.link(b)
}

func (v *visitor) MatchesNode(node *MatchesNode) {
	b := v.pop()
	a := v.pop()
	v.push("matches")
	v.link(a)
	v.link(b)
}

func (v *visitor) PropertyNode(node *PropertyNode) {
	a := v.pop()
	v.push(fmt.Sprintf(".%v", node.Property))
	v.link(a)
}

func (v *visitor) IndexNode(node *IndexNode) {
	b := v.pop()
	a := v.pop()
	v.push("[...]")
	v.link(a)
	v.link(b)
}

func (v *visitor) MethodNode(node *MethodNode) {
	args := make([]int, 0)
	for range node.Arguments {
		args = append(args, v.pop())
	}
	a := v.pop()
	v.push(fmt.Sprintf(".%v(...)", node.Method))
	v.link(a)
	for i := len(args) - 1; i >= 0; i-- {
		v.link(args[i])
	}
}

func (v *visitor) FunctionNode(node *FunctionNode) {
	args := make([]int, 0)
	for range node.Arguments {
		args = append(args, v.pop())
	}
	v.push(fmt.Sprintf("%v(...)", node.Name))
	for i := len(args) - 1; i >= 0; i-- {
		v.link(args[i])
	}
}

func (v *visitor) BuiltinNode(node *BuiltinNode) {
	args := make([]int, 0)
	for range node.Arguments {
		args = append(args, v.pop())
	}
	v.push(fmt.Sprintf("%v(...)", node.Name))
	for i := len(args) - 1; i >= 0; i-- {
		v.link(args[i])
	}
}

func (v *visitor) ClosureNode(node *ClosureNode) {
	a := v.pop()
	v.push("func {...}")
	v.link(a)
}

func (v *visitor) PointerNode(node *PointerNode) {
	v.push("#")
}

func (v *visitor) ConditionalNode(node *ConditionalNode) {
	e2 := v.pop()
	e1 := v.pop()
	c := v.pop()
	v.push("? :")
	v.link(c)
	v.link(e1)
	v.link(e2)
}

func (v *visitor) ArrayNode(node *ArrayNode) {
	n := make([]int, 0)
	for range node.Nodes {
		n = append(n, v.pop())
	}
	v.push("[...]")
	for i := len(n) - 1; i >= 0; i-- {
		v.link(n[i])
	}
}

func (v *visitor) MapNode(node *MapNode) {
	n := make([]int, 0)
	for range node.Pairs {
		n = append(n, v.pop())
	}
	v.push("{...}")
	for i := len(n) - 1; i >= 0; i-- {
		v.link(n[i])
	}
}

func (v *visitor) PairNode(node *PairNode) {
	a := v.pop()
	v.push(fmt.Sprintf("%q:", node.Key.(*StringNode).Value))
	v.link(a)
}
