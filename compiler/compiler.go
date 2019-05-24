package compiler

import (
	"encoding/binary"
	"fmt"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/vm"
	"math"
)

func Compile(node ast.Node) (program vm.Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	c := &compiler{}
	c.compile(node)

	program = vm.Program{
		Bytecode: c.bytecode,
		Constant: c.constant,
	}
	return
}

type compiler struct {
	bytecode []byte
	constant []interface{}
}

func (c *compiler) emit(op byte, b ...byte) {
	c.bytecode = append(c.bytecode, op)
	c.bytecode = append(c.bytecode, b...)
}

func (c *compiler) makeConstant(i interface{}) []byte {
	c.constant = append(c.constant, i)

	if len(c.constant) > math.MaxUint16 {
		panic("exceeded constants max space limit")
	}

	return i64(int64(len(c.constant) - 1))
}

func (c *compiler) placeholder() []int {
	c.emit(0xFF, 0xFF)
	return []int{len(c.bytecode) - 2, len(c.bytecode) - 1}
}

func (c *compiler) patchJump(placeholder []int) {
	current := len(c.bytecode) - 1
	b := i64(int64(current))

	for i, ip := range placeholder {
		c.bytecode[ip] = b[i]
	}
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
	c.emit(vm.OpNil)
}

func (c *compiler) IdentifierNode(node *ast.IdentifierNode) {
	c.emit(vm.OpFetch, c.makeConstant(node.Value)...)
}

func (c *compiler) IntegerNode(node *ast.IntegerNode) {
	if node.Value <= math.MaxUint16 {
		c.emit(vm.OpPush, i64(node.Value)...)
	} else {
		c.emit(vm.OpLoad, c.makeConstant(node.Value)...)
	}
}

func (c *compiler) FloatNode(node *ast.FloatNode) {
	c.emit(vm.OpLoad, c.makeConstant(node.Value)...)
}

func (c *compiler) BoolNode(node *ast.BoolNode) {
	if node.Value {
		c.emit(vm.OpTrue)
	} else {
		c.emit(vm.OpFalse)
	}
}

func (c *compiler) StringNode(node *ast.StringNode) {
	c.emit(vm.OpLoad, c.makeConstant(node.Value)...)
}

func (c *compiler) UnaryNode(node *ast.UnaryNode) {
	c.compile(node.Node)

	switch node.Operator {

	case "!", "not":
		c.emit(vm.OpNot)

	case "+":
		// Do nothing

	case "-":
		c.emit(vm.OpNegate)

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))
	}
}

func (c *compiler) BinaryNode(node *ast.BinaryNode) {

	switch node.Operator {
	case "==":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(vm.OpEqual)

	case "!=":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(vm.OpNotEqual)

	case "or", "||":
		c.compile(node.Left)
		c.emit(vm.OpJumpIfTrue)
		end := c.placeholder()
		c.emit(vm.OpPop)
		c.compile(node.Right)
		c.patchJump(end)

	case "and", "&&":
		c.compile(node.Left)
		c.emit(vm.OpJumpIfFalse)
		end := c.placeholder()
		c.emit(vm.OpPop)
		c.compile(node.Right)
		c.patchJump(end)

	case "in", "not in":

	case "<", ">", ">=", "<=":

	case "/", "-", "*", "**":

	case "%":

	case "+":

	case "contains":

	case "..":

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))

	}
}

func (c *compiler) MatchesNode(node *ast.MatchesNode) {}

func (c *compiler) PropertyNode(node *ast.PropertyNode) {}

func (c *compiler) IndexNode(node *ast.IndexNode) {}

func (c *compiler) MethodNode(node *ast.MethodNode) {}

func (c *compiler) FunctionNode(node *ast.FunctionNode) {}

func (c *compiler) BuiltinNode(node *ast.BuiltinNode) {}

func (c *compiler) ClosureNode(node *ast.ClosureNode) {}

func (c *compiler) PointerNode(node *ast.PointerNode) {}

func (c *compiler) ConditionalNode(node *ast.ConditionalNode) {}

func (c *compiler) ArrayNode(node *ast.ArrayNode) {}

func (c *compiler) MapNode(node *ast.MapNode) {}

func i64(i int64) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}
