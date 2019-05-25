package compiler

import (
	"encoding/binary"
	"fmt"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/parser"
	. "github.com/antonmedv/expr/vm"
	"math"
)

func Compile(tree *parser.Tree) (program Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	c := &compiler{
		index: make(map[interface{}]int),
	}
	c.compile(tree.Node)

	program = Program{
		Bytecode: c.bytecode,
		Constant: c.constant,
	}
	return
}

type compiler struct {
	bytecode []byte
	constant []interface{}
	index    map[interface{}]int
}

func (c *compiler) emit(op byte, b ...byte) int {
	c.bytecode = append(c.bytecode, op)
	current := len(c.bytecode)
	c.bytecode = append(c.bytecode, b...)
	return current
}

func (c *compiler) makeConstant(i interface{}) []byte {
	if p, ok := c.index[i]; ok {
		return encode(p)
	}

	c.constant = append(c.constant, i)
	p := len(c.constant) - 1
	c.index[i] = p

	if len(c.constant) > math.MaxUint16 {
		panic("exceeded constants max space limit")
	}

	return encode(p)
}

func (c *compiler) placeholder() []byte {
	return []byte{0xFF, 0xFF}
}

func (c *compiler) patchJump(placeholder int) {
	offset := len(c.bytecode) - 2 - placeholder
	b := encode(offset)
	c.bytecode[placeholder] = b[0]
	c.bytecode[placeholder+1] = b[1]
}

func (c *compiler) calcBackwardJump(to int) []byte {
	return encode(len(c.bytecode) + 1 + 2 - to)
}

func (c *compiler) compile(node ast.Node) {
	switch n := node.(type) {
	case *ast.NilNode:
		c.NilNode(n)
	case *ast.IdentifierNode:
		c.IdentifierNode(n)
	case *ast.IntegerNode:
		c.IntegerNode(n)
	case *ast.FloatNode:
		c.FloatNode(n)
	case *ast.BoolNode:
		c.BoolNode(n)
	case *ast.StringNode:
		c.StringNode(n)
	case *ast.UnaryNode:
		c.UnaryNode(n)
	case *ast.BinaryNode:
		c.BinaryNode(n)
	case *ast.MatchesNode:
		c.MatchesNode(n)
	case *ast.PropertyNode:
		c.PropertyNode(n)
	case *ast.IndexNode:
		c.IndexNode(n)
	case *ast.MethodNode:
		c.MethodNode(n)
	case *ast.FunctionNode:
		c.FunctionNode(n)
	case *ast.BuiltinNode:
		c.BuiltinNode(n)
	case *ast.ClosureNode:
		c.ClosureNode(n)
	case *ast.PointerNode:
		c.PointerNode(n)
	case *ast.ConditionalNode:
		c.ConditionalNode(n)
	case *ast.ArrayNode:
		c.ArrayNode(n)
	case *ast.MapNode:
		c.MapNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}

func (c *compiler) NilNode(node *ast.NilNode) {
	c.emit(OpNil)
}

func (c *compiler) IdentifierNode(node *ast.IdentifierNode) {
	c.emit(OpFetch, c.makeConstant(node.Value)...)
}

func (c *compiler) IntegerNode(node *ast.IntegerNode) {
	if node.Value <= math.MaxUint16 {
		c.emit(OpPush, encode(int(node.Value))...)
	} else {
		c.emit(OpConst, c.makeConstant(node.Value)...)
	}
}

func (c *compiler) FloatNode(node *ast.FloatNode) {
	c.emit(OpConst, c.makeConstant(node.Value)...)
}

func (c *compiler) BoolNode(node *ast.BoolNode) {
	if node.Value {
		c.emit(OpTrue)
	} else {
		c.emit(OpFalse)
	}
}

func (c *compiler) StringNode(node *ast.StringNode) {
	c.emit(OpConst, c.makeConstant(node.Value)...)
}

func (c *compiler) UnaryNode(node *ast.UnaryNode) {
	c.compile(node.Node)

	switch node.Operator {

	case "!", "not":
		c.emit(OpNot)

	case "+":
		// Do nothing

	case "-":
		c.emit(OpNegate)

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))
	}
}

func (c *compiler) BinaryNode(node *ast.BinaryNode) {

	switch node.Operator {
	case "==":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpEqual)

	case "!=":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpEqual)
		c.emit(OpNot)

	case "or", "||":
		c.compile(node.Left)
		end := c.emit(OpJumpIfTrue, c.placeholder()...)
		c.emit(OpPop)
		c.compile(node.Right)
		c.patchJump(end)

	case "and", "&&":
		c.compile(node.Left)
		end := c.emit(OpJumpIfFalse, c.placeholder()...)
		c.emit(OpPop)
		c.compile(node.Right)
		c.patchJump(end)

	case "in":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpIn)

	case "not in":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpIn)
		c.emit(OpNot)

	case "<":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpLess)

	case ">":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpMore)

	case ">=":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpMoreOrEqual)

	case "<=":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpLessOrEqual)

	case "+":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpAdd)

	case "-":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpSubtract)

	case "*":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpMultiply)

	case "/":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpDivide)

	case "%":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpModulo)

	case "**":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpExponent)

	case "contains":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpContains)

	case "..":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpRange)

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))

	}
}

