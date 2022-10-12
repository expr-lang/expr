package main

import (
	"fmt"
	"os"

	. "github.com/antonmedv/expr/ast"
)

const format = `digraph {
	ranksep=.3;
	node [shape=oval, fontname="Helvetica-bold"];
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

func (v *visitor) Enter(node *Node) {}
func (v *visitor) Exit(ref *Node) {
	switch node := (*ref).(type) {

	case *NilNode:
		v.push("nil")

	case *IdentifierNode:
		v.push(node.Value)

	case *IntegerNode:
		v.push(fmt.Sprintf("%v", node.Value))

	case *FloatNode:
		v.push(fmt.Sprintf("%v", node.Value))

	case *BoolNode:
		v.push(fmt.Sprintf("%v", node.Value))

	case *StringNode:
		v.push(fmt.Sprintf("%q", node.Value))

	case *UnaryNode:
		n := v.pop()
		v.push(node.Operator)
		v.link(n)

	case *BinaryNode:
		b := v.pop()
		a := v.pop()
		v.push(node.Operator)
		v.link(a)
		v.link(b)

	case *MatchesNode:
		b := v.pop()
		a := v.pop()
		v.push("matches")
		v.link(a)
		v.link(b)

	case *PropertyNode:
		a := v.pop()
		if !node.NilSafe {
			v.push(fmt.Sprintf(".%v", node.Property))
		} else {
			v.push(fmt.Sprintf("?.%v", node.Property))
		}
		v.link(a)

	case *IndexNode:
		b := v.pop()
		a := v.pop()
		v.push(fmt.Sprintf("%T", node))
		v.link(a)
		v.link(b)

	case *MethodNode:
		args := make([]int, 0)
		for range node.Arguments {
			args = append(args, v.pop())
		}
		a := v.pop()
		if !node.NilSafe {
			v.push(fmt.Sprintf(".%v(...)", node.Method))
		} else {
			v.push(fmt.Sprintf("?.%v(...)", node.Method))
		}
		v.link(a)
		for i := len(args) - 1; i >= 0; i-- {
			v.link(args[i])
		}

	case *FunctionNode:
		args := make([]int, 0)
		for range node.Arguments {
			args = append(args, v.pop())
		}
		v.push(fmt.Sprintf("%v(...)", node.Name))
		for i := len(args) - 1; i >= 0; i-- {
			v.link(args[i])
		}

	case *BuiltinNode:
		args := make([]int, 0)
		for range node.Arguments {
			args = append(args, v.pop())
		}
		v.push(fmt.Sprintf("%v", node.Name))
		for i := len(args) - 1; i >= 0; i-- {
			v.link(args[i])
		}

	case *ClosureNode:
		a := v.pop()
		v.push(fmt.Sprintf("%T", node))
		v.link(a)

	case *PointerNode:
		v.push("#")

	case *ConditionalNode:
		e2 := v.pop()
		e1 := v.pop()
		c := v.pop()
		v.push(fmt.Sprintf("%T", node))
		v.link(c)
		v.link(e1)
		v.link(e2)

	case *ArrayNode:
		n := make([]int, 0)
		for range node.Nodes {
			n = append(n, v.pop())
		}
		v.push("[...]")
		for i := len(n) - 1; i >= 0; i-- {
			v.link(n[i])
		}

	case *MapNode:
		n := make([]int, 0)
		for range node.Pairs {
			n = append(n, v.pop())
		}
		v.push("{...}")
		for i := len(n) - 1; i >= 0; i-- {
			v.link(n[i])
		}

	case *PairNode:
		b := v.pop()
		a := v.pop()
		v.push(fmt.Sprintf("%T", node))
		v.link(a)
		v.link(b)

	default:
		v.push(fmt.Sprintf("%T", node))
	}
}
