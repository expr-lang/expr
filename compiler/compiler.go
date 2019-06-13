package compiler

import (
	"encoding/binary"
	"fmt"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/internal/file"
	"github.com/antonmedv/expr/parser"
	. "github.com/antonmedv/expr/vm"
	"math"
	"reflect"
)

func Compile(tree *parser.Tree, ops ...OptionFn) (program *Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	c := &compiler{
		index: make(map[interface{}]uint16),
	}

	for _, op := range ops {
		op(c)
	}

	c.compile(tree.Node)

	program = &Program{
		Source:    tree.Source,
		Locations: c.locations,
		Constants: c.constants,
		Bytecode:  c.bytecode,
	}
	return
}

type compiler struct {
	locations   []file.Location
	constants   []interface{}
	bytecode    []byte
	index       map[interface{}]uint16
	mapEnv      bool
	currentNode ast.Node
}

// OptionFn for configuring expr.
type OptionFn func(c *compiler)

func MapEnv() OptionFn {
	return func(c *compiler) {
		c.mapEnv = true
	}
}

func (c *compiler) emit(op byte, b ...byte) int {
	c.bytecode = append(c.bytecode, op)
	current := len(c.bytecode)
	c.bytecode = append(c.bytecode, b...)

	for i := 0; i < 1+len(b); i++ {
		c.locations = append(c.locations, c.currentNode.GetLocation())
	}

	return current
}

func (c *compiler) makeConstant(i interface{}) []byte {
	hashable := true
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		hashable = false
	}

	if hashable {
		if p, ok := c.index[i]; ok {
			return encode(p)
		}
	}

	c.constants = append(c.constants, i)
	if len(c.constants) > math.MaxUint16 {
		panic("exceeded constants max space limit")
	}

	p := uint16(len(c.constants) - 1)
	if hashable {
		c.index[i] = p
	}
	return encode(p)
}

func (c *compiler) placeholder() []byte {
	return []byte{0xFF, 0xFF}
}

func (c *compiler) patchJump(placeholder int) {
	offset := len(c.bytecode) - 2 - placeholder
	b := encode(uint16(offset))
	c.bytecode[placeholder] = b[0]
	c.bytecode[placeholder+1] = b[1]
}

func (c *compiler) calcBackwardJump(to int) []byte {
	return encode(uint16(len(c.bytecode) + 1 + 2 - to))
}

func (c *compiler) compile(node ast.Node) {
	c.currentNode = node
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
	v := c.makeConstant(node.Value)
	if c.mapEnv {
		c.emit(OpFetchMap, v...)
	} else {
		c.emit(OpFetch, v...)
	}
}

