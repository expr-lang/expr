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

func (c *compiler) analyzeCommonExpr(node ast.Node) {
	switch n := node.(type) {
	case *ast.NilNode:
		c.analyzeCommonNilNode(n)
	case *ast.IdentifierNode:
		c.analyzeCommonIdentifierNode(n)
	case *ast.IntegerNode:
		c.analyzeCommonIntegerNode(n)
	case *ast.FloatNode:
		c.analyzeCommonFloatNode(n)
	case *ast.BoolNode:
		c.analyzeCommonBoolNode(n)
	case *ast.StringNode:
		c.analyzeCommonStringNode(n)
	case *ast.ConstantNode:
		c.analyzeCommonConstantNode(n)
	case *ast.UnaryNode:
		c.analyzeCommonUnaryNode(n)
	case *ast.BinaryNode:
		c.analyzeCommonBinaryNode(n)
	case *ast.ChainNode:
		c.analyzeCommonChainNode(n)
	case *ast.MemberNode:
		c.analyzeCommonMemberNode(n)
	case *ast.SliceNode:
		c.analyzeCommonSliceNode(n)
	case *ast.CallNode:
		c.checkCallNode(n)
	case *ast.BuiltinNode:
		c.analyzeCommonBuiltinNode(n)
	case *ast.ClosureNode:
		c.analyzeCommonClosureNode(n)
	case *ast.PointerNode:
		c.analyzeCommonPointerNode(n)
	case *ast.ConditionalNode:
		c.analyzeCommonConditionalNode(n)
	case *ast.ArrayNode:
		c.analyzeCommonArrayNode(n)
	case *ast.MapNode:
		c.analyzeCommonMapNode(n)
	case *ast.PairNode:
		// do nothing
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}

func (c *compiler) analyzeCommonNilNode(n *ast.NilNode) {
	n.SetSubExpr("nil")
}

func (c *compiler) analyzeCommonIdentifierNode(n *ast.IdentifierNode) {
	n.SetSubExpr(n.Value)
}

func (c *compiler) analyzeCommonIntegerNode(n *ast.IntegerNode) {
	n.SetSubExpr(strconv.FormatInt(int64(n.Value), 10))
}

func (c *compiler) analyzeCommonFloatNode(n *ast.FloatNode) {
	n.SetSubExpr(strconv.FormatFloat(n.Value, 'f', 10, 64))
}

func (c *compiler) analyzeCommonBoolNode(n *ast.BoolNode) {
	n.SetSubExpr(strconv.FormatBool(n.Value))
}

func (c *compiler) analyzeCommonStringNode(n *ast.StringNode) {
	n.SetSubExpr(strconv.Quote(n.Value))
}

func (c *compiler) analyzeCommonConstantNode(n *ast.ConstantNode) {
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

func (c *compiler) analyzeCommonUnaryNode(n *ast.UnaryNode) {
	c.analyzeCommonExpr(n.Node)
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

func (c *compiler) analyzeCommonBinaryNode(n *ast.BinaryNode) {
	c.analyzeCommonExpr(n.Left)
	c.analyzeCommonExpr(n.Right)

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
		c.countCommonExpr(n.SubExpr(), n.Location())
	}
}

func (c *compiler) analyzeCommonChainNode(n *ast.ChainNode) {
	c.analyzeCommonExpr(n.Node)
	n.SetSubExpr(n.Node.SubExpr())
}

func (c *compiler) analyzeCommonMemberNode(n *ast.MemberNode) {
	c.analyzeCommonExpr(n.Node)
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
			c.analyzeCommonExpr(n.Property)
			buf.WriteString("[")
			buf.WriteString(n.Property.SubExpr())
			buf.WriteString("]")
		}
	}
	n.SetSubExpr(buf.String())
}

func (c *compiler) analyzeCommonSliceNode(n *ast.SliceNode) {
	buf := strings.Builder{}
	c.analyzeCommonExpr(n.Node)
	buf.WriteString(n.Node.SubExpr())
	buf.WriteString("[")
	toStr := ""
	fromStr := ""
	if n.To != nil {
		c.analyzeCommonExpr(n.To)
		toStr = n.To.SubExpr()
	}
	if n.From != nil {
		c.analyzeCommonExpr(n.From)
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
	c.analyzeCommonExpr(n.Callee)
	buf.WriteString(n.Callee.SubExpr())
	buf.WriteString("(")
	for i, arg := range n.Arguments {
		c.analyzeCommonExpr(arg)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.SubExpr())
	}
	buf.WriteString(")")
	n.SetSubExpr(buf.String())
	c.countCommonExpr(n.SubExpr(), n.Location())
}

func (c *compiler) analyzeCommonBuiltinNode(n *ast.BuiltinNode) {
	buf := strings.Builder{}
	buf.WriteString(n.Name)
	buf.WriteString("(")
	for i, arg := range n.Arguments {
		c.analyzeCommonExpr(arg)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.SubExpr())
	}
	buf.WriteString(")")
	n.SetSubExpr(buf.String())
}

func (c *compiler) analyzeCommonClosureNode(n *ast.ClosureNode) {
	c.analyzeCommonExpr(n.Node)
	// omit '{' / '}'
	n.SetSubExpr(n.Node.SubExpr())
}

func (c *compiler) analyzeCommonPointerNode(n *ast.PointerNode) {
	// do nothing
	n.SetSubExpr("#")
}

func (c *compiler) analyzeCommonConditionalNode(n *ast.ConditionalNode) {
	c.analyzeCommonExpr(n.Cond)
	c.analyzeCommonExpr(n.Exp1)
	c.analyzeCommonExpr(n.Exp2)
	buf := strings.Builder{}
	buf.WriteString(n.Cond.SubExpr())
	buf.WriteString(" ? ")
	buf.WriteString(n.Exp1.SubExpr())
	buf.WriteString(" : ")
	buf.WriteString(n.Exp2.SubExpr())
	n.SetSubExpr(buf.String())
}

func (c *compiler) analyzeCommonArrayNode(n *ast.ArrayNode) {
	buf := strings.Builder{}
	buf.WriteString("[")
	for i, node := range n.Nodes {
		c.analyzeCommonExpr(node)
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(node.SubExpr())
	}
	buf.WriteString("]")
	n.SetSubExpr(buf.String())
}

func (c *compiler) analyzeCommonMapNode(n *ast.MapNode) {
	pairs := make([]*ast.PairNode, 0)
	for i := range n.Pairs {
		pair := n.Pairs[i].(*ast.PairNode)
		c.analyzeCommonExpr(pair.Key)
		c.analyzeCommonExpr(pair.Value)
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

func (c *compiler) countCommonExpr(subExpr string, loc file.Location) {
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

func (c *compiler) needReuseCommon(n ast.Node) (bool, bool, int) {
	needReuseCommon, isFirstOccur, exprUniqId := false, false, -1
	if c.exprRecords != nil {
		expr := n.SubExpr()
		hash := fmt.Sprintf("%x", sha1.Sum([]byte(expr)))
		cs, ok := c.exprRecords[hash]
		if ok && cs.cnt > 1 {
			if cs.id == -1 {
				cs.id = c.commonExprInc
				cs.loc = n.Location()
				c.commonExpr[cs.id] = expr
				c.commonExprInc += 1
			}
			needReuseCommon = true
			isFirstOccur = n.Location().Line == cs.loc.Line && n.Location().Column == cs.loc.Column
			exprUniqId = cs.id
		}
	}
	return needReuseCommon, isFirstOccur, exprUniqId
}
