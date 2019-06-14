package parser

//go:generate antlr -Dlanguage=Go -listener -visitor -o gen -package gen Expr.g4

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"gopkg.in/antonmedv/expr.v2/ast"
	"gopkg.in/antonmedv/expr.v2/internal/file"
	"gopkg.in/antonmedv/expr.v2/parser/gen"
	"regexp"
	"strconv"
	"strings"
)

type Tree struct {
	Node   ast.Node
	Source *file.Source
}

func Parse(input string) (*Tree, error) {
	source := file.NewSource(input)
	is := antlr.NewInputStream(input)

	lexer := gen.NewExprLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	expr := gen.NewExprParser(stream)

	p := &parser{
		errors: file.NewErrors(source),
	}

	lexer.RemoveErrorListeners()
	expr.RemoveErrorListeners()
	lexer.AddErrorListener(p)
	expr.AddErrorListener(p)

	antlr.ParseTreeWalkerDefault.Walk(p, expr.Start())

	if p.errors.HasError() {
		return nil, fmt.Errorf("%v", p.errors.First())
	}
	if len(p.stack) == 0 {
		return nil, fmt.Errorf("empty stack")
	}
	if len(p.stack) > 1 {
		return nil, fmt.Errorf("too long stack")
	}
	return &Tree{
		Node:   p.stack[0],
		Source: source,
	}, nil
}

type parser struct {
	*gen.BaseExprListener
	stack   []ast.Node
	errors  *file.Errors
	closure bool
}

func (p *parser) push(node ast.Node) ast.Node {
	p.stack = append(p.stack, node)
	return node
}

func (p *parser) pop(ctx antlr.ParserRuleContext, it ...string) ast.Node {
	if len(p.stack) == 0 {
		if len(it) == 0 {
			p.reportError(ctx, "parse error: the expression lacks something")
		} else {
			p.reportError(ctx, "parse error: the expression lacks "+strings.Join(it, ", "))
		}
		return &ast.NilNode{}
	}
	node := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return node
}

func (p *parser) reportError(ctx antlr.ParserRuleContext, format string, args ...interface{}) {
	p.errors.ReportError(location(ctx), format, args...)
}

func (p *parser) EnterIdentifierExpression(ctx *gen.IdentifierExpressionContext) {
	p.push(&ast.IdentifierNode{Value: ctx.GetText()}).SetLocation(location(ctx))
}

func (p *parser) EnterPointerExpression(ctx *gen.PointerExpressionContext) {
	p.push(&ast.PointerNode{}).SetLocation(location(ctx))
}

func (p *parser) EnterStringLiteral(ctx *gen.StringLiteralContext) {
	p.push(&ast.StringNode{Value: unquotes(ctx.GetText())}).SetLocation(location(ctx))
}

func (p *parser) EnterIntegerLiteral(ctx *gen.IntegerLiteralContext) {
	if node := ctx.IntegerLiteral(); node != nil {
		text := node.GetText()
		text = strings.Replace(text, "_", "", -1)
		i, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			p.reportError(ctx, "parse error: invalid int literal")
			return
		}
		p.push(&ast.IntegerNode{Value: int(i)}).SetLocation(location(ctx))
	} else if node := ctx.HexIntegerLiteral(); node != nil {
		text := node.GetText()
		i, err := strconv.ParseInt(text, 0, 64)
		if err != nil {
			p.reportError(ctx, "parse error: invalid hex literal")
			return
		}
		p.push(&ast.IntegerNode{Value: int(i)}).SetLocation(location(ctx))
	} else {
		p.reportError(ctx, "parse error: invalid octal literal")
	}
}

func (p *parser) EnterFloatExpression(ctx *gen.FloatExpressionContext) {
	f, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		p.reportError(ctx, "parse error: invalid float literal")
		return
	}
	p.push(&ast.FloatNode{Value: f}).SetLocation(location(ctx))
}