func (c *compiler) IntegerNode(node *ast.IntegerNode) {
	t := node.GetType()
	if t == nil {
		if node.Value <= math.MaxUint16 {
			c.emit(OpPush, encode(uint16(node.Value))...)
		} else {
			c.emit(OpConst, c.makeConstant(node.Value)...)
		}
		return
	}

	switch t.Kind() {
	case reflect.Float32:
		c.emit(OpConst, c.makeConstant(float32(node.Value))...)
	case reflect.Float64:
		c.emit(OpConst, c.makeConstant(float64(node.Value))...)

	case reflect.Int:
		if node.Value <= math.MaxUint16 {
			c.emit(OpPush, encode(uint16(node.Value))...)
		} else {
			c.emit(OpConst, c.makeConstant(node.Value)...)
		}
	case reflect.Int8:
		c.emit(OpConst, c.makeConstant(int8(node.Value))...)
	case reflect.Int16:
		c.emit(OpConst, c.makeConstant(int16(node.Value))...)
	case reflect.Int32:
		c.emit(OpConst, c.makeConstant(int32(node.Value))...)
	case reflect.Int64:
		c.emit(OpConst, c.makeConstant(int64(node.Value))...)

	case reflect.Uint:
		c.emit(OpConst, c.makeConstant(uint(node.Value))...)
	case reflect.Uint8:
		c.emit(OpConst, c.makeConstant(uint8(node.Value))...)
	case reflect.Uint16:
		c.emit(OpConst, c.makeConstant(uint16(node.Value))...)
	case reflect.Uint32:
		c.emit(OpConst, c.makeConstant(uint32(node.Value))...)
	case reflect.Uint64:
		c.emit(OpConst, c.makeConstant(uint64(node.Value))...)

	default:
		panic(fmt.Sprintf("cannot compile %v to %v", node.Value, node.GetType()))
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

		l := kind(node.Left)
		r := kind(node.Right)

		if l == reflect.String && r == reflect.String {
			c.emit(OpEqualString)
		} else {
			c.emit(OpEqual)
		}

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

	case "startsWith":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpStartsWith)

	case "endsWith":
		c.compile(node.Left)
		c.compile(node.Right)
		c.emit(OpEndsWith)

	case "..":
		min, ok1 := node.Left.(*ast.IntegerNode)
		max, ok2 := node.Right.(*ast.IntegerNode)
		if ok1 && ok2 {
			// Create range on compile time to avoid unnecessary work at runtime.
			size := max.Value - min.Value + 1
			rng := make([]interface{}, size)
			for i := range rng {
				rng[i] = min.Value + i
			}
			c.emit(OpConst, c.makeConstant(rng)...)
			return
		}
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

	case "all":
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		var loopBreak int
		c.emitLoop(func() {
			c.compile(node.Arguments[1])
			loopBreak = c.emit(OpJumpIfFalse, c.placeholder()...)
			c.emit(OpPop)
		})
		c.emit(OpTrue)
		c.patchJump(loopBreak)
		c.emit(OpEnd)

	case "none":
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		var loopBreak int
		c.emitLoop(func() {
			c.compile(node.Arguments[1])
			c.emit(OpNot)
			loopBreak = c.emit(OpJumpIfFalse, c.placeholder()...)
			c.emit(OpPop)
		})
		c.emit(OpTrue)
		c.patchJump(loopBreak)
		c.emit(OpEnd)

	case "any":
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		var loopBreak int
		c.emitLoop(func() {
			c.compile(node.Arguments[1])
			loopBreak = c.emit(OpJumpIfTrue, c.placeholder()...)
			c.emit(OpPop)
		})
		c.emit(OpFalse)
		c.patchJump(loopBreak)
		c.emit(OpEnd)

	case "one":
		count := c.makeConstant("count")
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		c.emit(OpPush, encode(0)...)
		c.emit(OpStore, count...)
		c.emitLoop(func() {
			c.compile(node.Arguments[1])
			c.emitCond(func() {
				c.emit(OpLoad, count...)
				c.emit(OpInc)
				c.emit(OpStore, count...)
			})
		})
		c.emit(OpLoad, count...)
		c.emit(OpPush, encode(1)...)
		c.emit(OpEqual)
		c.emit(OpEnd)

	case "filter":
		count := c.makeConstant("count")
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		c.emit(OpPush, encode(0)...)
		c.emit(OpStore, count...)
		c.emitLoop(func() {
			c.compile(node.Arguments[1])
			c.emitCond(func() {
				c.emit(OpLoad, count...)
				c.emit(OpInc)
				c.emit(OpStore, count...)

				c.emit(OpLoad, c.makeConstant("array")...)
				c.emit(OpLoad, c.makeConstant("i")...)
				c.emit(OpIndex)
			})
		})
		c.emit(OpLoad, count...)
		c.emit(OpEnd)
		c.emit(OpArray)

	case "map":
		c.compile(node.Arguments[0])
		c.emit(OpBegin)
		size := c.emitLoop(func() {
			c.compile(node.Arguments[1])
		})
		c.emit(OpLoad, size...)
		c.emit(OpEnd)
		c.emit(OpArray)

	default:
		panic(fmt.Sprintf("unknown builtin %v", node.Name))
	}
}

func (c *compiler) emitCond(body func()) {
	noop := c.emit(OpJumpIfFalse, c.placeholder()...)
	c.emit(OpPop)

	body()

	jmp := c.emit(OpJump, c.placeholder()...)
	c.patchJump(noop)
	c.emit(OpPop)
	c.patchJump(jmp)
}

func (c *compiler) emitLoop(body func()) []byte {
	i := c.makeConstant("i")
	size := c.makeConstant("size")
	array := c.makeConstant("array")

	c.emit(OpLen)
	c.emit(OpStore, size...)
	c.emit(OpStore, array...)
	c.emit(OpPush, encode(0)...)
	c.emit(OpStore, i...)

	cond := len(c.bytecode)
	c.emit(OpLoad, i...)
	c.emit(OpLoad, size...)
	c.emit(OpLess)
	end := c.emit(OpJumpIfFalse, c.placeholder()...)
	c.emit(OpPop)

	body()

	c.emit(OpLoad, i...)
	c.emit(OpInc)
	c.emit(OpStore, i...)
	c.emit(OpJumpBackward, c.calcBackwardJump(cond)...)

	c.patchJump(end)
	c.emit(OpPop)

	return size
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

	c.emit(OpPush, encode(uint16(len(node.Nodes)))...)
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

	c.emit(OpPush, encode(uint16(len(node.Pairs)))...)
	c.emit(OpMap)
}

func encode(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

func kind(node ast.Node) reflect.Kind {
	t := node.GetType()
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}
