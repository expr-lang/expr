package compiler

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/file"
)

func (c *compiler) countCommonExpr(node ast.Node) {
	switch n := node.(type) {
	case *ast.NilNode:
		c.commonCommonNilNode(n)
	case *ast.IdentifierNode:
		c.countCommonIdentifierNode(n)
	case *ast.IntegerNode:
		c.countCommonIntegerNode(n)
	case *ast.FloatNode:
		c.countCommonFloatNode(n)
	case *ast.BoolNode:
		c.countCommonBoolNode(n)
	case *ast.StringNode:
		c.countCommonStringNode(n)
	case *ast.ConstantNode:
		c.countCommonConstantNode(n)
	case *ast.UnaryNode:
		c.countCommonUnaryNode(n)
	case *ast.BinaryNode:
		c.countCommonBinaryNode(n)
	case *ast.ChainNode:
		c.countCommonChainNode(n)
	case *ast.MemberNode:
		c.countCommonMemberNode(n)
	case *ast.SliceNode:
		c.countCommonSliceNode(n)
	case *ast.CallNode:
		c.checkCallNode(n)
	case *ast.BuiltinNode:
		c.countCommonBuiltinNode(n)
	case *ast.ClosureNode:
		c.countCommonClosureNode(n)
	case *ast.PointerNode:
		c.countCommonPointerNode(n)
	case *ast.ConditionalNode:
		c.countCommonConditionalNode(n)
	case *ast.ArrayNode:
		c.countCommonArrayNode(n)
	case *ast.MapNode:
		c.countCommonMapNode(n)
	case *ast.PairNode:
		// do nothing
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}

func (c *compiler) commonCommonNilNode(n *ast.NilNode) {
	n.SetSubExpr("nil")
}

func (c *compiler) countCommonIdentifierNode(n *ast.IdentifierNode) {
	n.SetSubExpr(n.Value)
}

func (c *compiler) countCommonIntegerNode(n *ast.IntegerNode) {
	n.SetSubExpr(strconv.FormatInt(int64(n.Value), 10))
}

func (c *compiler) countCommonFloatNode(n *ast.FloatNode) {
	n.SetSubExpr(strconv.FormatFloat(n.Value, 'f', 10, 64))
}

func (c *compiler) countCommonBoolNode(n *ast.BoolNode) {
	n.SetSubExpr(strconv.FormatBool(n.Value))
}

func (c *compiler) countCommonStringNode(n *ast.StringNode) {
	n.SetSubExpr(strconv.Quote(n.Value))
}

func (c *compiler) countCommonConstantNode(n *ast.ConstantNode) {
	var s string
	switch n.Value.(type) {
	case string:
		s = strconv.Quote(n.Value.(string))
	case int:
		s = strconv.FormatInt(int64(n.Value.(int)), 10)
	case int8:
		s = strconv.FormatInt(int64(n.Value.(int8)), 10)
	case int16:
		s = strconv.FormatInt(int64(n.Value.(int16)), 10)
	case int32:
		s = strconv.FormatInt(int64(n.Value.(int32)), 10)
	case int64:
		s = strconv.FormatInt(n.Value.(int64), 10)
	case uint:
		s = strconv.FormatUint(uint64(n.Value.(uint)), 10)
	case uint8:
		s = strconv.FormatUint(uint64(n.Value.(uint8)), 10)
	case uint16:
		s = strconv.FormatUint(uint64(n.Value.(uint16)), 10)
	case uint32:
		s = strconv.FormatUint(uint64(n.Value.(uint32)), 10)
	case uint64:
		s = strconv.FormatUint(uint64(n.Value.(uint64)), 10)
	case float32:
		s = strconv.FormatFloat(float64(n.Value.(float32)), 'f', 10, 32)
	case float64:
		s = strconv.FormatFloat(n.Value.(float64), 'f', 10, 64)
	default:
		s = fmt.Sprintf("%+v", n.Value)
	}
	n.SetSubExpr(s)
}

func (c *compiler) countCommonUnaryNode(n *ast.UnaryNode) {
	c.countCommonExpr(n.Node)
	buf := strings.Builder{}
	switch n.Operator {
	case "+":
		buf.WriteString("")
	case "-":
		buf.WriteString("-(")
	case "not", "!":
		buf.WriteString("not (")
	}
	buf.WriteString(n.Node.SubExpr())
	buf.WriteString(")")
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonBinaryNode(n *ast.BinaryNode) {
	c.countCommonExpr(n.Left)
	c.countCommonExpr(n.Right)

	ls := n.Left.SubExpr()
	rs := n.Right.SubExpr()
	_, lw := n.Left.(*ast.BinaryNode)
	_, rw := n.Right.(*ast.BinaryNode)
	op := n.Operator
	switch op {
	case "==", "!=", "and", "or", "+", "*", "||", "&&", ">=", "<=": // right / left can be swap
		if op == ">=" || op == "<=" || rs <= ls {
			ls, rs = rs, ls
			lw, rw = rw, lw
		}
		if op == "&&" {
			op = "and"
		} else if op == "||" {
			op = "or"
		} else if op == ">=" {
			op = "<"
		} else if op == "<=" {
			op = ">"
		}
	case "**", "^":
		op = "**"
	default:
		// do nothing
	}

	buf := strings.Builder{}
	if lw {
		buf.WriteString("(")
		buf.WriteString(ls)
		buf.WriteString(")")
	} else {
		buf.WriteString(ls)
	}
	buf.WriteString(" " + op + " ")
	if rw {
		buf.WriteString("(")
		buf.WriteString(rs)
		buf.WriteString(")")
	} else {
		buf.WriteString(rs)
	}
	n.SetSubExpr(buf.String())

	switch n.Operator {
	case "??", "and", "or", "||", "&&":
		// do nothing
	default:
		c.emitSubExpr(n.SubExpr(), n.Location())
	}
}

func (c *compiler) countCommonChainNode(n *ast.ChainNode) {
	c.countCommonExpr(n.Node)
	n.SetSubExpr(n.Node.SubExpr())
}

func (c *compiler) countCommonMemberNode(n *ast.MemberNode) {
	c.countCommonExpr(n.Node)
	buf := strings.Builder{}
	buf.WriteString(n.Node.SubExpr())
	if n.Method {
		buf.WriteString(".")
		buf.WriteString(n.Name)
	} else {
		if len(n.FieldIndex) > 0 {
			if n.Optional {
				buf.WriteString("?")
			}
			buf.WriteString(".")
			buf.WriteString(n.Name)
		} else {
			c.countCommonExpr(n.Property)
			buf.WriteString("[")
			buf.WriteString(n.Property.SubExpr())
			buf.WriteString("]")
		}
	}
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonSliceNode(n *ast.SliceNode) {
	buf := strings.Builder{}
	c.countCommonExpr(n.Node)
	buf.WriteString(n.Node.SubExpr())
	buf.WriteString("[")
	toStr := ""
	fromStr := ""
	if n.To != nil {
		c.countCommonExpr(n.To)
		toStr = n.To.SubExpr()
	}
	if n.From != nil {
		c.countCommonExpr(n.From)
		fromStr = n.From.SubExpr()
	}
	buf.WriteString(fromStr)
	buf.WriteString(":")
	buf.WriteString(toStr)
	buf.WriteString("]")
	n.SetSubExpr(buf.String())
}

func (c *compiler) checkCallNode(n *ast.CallNode) {
	buf := strings.Builder{}
	c.countCommonExpr(n.Callee)
	buf.WriteString(n.Callee.SubExpr())
	buf.WriteString("(")
	for i, arg := range n.Arguments {
		c.countCommonExpr(arg)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.SubExpr())
	}
	buf.WriteString(")")
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonBuiltinNode(n *ast.BuiltinNode) {
	buf := strings.Builder{}
	buf.WriteString(n.Name)
	buf.WriteString("(")
	for i, arg := range n.Arguments {
		c.countCommonExpr(arg)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.SubExpr())
	}
	buf.WriteString(")")
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonClosureNode(n *ast.ClosureNode) {
	c.countCommonExpr(n.Node)
	// omit '{' / '}'
	n.SetSubExpr(n.Node.SubExpr())
}

func (c *compiler) countCommonPointerNode(n *ast.PointerNode) {
	// do nothing
	n.SetSubExpr("#")
}

func (c *compiler) countCommonConditionalNode(n *ast.ConditionalNode) {
	c.countCommonExpr(n.Cond)
	c.countCommonExpr(n.Exp1)
	c.countCommonExpr(n.Exp2)
	buf := strings.Builder{}
	buf.WriteString(n.Cond.SubExpr())
	buf.WriteString(" ? ")
	buf.WriteString(n.Exp1.SubExpr())
	buf.WriteString(" : ")
	buf.WriteString(n.Exp2.SubExpr())
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonArrayNode(n *ast.ArrayNode) {
	buf := strings.Builder{}
	buf.WriteString("[")
	for i, node := range n.Nodes {
		c.countCommonExpr(node)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(node.SubExpr())
	}
	buf.WriteString("]")
	n.SetSubExpr(buf.String())
}

func (c *compiler) countCommonMapNode(n *ast.MapNode) {
	pairs := make([]*ast.PairNode, 0)
	for i := range n.Pairs {
		pair := n.Pairs[i].(*ast.PairNode)
		c.countCommonExpr(pair.Key)
		c.countCommonExpr(pair.Value)
		pairs = append(pairs, pair)
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key.SubExpr() < pairs[j].Key.SubExpr()
	})
	buf := strings.Builder{}
	buf.WriteString("{")
	for i, pair := range pairs {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(pair.Key.SubExpr())
		buf.WriteString(":")
		buf.WriteString(pair.Value.SubExpr())
	}
	buf.WriteString("}")
	n.SetSubExpr(buf.String())
}

func (c *compiler) emitSubExpr(subExpr string, loc file.Location) {
	if c.exprRecords == nil || subExpr == "" {
		return
	}
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(subExpr)))
	if cs, ok := c.exprRecords[hash]; !ok {
		c.exprRecords[hash] = &exprRecord{cnt: 1, id: -1}
	} else {
		cs.cnt = cs.cnt + 1
	}
}

func (c *compiler) needCacheCommon(n ast.Node) (bool, int) {
	needCacheCommon, exprUniqId := false, -1
	if c.exprRecords != nil {
		expr := n.SubExpr()
		hash := fmt.Sprintf("%x", sha1.Sum([]byte(expr)))
		cs, ok := c.exprRecords[hash]
		if ok && cs.cnt > 1 {
			if cs.id == -1 {
				cs.id = c.commonExprInc
				c.commonExpr[cs.id] = expr
				c.commonExprInc += 1
			}
			needCacheCommon = true
			exprUniqId = cs.id
		}
	}
	return needCacheCommon, exprUniqId
}
