package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	. "expr/ast"
	"expr/conf"
	"expr/file"
	. "expr/parser/lexer"
	"expr/parser/utils"
)

type arg byte

const (
	expr arg = 1 << iota
	predicate
)

const optional arg = 1 << 7

type parser struct {
	tokens  []Token
	current Token
	pos     int
	err     *file.Error
	depth   int // predicate call depth
	config  *conf.Config
}

type Tree struct {
	Node   Node
	Source file.Source
}

func Parse(input string) (*Tree, error) {
	return ParseWithConfig(input, &conf.Config{
		Disabled: map[string]bool{},
	})
}

func ParseWithConfig(input string, config *conf.Config) (*Tree, error) {
	source := file.NewSource(input)

	tokens, err := Lex(source)
	if err != nil {
		return nil, err
	}

	p := &parser{
		tokens:  tokens,
		current: tokens[0],
		config:  config,
	}

	node := p.parseExpression()

	if !p.current.Is(EOF) {
		p.error("unexpected token %v", p.current)
	}

	tree := &Tree{
		Node:   node,
		Source: source,
	}

	if p.err != nil {
		return tree, p.err.Bind(source)
	}

	return tree, nil
}

func (p *parser) error(format string, args ...any) {
	p.errorAt(p.current, format, args...)
}

func (p *parser) errorAt(token Token, format string, args ...any) {
	if p.err == nil { // show first error
		p.err = &file.Error{
			Location: token.Location,
			Message:  fmt.Sprintf(format, args...),
		}
	}
}

func (p *parser) next() {
	p.pos++
	if p.pos >= len(p.tokens) {
		p.error("unexpected end of expression")
		return
	}
	p.current = p.tokens[p.pos]
}

func (p *parser) expect(kind Kind, values ...string) {
	if p.current.Is(kind, values...) {
		p.next()
		return
	}
	p.error("unexpected token %v", p.current)
}

// parse functions

func (p *parser) parseExpression() Node {
	nodeLeft := p.parsePrimary()
	return nodeLeft
}

func (p *parser) parsePrimary() Node {
	token := p.current

	if token.Is(Bracket, "(") {
		p.next()
		expr := p.parseExpression()
		p.expect(Bracket, ")") // "an opened parenthesis is not properly closed"
		return p.parsePostfixExpression(expr)
	}

	return p.parseSecondary()
}

func (p *parser) parseSecondary() Node {
	var node Node
	token := p.current

	switch token.Kind {

	case Identifier:
		p.next()
		switch token.Value {
		case "true":
			if p.current.Is(Bracket, "(") {
				node = p.parseCall(token, []Node{})
			} else {
				node := &BoolNode{Value: true}
				node.SetLocation(token.Location)
				return node
			}
		case "false":
			if p.current.Is(Bracket, "(") {
				node = p.parseCall(token, []Node{})
			} else {
				node := &BoolNode{Value: false}
				node.SetLocation(token.Location)
				return node
			}
		case "nil":
			node := &NilNode{}
			node.SetLocation(token.Location)
			return node
		default:
			if p.current.Is(Bracket, "(") {
				node = p.parseCall(token, []Node{})
			} else {
				node = &IdentifierNode{Value: token.Value}
				node.SetLocation(token.Location)
			}
		}

	case Number:
		p.next()
		value := strings.Replace(token.Value, "_", "", -1)
		var node Node
		valueLower := strings.ToLower(value)
		switch {
		case strings.HasPrefix(valueLower, "0x"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid hex literal: %v", err)
			}
			node = p.toIntegerNode(number)
		case strings.ContainsAny(valueLower, ".e"):
			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				p.error("invalid float literal: %v", err)
			}
			node = p.toFloatNode(number)
		case strings.HasPrefix(valueLower, "0b"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid binary literal: %v", err)
			}
			node = p.toIntegerNode(number)
		case strings.HasPrefix(valueLower, "0o"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid octal literal: %v", err)
			}
			node = p.toIntegerNode(number)
		default:
			number, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				p.error("invalid integer literal: %v", err)
			}
			node = p.toIntegerNode(number)
		}
		if node != nil {
			node.SetLocation(token.Location)
		}
		return node
	case String:
		p.next()
		node = &StringNode{Value: token.Value}
		node.SetLocation(token.Location)

	default:
		if token.Is(Bracket, "[") {
			node = p.parseArrayExpression(token)
		} else if token.Is(Bracket, "{") {
			node = p.parseMapExpression(token)
		} else {
			p.error("unexpected token %v", token)
		}
	}

	return p.parsePostfixExpression(node)
}

func (p *parser) toIntegerNode(number int64) Node {
	if number > math.MaxInt {
		p.error("integer literal is too large")
		return nil
	}
	return &IntegerNode{Value: int(number)}
}

func (p *parser) toFloatNode(number float64) Node {
	if number > math.MaxFloat64 {
		p.error("float literal is too large")
		return nil
	}
	return &FloatNode{Value: number}
}

func (p *parser) parseCall(token Token, arguments []Node) Node {
	var node Node

	callee := &IdentifierNode{Value: token.Value}
	callee.SetLocation(token.Location)
	node = &CallNode{
		Callee:    callee,
		Arguments: p.parseArguments(arguments),
	}
	node.SetLocation(token.Location)

	return node
}