func (c *compiler) MatchesNode(node *ast.MatchesNode) {
	if node.Regexp != nil {
		c.compile(node.Left)
		c.emit(OpMatchesConst, c.makeConstant(node.Regexp)...)
		return
	}
	c.compile(node.Left)
	c.compile(node.Right)
	c.emit(OpMatches)
}

func (c *compiler) PropertyNode(node *ast.PropertyNode) {
	c.compile(node.Node)
	c.emit(OpProperty, c.makeConstant(node.Property)...)
}

func (c *compiler) IndexNode(node *ast.IndexNode) {
	c.compile(node.Node)
	c.compile(node.Index)
	c.emit(OpIndex)
}

func (c *compiler) MethodNode(node *ast.MethodNode) {
	c.compile(node.Node)
	for _, arg := range node.Arguments {
		c.compile(arg)
	}
	c.emit(OpMethod, c.makeConstant(Call{Name: node.Method, Size: len(node.Arguments)})...)
}

func (c *compiler) FunctionNode(node *ast.FunctionNode) {
	for _, arg := range node.Arguments {
		c.compile(arg)
	}
	c.emit(OpCall, c.makeConstant(Call{Name: node.Name, Size: len(node.Arguments)})...)
}

func (c *compiler) BuiltinNode(node *ast.BuiltinNode) {
	switch node.Name {
	case "len":
		c.compile(node.Arguments[0])
		c.emit(OpLen)
		c.emit(OpPop)

	case "all":

	case "none":

	case "any":

	case "one":

	case "filter":
		i := c.makeConstant("i")
		size := c.makeConstant("size")
		array := c.makeConstant("array")
		count := c.makeConstant("count")

		c.compile(node.Arguments[0])

		c.emit(OpBegin)
		c.emit(OpLen)
		c.emit(OpPush, encode(0)...)
		c.emit(OpStore, i...)
		c.emit(OpStore, size...)
		c.emit(OpStore, array...)

		c.emit(OpPush, encode(0)...)
		c.emit(OpStore, count...)

		cond := len(c.bytecode)
		c.emit(OpLoad, i...)
		c.emit(OpLoad, size...)
		c.emit(OpLess)
		end := c.emit(OpJumpIfFalse, c.placeholder()...)
		c.emit(OpPop)

		c.compile(node.Arguments[1])

		noInc := c.emit(OpJumpIfFalse, c.placeholder()...)
		c.emit(OpPop)

		c.emit(OpLoad, count...)
		c.emit(OpInc)
		c.emit(OpStore, count...)

		c.emit(OpLoad, c.makeConstant("array")...)
		c.emit(OpLoad, c.makeConstant("i")...)
		c.emit(OpIndex)
		jmp := c.emit(OpJump, c.placeholder()...)

		c.patchJump(noInc)
		c.emit(OpPop)

		c.patchJump(jmp)
		c.emit(OpLoad, i...)
		c.emit(OpInc)
		c.emit(OpStore, i...)
		c.emit(OpJumpBackward, c.calcBackwardJump(cond)...)

		c.patchJump(end)
		c.emit(OpPop)
		c.emit(OpLoad, count...)
		c.emit(OpEnd)

		c.emit(OpArray)

	case "map":

	default:
		panic(fmt.Sprintf("unknown builtin %v", node.Name))
	}
}

func (c *compiler) ClosureNode(node *ast.ClosureNode) {
	c.compile(node.Node)
}

func (c *compiler) PointerNode(node *ast.PointerNode) {
	c.emit(OpLoad, c.makeConstant("array")...)
	c.emit(OpLoad, c.makeConstant("i")...)
	c.emit(OpIndex)
}

func (c *compiler) ConditionalNode(node *ast.ConditionalNode) {
	c.compile(node.Cond)
	otherwise := c.emit(OpJumpIfFalse, c.placeholder()...)

	c.emit(OpPop)
	c.compile(node.Exp1)
	end := c.emit(OpJump, c.placeholder()...)

	c.patchJump(otherwise)
	c.emit(OpPop)
	c.compile(node.Exp2)

	c.patchJump(end)
}

func (c *compiler) ArrayNode(node *ast.ArrayNode) {
	for _, node := range node.Nodes {
		c.compile(node)
	}

	if len(node.Nodes) > math.MaxUint16 {
		panic("too big array")
	}

	c.emit(OpPush, encode(len(node.Nodes))...)
	c.emit(OpArray)
}

func (c *compiler) MapNode(node *ast.MapNode) {
	for _, pair := range node.Pairs {
		c.compile(pair.Key)
		c.compile(pair.Value)
	}

	if len(node.Pairs) > math.MaxUint16 {
		panic("too big array")
	}

	c.emit(OpPush, encode(len(node.Pairs))...)
	c.emit(OpMap)
}

func encode(i int) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}
