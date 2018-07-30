package expr

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type associativity int

const (
	left associativity = iota + 1
	right
)

type info struct {
	precedence    int
	associativity associativity
}

var unaryOperators = map[string]info{
	"not": {50, left},
	"!":   {50, left},
	"-":   {500, left},
	"+":   {500, left},
}

var binaryOperators = map[string]info{
	"or":      {10, left},
	"||":      {10, left},
	"and":     {15, left},
	"&&":      {15, left},
	"|":       {16, left},
	"^":       {17, left},
	"&":       {18, left},
	"==":      {20, left},
	"!=":      {20, left},
	"<":       {20, left},
	">":       {20, left},
	">=":      {20, left},
	"<=":      {20, left},
	"not in":  {20, left},
	"in":      {20, left},
	"matches": {20, left},
	"..":      {25, left},
	"+":       {30, left},
	"-":       {30, left},
	"~":       {40, left},
	"*":       {60, left},
	"/":       {60, left},
	"%":       {60, left},
	"**":      {200, right},
}

var builtins = map[string]bool{
	"len": true,
}

type parser struct {
	input    string
	tokens   []token
	position int
	current  token
	options  *options
}

type options struct {
	names map[string]struct{}
	funcs map[string]struct{}
}

// OptionFn for configuring parser.
type OptionFn func(p *options)

// Parse parses input into ast.
func Parse(input string, ops ...OptionFn) (Node, error) {
	tokens, err := lex(input)
	if err != nil {
		return nil, err
	}

	p := &parser{
		input:   input,
		tokens:  tokens,
		current: tokens[0],
		options: &options{},
	}

	for _, op := range ops {
		op(p.options)
	}

	node, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	if !p.isEOF() {
		return nil, p.errorf("unexpected token %v", p.current)
	}
	return node, nil
}

// Names sets list of allowed names.
func Names(names ...string) OptionFn {
	return func(o *options) {
		set := make(map[string]struct{})
		for _, name := range names {
			set[name] = struct{}{}
		}
		o.names = set
	}
}

// Funcs sets list of allowed function.
func Funcs(funcs ...string) OptionFn {
	return func(o *options) {
		set := make(map[string]struct{})
		for _, fn := range funcs {
			set[fn] = struct{}{}
		}
		o.funcs = set
	}
}

func (p *parser) errorf(format string, args ...interface{}) error {
	return &syntaxError{
		message: fmt.Sprintf(format, args...),
		input:   p.input,
		pos:     p.current.pos,
	}
}

func (p *parser) next() error {
	p.position++
	if p.position >= len(p.tokens) {
		return p.errorf("unexpected end of expression")
	}
	p.current = p.tokens[p.position]
	return nil
}

func (p *parser) expect(kind tokenKind, values ...string) error {
	if p.current.is(kind, values...) {
		return p.next()
	}
	return p.errorf("unexpected token %v", p.current)
}

func (p *parser) isEOF() bool {
	return p.current.is(eof)
}

// parse functions

func (p *parser) parseExpression(precedence int) (Node, error) {
	node, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	token := p.current
	for token.is(operator) {
		if op, ok := binaryOperators[token.value]; ok {
			if op.precedence >= precedence {
				if err = p.next(); err != nil {
					return nil, err
				}

				var expr Node
				if op.associativity == left {
					expr, err = p.parseExpression(op.precedence + 1)
					if err != nil {
						return nil, err
					}
				} else {
					expr, err = p.parseExpression(op.precedence)
					if err != nil {
						return nil, err
					}
				}

				node = binaryNode{operator: token.value, left: node, right: expr}
				token = p.current
				continue
			}
		}
		break
	}

	if precedence == 0 {
		node, err = p.parseConditionalExpression(node)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

func (p *parser) parsePrimary() (Node, error) {
	token := p.current

	if token.is(operator) {
		if op, ok := unaryOperators[token.value]; ok {
			if err := p.next(); err != nil {
				return nil, err
			}
			expr, err := p.parseExpression(op.precedence)
			if err != nil {
				return nil, err
			}

			return p.parsePostfixExpression(unaryNode{operator: token.value, node: expr})
		}
	}

	if token.is(punctuation, "(") {
		if err := p.next(); err != nil {
			return nil, err
		}

		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}

		err = p.expect(punctuation, ")")
		if err != nil {
			return nil, p.errorf("an opened parenthesis is not properly closed")
		}

		return p.parsePostfixExpression(expr)
	}

	return p.parsePrimaryExpression()
}

func (p *parser) parseConditionalExpression(node Node) (Node, error) {
	var err error
	var expr1, expr2 Node
	for p.current.is(punctuation, "?") {
		if err := p.next(); err != nil {
			return nil, err
		}
		if !p.current.is(punctuation, ":") {
			expr1, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}
			if err := p.expect(punctuation, ":"); err != nil {
				return nil, err
			}
			expr2, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}
		} else {
			if err := p.next(); err != nil {
				return nil, err
			}
			expr1 = node
			expr2, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}
		}

		node = conditionalNode{node, expr1, expr2}
	}
	return node, nil
}

