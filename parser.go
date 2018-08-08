package expr

import (
	"fmt"
	"reflect"
	"regexp"
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
	strict   bool
	types    map[string]Type
}

// OptionFn for configuring parser.
type OptionFn func(p *parser)

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
		types:   make(map[string]Type),
	}

	for _, op := range ops {
		op(p)
	}

	node, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	if !p.isEOF() {
		return nil, p.errorf("unexpected token %v", p.current)
	}

	if p.strict {
		_, err = node.Type(p.types)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

// Define sets variable for type checks during parsing.
func Define(name string, t interface{}) OptionFn {
	return func(p *parser) {
		p.strict = true
		p.types[name] = reflect.TypeOf(t)
	}
}

// With sets variables for type checks during parsing.
// If struct is passed, all fields will be treated as variables.
// If map is passed, all items will be treated as variables
// (key as name, value as type).
func With(i interface{}) OptionFn {
	return func(p *parser) {
		p.strict = true
		v := reflect.ValueOf(i)
		t := reflect.TypeOf(i)
		t = dereference(t)
		if t == nil {
			return
		}

		switch t.Kind() {
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				p.types[f.Name] = f.Type
			}
		case reflect.Map:
			for _, key := range v.MapKeys() {
				value := v.MapIndex(key)
				if key.Kind() == reflect.String && value.IsValid() && value.CanInterface() {
					p.types[key.String()] = reflect.TypeOf(value.Interface())
				}
			}
		}
	}
}

func (p *parser) errorf(format string, args ...interface{}) *syntaxError {
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

				if token.is(operator, "matches") {
					var r *regexp.Regexp
					if s, ok := expr.(textNode); ok {
						r, err = regexp.Compile(s.value)
						if err != nil {
							return nil, p.errorf("%v", err)
						}
					}
					node = matchesNode{r: r, left: node, right: expr}
				} else {
					node = binaryNode{operator: token.value, left: node, right: expr}
				}
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
			node, err = p.parseNameExpression(token)
			if err != nil {
				return nil, err
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
			return nil, p.errorf("unexpected token %v", token).at(token)
		}
	}

	return p.parsePostfixExpression(node)
}

func (p *parser) parseNameExpression(token token) (Node, error) {
	var node Node
	if p.current.is(punctuation, "(") {
		arguments, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		if _, ok := builtins[token.value]; ok {
			node = builtinNode{name: token.value, arguments: arguments}
		} else {
			node = functionNode{name: token.value, arguments: arguments}
		}
	} else {
		node = nameNode{name: token.value}
	}
	return node, nil
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
				return nil, p.errorf("expected name").at(token)
			}

			if p.current.is(punctuation, "(") {
				arguments, err := p.parseArguments()
				if err != nil {
					return nil, err
				}
				node = methodNode{node: node, method: token.value, arguments: arguments}
			} else {
				node = propertyNode{node: node, property: token.value}
			}

		} else if token.value == "[" {

			if err := p.next(); err != nil {
				return nil, err
			}

			arg, err := p.parseExpression(0)
			if err != nil {
				return nil, err
			}

			node = indexNode{node: node, index: arg}

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
