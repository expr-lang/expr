package parser

//go:generate antlr -Dlanguage=Go -listener -visitor -o gen -package gen Expr.g4

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/internal/file"
	"github.com/antonmedv/expr/parser/gen"
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

func (p *parser) EnterIdentifier(ctx *gen.IdentifierContext) {
	p.push(&ast.IdentifierNode{Value: ctx.GetText()}).SetLocation(location(ctx))
}

func (p *parser) EnterPointer(ctx *gen.PointerContext) {
	p.push(&ast.PointerNode{}).SetLocation(location(ctx))
}

func (p *parser) EnterString(ctx *gen.StringContext) {
	var value string
	if s, err := unescape(ctx.GetText()); err == nil {
		value = s
	} else {
		p.reportError(ctx, "parse error: %v", err)
		return
	}
	node := &ast.StringNode{
		Value: value,
	}
	p.push(node).SetLocation(location(ctx))
}

func (p *parser) EnterInteger(ctx *gen.IntegerContext) {
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

func (p *parser) EnterFloat(ctx *gen.FloatContext) {
	f, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		p.reportError(ctx, "parse error: invalid float literal")
		return
	}
	p.push(&ast.FloatNode{Value: f}).SetLocation(location(ctx))
}

func (p *parser) EnterBoolean(ctx *gen.BooleanContext) {
	b, err := strconv.ParseBool(ctx.GetText())
	if err != nil {
		p.reportError(ctx, "parse error: invalid boolean literal")
		return
	}
	p.push(&ast.BoolNode{Value: b}).SetLocation(location(ctx))
}

func (p *parser) EnterNil(ctx *gen.NilContext) {
	p.push(&ast.NilNode{}).SetLocation(location(ctx))
}

func (p *parser) ExitUnary(ctx *gen.UnaryContext) {
	p.push(&ast.UnaryNode{
		Operator: ctx.GetOp().GetText(),
		Node:     p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitBinary(ctx *gen.BinaryContext) {
	p.push(&ast.BinaryNode{
		Operator: ctx.GetOp().GetText(),
		Right:    p.pop(ctx),
		Left:     p.pop(ctx),
	}).SetLocation(locationToken(ctx.GetOp()))
}

func (p *parser) ExitMatches(ctx *gen.MatchesContext) {
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

func (p *parser) ExitCall(ctx *gen.CallContext) {
	expr := ctx.GetChild(0)
	args := ctx.GetArgs()

	var list []gen.IExprContext
	if args != nil {
		list = args.GetList()
	}
	arguments := p.arguments(ctx, list)

	switch c := expr.(type) {
	case *gen.IdentifierContext:
		p.push(&ast.FunctionNode{
			Arguments: arguments,
			Name:      p.pop(ctx).(*ast.IdentifierNode).Value,
		}).SetLocation(location(ctx))
	case *gen.MemberDotContext:
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

func (p *parser) ExitMemberIndex(ctx *gen.MemberIndexContext) {
	p.push(&ast.IndexNode{
		Index: p.pop(ctx),
		Node:  p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitSlice(ctx *gen.SliceContext) {
	var a, b ast.Node
	if ctx.GetB() != nil {
		b = p.pop(ctx)
	}
	if ctx.GetA() != nil {
		a = p.pop(ctx)
	}
	node := p.pop(ctx)
	p.push(&ast.SliceNode{
		Node: node,
		From: a,
		To:   b,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitMemberDot(ctx *gen.MemberDotContext) {
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

func (p *parser) ExitTernary(ctx *gen.TernaryContext) {
	expr2 := p.pop(ctx)
	expr1 := p.pop(ctx)
	cond := p.pop(ctx)
	p.push(&ast.ConditionalNode{
		Exp2: expr2,
		Exp1: expr1,
		Cond: cond,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitArrayLiteral(ctx *gen.ArrayLiteralContext) {
	list := ctx.GetList()
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

	nodes := make([]ast.Node, 0)
	for range e.GetList() {
		nodes = append([]ast.Node{p.pop(ctx).(*ast.PairNode)}, nodes...)
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
		s2, err := unescape(str.GetText())
		if err != nil {
			p.reportError(ctx, "parse error: %v", err)
			return
		}
		s = s2
	} else {
		p.reportError(ctx, "parse error: invalid key type")
		return
	}

	key := &ast.StringNode{Value: s}
	key.SetLocation(location(ctx))

	p.push(&ast.PairNode{
		Key:   key,
		Value: value,
	}).SetLocation(location(ctx))
}

func (p *parser) ExitBuiltinLen(ctx *gen.BuiltinLenContext) {
	name := ctx.GetName().GetText()
	node := p.pop(ctx.GetE())

	p.push(&ast.BuiltinNode{
		Name:      name,
		Arguments: []ast.Node{node},
	}).SetLocation(location(ctx))
}

func (p *parser) ExitBuiltin(ctx *gen.BuiltinContext) {
	name := ctx.GetName().GetText()
	closure := p.pop(ctx.GetC())
	node := p.pop(ctx.GetE())

	p.push(&ast.BuiltinNode{
		Name:      name,
		Arguments: []ast.Node{node, closure},
	}).SetLocation(location(ctx))
}

func (p *parser) EnterClosure(ctx *gen.ClosureContext) {
	p.closure = true
}

func (p *parser) ExitClosure(ctx *gen.ClosureContext) {
	p.closure = false
	p.push(&ast.ClosureNode{
		Node: p.pop(ctx),
	}).SetLocation(location(ctx))
}

func (p *parser) ExitClosureMemberDot(ctx *gen.ClosureMemberDotContext) {
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