func (p *parser) parsePrimaryExpression() (Node, error) {
	var err error
	var node Node
	token := p.current
	switch token.kind {
	case name:
		if err := p.next(); err != nil {
			return nil, err
		}
		switch token.value {
		case "true":
			return boolNode{value: true}, nil
		case "false":
			return boolNode{value: false}, nil
		case "nil":
			return nilNode{}, nil
		default:
			if p.current.is(punctuation, "(") {
				if _, ok := builtins[token.value]; ok {
					arguments, err := p.parseArguments()
					if err != nil {
						return nil, err
					}
					node = builtinNode{name: token.value, arguments: arguments}
				} else {
					if p.options.funcs != nil {
						if _, ok := p.options.funcs[token.value]; !ok {
							return nil, p.errorf("unknown func %v", token.value)
						}
					}
					arguments, err := p.parseArguments()
					if err != nil {
						return nil, err
					}
					node = functionNode{name: token.value, arguments: arguments}
				}
			} else {
				if p.options.names != nil {
					if _, ok := p.options.names[token.value]; !ok {
						return nil, p.errorf("unknown name %v", token.value)
					}
				}
				node = nameNode{name: token.value}
			}
		}

	case number:
		if err := p.next(); err != nil {
			return nil, err
		}
		number, err := strconv.ParseFloat(token.value, 64)
		if err != nil {
			return nil, p.errorf("%v", err)
		}
		return numberNode{value: number}, nil

	case text:
		if err := p.next(); err != nil {
			return nil, err
		}
		return textNode{value: token.value}, nil

	default:
		if token.is(punctuation, "[") {
			node, err = p.parseArrayExpression()
			if err != nil {
				return nil, err
			}
		} else if token.is(punctuation, "{") {
			node, err = p.parseMapExpression()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, p.errorf("unexpected token %v", token)
		}
	}

	return p.parsePostfixExpression(node)
}

func (p *parser) parseArrayExpression() (Node, error) {
	nodes, err := p.parseList("array items", "[", "]")
	if err != nil {
		return nil, err
	}
	return arrayNode{nodes}, nil
}

func (p *parser) parseMapExpression() (Node, error) {
	err := p.expect(punctuation, "{")
	if err != nil {
		return nil, err
	}

	nodes := make([]pairNode, 0)
	for !p.current.is(punctuation, "}") {
		if len(nodes) > 0 {
			err = p.expect(punctuation, ",")
			if err != nil {
				return nil, p.errorf("a map value must be followed by a comma")
			}
		}

		var key Node
		// a map key can be:
		//  * a number
		//  * a text
		//  * a name, which is equivalent to a string
		//  * an expression, which must be enclosed in parentheses -- (1 + 2)
		if p.current.is(number) || p.current.is(text) || p.current.is(name) {
			key = identifierNode{p.current.value}
			if err := p.next(); err != nil {
				return nil, err
			}
		} else if p.current.is(punctuation, "(") {
			key, err = p.parseExpression(0)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, p.errorf("a map key must be a quoted string, a number, a name, or an expression enclosed in parentheses (unexpected token %v)", p.current)
		}

		err = p.expect(punctuation, ":")
		if err != nil {
			return nil, p.errorf("a map key must be followed by a colon (:)")
		}

		node, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, pairNode{key, node})
	}

	err = p.expect(punctuation, "}")
	if err != nil {
		return nil, err
	}

	return mapNode{nodes}, nil
}

func (p *parser) parsePostfixExpression(node Node) (Node, error) {
	token := p.current
	for token.is(punctuation) {
		if token.value == "." {
			if err := p.next(); err != nil {
				return nil, err
			}
			token = p.current
			if err := p.next(); err != nil {
				return nil, err
			}

			if token.kind != name &&
				// Operators like "not" and "matches" are valid method or property names,
				//
				// In other words, besides name token kind, operator kind could also be parsed as a property or method.
				// This is because operators are processed by the lexer prior to names. So "not" in "foo.not()"
				// or "matches" in "foo.matches" will be recognized as an operator first. But in fact, "not"
				// and "matches" in such expressions shall be parsed as method or property names.
				//
				// And this ONLY works if the operator consists of valid characters for a property or method name.
				//
				// Other types, such as text kind and number kind, can't be parsed as property nor method names.
				//
				// As a result, if token is NOT an operator OR token.value is NOT a valid property or method name,
				// an error shall be returned.
				(token.kind != operator || !isValidIdentifier(token.value)) {
				return nil, p.errorf("expected name")
			}

			property := identifierNode{value: token.value}

			if p.current.is(punctuation, "(") {
				arguments, err := p.parseArguments()
				if err != nil {
					return nil, err
				}
				node = methodNode{node: node, property: property, arguments: arguments}
			} else {
				node = propertyNode{node: node, property: property}
			}

		} else if token.value == "[" {

			if err := p.next(); err != nil {
				return nil, err
			}

			arg, err := p.parseExpression(0)
			if err != nil {
				return err, nil
			}

			node = propertyNode{node: node, property: arg}

			err = p.expect(punctuation, "]")
			if err != nil {
				return nil, err
			}

		} else {
			break
		}

		token = p.current
	}
	return node, nil
}

func isValidIdentifier(str string) bool {
	if len(str) == 0 {
		return false
	}
	h, w := utf8.DecodeRuneInString(str)
	if !isAlphabetic(h) {
		return false
	}
	for _, r := range str[w:] {
		if !isAlphaNumeric(r) {
			return false
		}
	}
	return true
}

func (p *parser) parseArguments() ([]Node, error) {
	return p.parseList("arguments", "(", ")")
}

func (p *parser) parseList(what, start, end string) ([]Node, error) {
	err := p.expect(punctuation, start)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0)
	for !p.current.is(punctuation, end) {
		if len(nodes) > 0 {
			err = p.expect(punctuation, ",")
			if err != nil {
				return nil, p.errorf("%v must be separated by a comma", what)
			}
		}
		node, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	err = p.expect(punctuation, end)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