func (p *parser) parseArguments(arguments []Node) []Node {
	// If pipe operator is used, the first argument is the left-hand side
	// of the operator, so we do not parse it as an argument inside brackets.
	offset := len(arguments)

	p.expect(Bracket, "(")
	for !p.current.Is(Bracket, ")") && p.err == nil {
		if len(arguments) > offset {
			p.expect(Operator, ",")
		}
		node := p.parseExpression()
		arguments = append(arguments, node)
	}
	p.expect(Bracket, ")")

	return arguments
}

func (p *parser) parseArrayExpression(token Token) Node {
	nodes := make([]Node, 0)

	p.expect(Bracket, "[")
	for !p.current.Is(Bracket, "]") && p.err == nil {
		if len(nodes) > 0 {
			p.expect(Operator, ",")
			if p.current.Is(Bracket, "]") {
				goto end
			}
		}
		node := p.parseExpression()
		nodes = append(nodes, node)
	}
end:
	p.expect(Bracket, "]")

	node := &ArrayNode{Nodes: nodes}
	node.SetLocation(token.Location)
	return node
}

func (p *parser) parseMapExpression(token Token) Node {
	p.expect(Bracket, "{")

	nodes := make([]Node, 0)
	for !p.current.Is(Bracket, "}") && p.err == nil {
		if len(nodes) > 0 {
			p.expect(Operator, ",")
			if p.current.Is(Bracket, "}") {
				goto end
			}
			if p.current.Is(Operator, ",") {
				p.error("unexpected token %v", p.current)
			}
		}

		var key Node
		// Map key can be one of:
		//  * number
		//  * string
		//  * identifier, which is equivalent to a string
		//  * expression, which must be enclosed in parentheses -- (1 + 2)
		if p.current.Is(Number) || p.current.Is(String) || p.current.Is(Identifier) {
			key = &StringNode{Value: p.current.Value}
			key.SetLocation(token.Location)
			p.next()
		} else if p.current.Is(Bracket, "(") {
			key = p.parseExpression()
		} else {
			p.error("a map key must be a quoted string, a number, a identifier, or an expression enclosed in parentheses (unexpected token %v)", p.current)
		}

		p.expect(Operator, ":")

		node := p.parseExpression()
		pair := &PairNode{Key: key, Value: node}
		pair.SetLocation(token.Location)
		nodes = append(nodes, pair)
	}

end:
	p.expect(Bracket, "}")

	node := &MapNode{Pairs: nodes}
	node.SetLocation(token.Location)
	return node
}

func (p *parser) parsePostfixExpression(node Node) Node {
	postfixToken := p.current
	for (postfixToken.Is(Operator) || postfixToken.Is(Bracket)) && p.err == nil {
		optional := postfixToken.Value == "?."
	parseToken:
		if postfixToken.Value == "." || postfixToken.Value == "?." {
			p.next()

			propertyToken := p.current
			if optional && propertyToken.Is(Bracket, "[") {
				postfixToken = propertyToken
				goto parseToken
			}
			p.next()

			if propertyToken.Kind != Identifier &&
				// Operators like "not" and "matches" are valid methods or property names.
				(propertyToken.Kind != Operator || !utils.IsValidIdentifier(propertyToken.Value)) {
				p.error("expected name")
			}

			property := &StringNode{Value: propertyToken.Value}
			property.SetLocation(propertyToken.Location)

			chainNode, isChain := node.(*ChainNode)
			optional := postfixToken.Value == "?."

			if isChain {
				node = chainNode.Node
			}

			memberNode := &MemberNode{
				Node:     node,
				Property: property,
				Optional: optional,
			}
			memberNode.SetLocation(propertyToken.Location)

			if p.current.Is(Bracket, "(") {
				memberNode.Method = true
				node = &CallNode{
					Callee:    memberNode,
					Arguments: p.parseArguments([]Node{}),
				}
				node.SetLocation(propertyToken.Location)
			} else {
				node = memberNode
			}

			if isChain || optional {
				node = &ChainNode{Node: node}
			}

		} else if postfixToken.Value == "[" {
			p.next()
			var from, to Node

			if p.current.Is(Operator, ":") { // slice without from [:1]
				p.next()

				if !p.current.Is(Bracket, "]") { // slice without from and to [:]
					to = p.parseExpression()
				}

				node = &SliceNode{
					Node: node,
					To:   to,
				}
				node.SetLocation(postfixToken.Location)
				p.expect(Bracket, "]")

			} else {

				from = p.parseExpression()

				if p.current.Is(Operator, ":") {
					p.next()

					if !p.current.Is(Bracket, "]") { // slice without to [1:]
						to = p.parseExpression()
					}

					node = &SliceNode{
						Node: node,
						From: from,
						To:   to,
					}
					node.SetLocation(postfixToken.Location)
					p.expect(Bracket, "]")

				} else {
					// Slice operator [:] was not found,
					// it should be just an index node.
					node = &MemberNode{
						Node:     node,
						Property: from,
						Optional: optional,
					}
					node.SetLocation(postfixToken.Location)
					if optional {
						node = &ChainNode{Node: node}
					}
					p.expect(Bracket, "]")
				}
			}
		} else {
			break
		}
		postfixToken = p.current
	}
	return node
}