func (p *parser) EnterBooleanExpression(ctx *gen.BooleanExpressionContext) {
	b, err := strconv.ParseBool(ctx.GetText())
	if err != nil {
		p.reportError(ctx, "parse error: invalid boolean literal")
		return
	}
	p.push(&ast.BoolNode{Value: b}).SetLocation(location(ctx))
}

func (p *parser) EnterNilExpression(ctx *gen.NilExpressionContext) {
	p.push(&ast.NilNode{}).SetLocation(location(ctx))
}

func (p *parser) ExitUnaryExpression(ctx *gen.UnaryExpressionContext) {
	p.push(&ast.UnaryNode{
		Operator: ctx.GetOp().GetText(),
		Node:     p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitRangeExpression(ctx *gen.RangeExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: "..",
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitMultiplicativeExpression(ctx *gen.MultiplicativeExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitAdditiveExpression(ctx *gen.AdditiveExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitRelationalExpression(ctx *gen.RelationalExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitStartsWithExpression(ctx *gen.StartsWithExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: "startsWith",
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitEndsWithExpression(ctx *gen.EndsWithExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: "endsWith",
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitContainsExpression(ctx *gen.ContainsExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: "contains",
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitMatchesExpression(ctx *gen.MatchesExpressionContext) {
	right := p.pop(ctx)
	left := p.pop(ctx)
	node := &ast.MatchesNode{
		Right: right,
		Left:  left,
	}

	var err error
	var r *regexp.Regexp
	if s, ok := right.(*ast.StringNode); ok {
		r, err = regexp.Compile(s.Value)
		if err != nil {
			p.reportError(ctx.GetPattern(), "%v", err)
			return
		}
		node.Regexp = r
	}
	p.push(node).SetLocation(location(ctx))
}

func (p *parser) ExitInExpression(ctx *gen.InExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitEqualityExpression(ctx *gen.EqualityExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitLogicalExpression(ctx *gen.LogicalExpressionContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitCallExpression(ctx *gen.CallExpressionContext) {
	expr := ctx.GetChild(0)
	args := ctx.GetArgs()

	var list []gen.IExprContext
	if args != nil {
		list = args.GetList()
	}
	arguments := p.arguments(ctx, list)

	switch c := expr.(type) {
	case *gen.IdentifierExpressionContext:
		p.push(&ast.FunctionNode{
			Arguments: arguments,
			Name:      p.pop(ctx).(*ast.IdentifierNode).Value,
		}).SetLocation(location(ctx))
	case *gen.MemberDotExpressionContext:
		p.push(&ast.MethodNode{
			Arguments: arguments,
			Method:    c.GetName().GetText(),
			Node:      p.pop(ctx).(*ast.PropertyNode).Node,
		}).SetLocation(location(ctx))
	default:
		p.reportError(ctx, "parse error: undefined call expression")
	}
}

func (p *parser) arguments(ctx antlr.ParserRuleContext, list []gen.IExprContext) []ast.Node {
	args := make([]ast.Node, 0)
	for range list {
		args = append([]ast.Node{p.pop(ctx)}, args...)
	}
	return args
}

func (p *parser) ExitMemberIndexExpression(ctx *gen.MemberIndexExpressionContext) {
	p.push(&ast.IndexNode{
		Index: p.pop(ctx),
		Node:  p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitMemberDotExpression(ctx *gen.MemberDotExpressionContext) {
	var property string
	name := ctx.GetName()
	if name != nil {
		property = name.GetText()
	}
	p.push(&ast.PropertyNode{
		Node:     p.pop(ctx),
		Property: property,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitTernaryExpression(ctx *gen.TernaryExpressionContext) {
	expr2 := p.pop(ctx)
	expr1 := p.pop(ctx)
	cond := p.pop(ctx)
	p.push(&ast.ConditionalNode{
		Exp2: expr2,
		Exp1: expr1,
		Cond: cond,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitArrayLiteralExpression(ctx *gen.ArrayLiteralExpressionContext) {
	list := ctx.GetChild(0).(*gen.ArrayLiteralContext).GetList()
	nodes := make([]ast.Node, 0)
	for range list {
		nodes = append([]ast.Node{p.pop(ctx)}, nodes...)
	}
	p.push(&ast.ArrayNode{Nodes: nodes}).SetLocation(location(ctx))
}

func (p *parser) ExitMapLiteral(ctx *gen.MapLiteralContext) {
	e := ctx.GetE()
	if e == nil {
		p.push(&ast.MapNode{}).SetLocation(location(ctx))
		return
	}

	nodes := make([]*ast.PairNode, 0)
	for range e.GetList() {
		nodes = append([]*ast.PairNode{p.pop(ctx).(*ast.PairNode)}, nodes...)
	}
	p.push(&ast.MapNode{Pairs: nodes}).SetLocation(location(ctx))
}

func (p *parser) ExitPropertyAssignment(ctx *gen.PropertyAssignmentContext) {
	value := p.pop(ctx)
	name := ctx.GetName().(*gen.PropertyNameContext)

	var s string
	if id := name.Identifier(); id != nil {
		s = id.GetText()
	} else if str := name.StringLiteral(); str != nil {
		s = unquotes(str.GetText())
	} else {
		p.reportError(ctx, "parse error: invalid key type")
	}

	key := &ast.StringNode{Value: s}
	key.SetLocation(location(ctx))

	p.push(&ast.PairNode{
		Key:   key,
		Value: value,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitLenBuiltinExpression(ctx *gen.LenBuiltinExpressionContext) {
	p.push(&ast.BuiltinNode{
		Name: "len",
		Arguments: []ast.Node{
			p.pop(ctx.GetE(), "parameter"),
		},
	}).SetLocation(location(ctx))
}

func (p *parser) ExitBuiltinExpression(ctx *gen.BuiltinExpressionContext) {
	name := ctx.GetName().GetText()
	closure := p.pop(ctx.GetC())
	node := p.pop(ctx.GetE())

	p.push(&ast.BuiltinNode{
		Name:      name,
		Arguments: []ast.Node{node, closure},
	}).SetLocation(location(ctx))
}

func (p *parser) EnterClosureExpression(ctx *gen.ClosureExpressionContext) {
	p.closure = true
}

func (p *parser) ExitClosureExpression(ctx *gen.ClosureExpressionContext) {
	p.closure = false
	p.push(&ast.ClosureNode{
		Node: p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitClosureMemberDotExpression(ctx *gen.ClosureMemberDotExpressionContext) {
	if !p.closure {
		p.reportError(ctx, "parse error: dot property accessor can be only inside closure")
		return
	}
	var property string
	name := ctx.GetName()
	if name != nil {
		property = name.GetText()
	}
	pointer := &ast.PointerNode{}
	pointer.SetLocation(location(ctx))
	p.push(&ast.PropertyNode{
		Node:     pointer,
		Property: property,
	}).SetLocation(location(ctx))
}

func (p *parser) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	p.errors.ReportError(file.Location{Line: line, Column: column}, fmt.Sprintf("syntax error: %s", msg))
}

func (p *parser) ReportAmbiguity(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ bool, _ *antlr.BitSet, _ antlr.ATNConfigSet) {
}

func (p *parser) ReportAttemptingFullContext(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ *antlr.BitSet, _ antlr.ATNConfigSet) {
}

func (p *parser) ReportContextSensitivity(_ antlr.Parser, _ *antlr.DFA, _, _, _ int, _ antlr.ATNConfigSet) {
}

func unquotes(s string) string {
	if len(s) >= 2 {
		s = strings.Replace(s, string([]byte{'\\', s[0]}), string(s[0]), -1)
		s = strings.Replace(s, `\\`, `\`, -1)
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func location(ctx antlr.ParserRuleContext) file.Location {
	if ctx == nil {
		return file.Location{Line: 0, Column: 0}
	}

	token := ctx.GetStart()
	if token == nil {
		return file.Location{Line: 0, Column: 0}
	}

	return file.Location{Line: token.GetLine(), Column: token.GetColumn()}
}

func locationToken(token antlr.Token) file.Location {
	return file.Location{Line: token.GetLine(), Column: token.GetColumn()}
}
