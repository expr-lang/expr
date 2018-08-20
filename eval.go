package expr

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
)

// Eval parses and evaluates given input.
func Eval(input string, env interface{}) (interface{}, error) {
	node, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return Run(node, env)
}

// Run evaluates given ast.
func Run(node Node, env interface{}) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return node.Eval(env)
}

// eval functions

func (n nilNode) Eval(env interface{}) (interface{}, error) {
	return nil, nil
}

func (n identifierNode) Eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n numberNode) Eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n boolNode) Eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n textNode) Eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n nameNode) Eval(env interface{}) (interface{}, error) {
	v, ok := extract(env, n.name)
	if !ok {
		return nil, fmt.Errorf("undefined: %v", n)
	}
	return v, nil
}

func (n unaryNode) Eval(env interface{}) (interface{}, error) {
	val, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "not", "!":
		return !toBool(n, val), nil
	}

	v := toNumber(n, val)
	switch n.operator {
	case "-":
		return -v, nil
	case "+":
		return +v, nil
	}

	return nil, fmt.Errorf("implement unary %q operator", n.operator)
}

func (n binaryNode) Eval(env interface{}) (interface{}, error) {
	left, err := n.left.Eval(env)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "or", "||":
		if toBool(n.left, left) {
			return true, nil
		}
		right, err := n.right.Eval(env)
		if err != nil {
			return nil, err
		}
		return toBool(n.right, right), nil

	case "and", "&&":
		if toBool(n.left, left) {
			right, err := n.right.Eval(env)
			if err != nil {
				return nil, err
			}
			return toBool(n.right, right), nil
		}
		return false, nil
	}

	right, err := n.right.Eval(env)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "==":
		return equal(left, right), nil

	case "!=":
		return !equal(left, right), nil

	case "in":
		ok, err := contains(left, right)
		if err != nil {
			return nil, err
		}
		return ok, nil

	case "not in":
		ok, err := contains(left, right)
		if err != nil {
			return nil, err
		}
		return !ok, nil

	case "~":
		return toText(n.left, left) + toText(n.right, right), nil
	}

	// Next goes operators on numbers

	l, r := toNumber(n.left, left), toNumber(n.right, right)

	switch n.operator {
	case "|":
		return int(l) | int(r), nil

	case "^":
		return int(l) ^ int(r), nil

	case "&":
		return int(l) & int(r), nil

	case "<":
		return l < r, nil

	case ">":
		return l > r, nil

	case ">=":
		return l >= r, nil

	case "<=":
		return l <= r, nil

	case "+":
		return l + r, nil

	case "-":
		return l - r, nil

	case "*":
		return l * r, nil

	case "/":
		div := r
		if div == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return l / div, nil

	case "%":
		numerator := int64(l)
		denominator := int64(r)
		if denominator == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return float64(numerator % denominator), nil

	case "**":
		return math.Pow(l, r), nil

	case "..":
		return makeRange(int64(l), int64(r))
	}

	return nil, fmt.Errorf("implement %q operator", n.operator)
}

func makeRange(min, max int64) ([]float64, error) {
	size := max - min + 1
	if size > 1e6 {
		return nil, fmt.Errorf("range %v..%v exceeded max size of 1e6", min, max)
	}
	a := make([]float64, size)
	for i := range a {
		a[i] = float64(min + int64(i))
	}
	return a, nil
}

func (n matchesNode) Eval(env interface{}) (interface{}, error) {
	left, err := n.left.Eval(env)
	if err != nil {
		return nil, err
	}

	if n.r != nil {
		return n.r.MatchString(toText(n.left, left)), nil
	}

	right, err := n.right.Eval(env)
	if err != nil {
		return nil, err
	}

	matched, err := regexp.MatchString(toText(n.right, right), toText(n.left, left))
	if err != nil {
		return nil, err
	}
	return matched, nil
}

func (n propertyNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}
	p, ok := extract(v, n.property)
	if !ok {
		if isNil(v) {
			return nil, fmt.Errorf("%v is nil", n.node)
		}
		return nil, fmt.Errorf("%v undefined (type %T has no field %v)", n, v, n.property)
	}
	return p, nil
}

func (n indexNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}
	i, err := n.index.Eval(env)
	if err != nil {
		return nil, err
	}
	p, ok := extract(v, i)
	if !ok {
		return nil, fmt.Errorf("cannot get %q from %T: %v", i, v, n)
	}
	return p, nil
}

func (n methodNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}

	method, ok := extract(v, n.method)
	if !ok {
		return nil, fmt.Errorf("cannot get method %v from %T: %v", n.method, v, n)
	}

	in := make([]reflect.Value, 0)

	for _, a := range n.arguments {
		i, err := a.Eval(env)
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(i))
	}

	out := reflect.ValueOf(method).Call(in)

	if len(out) == 0 {
		return nil, nil
	} else if len(out) > 1 {
		return nil, fmt.Errorf("method %q must return only one value", n.method)
	}

	if out[0].IsValid() && out[0].CanInterface() {
		return out[0].Interface(), nil
	}

	return nil, nil
}

func (n builtinNode) Eval(env interface{}) (interface{}, error) {
	switch n.name {
	case "len":
		if len(n.arguments) == 0 {
			return nil, fmt.Errorf("missing argument: %v", n)
		}
		if len(n.arguments) > 1 {
			return nil, fmt.Errorf("too many arguments: %v", n)
		}

		i, err := n.arguments[0].Eval(env)
		if err != nil {
			return nil, err
		}

		switch reflect.TypeOf(i).Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			return float64(reflect.ValueOf(i).Len()), nil
		}
		return nil, fmt.Errorf("invalid argument %v (type %T)", n, i)
	}

	return nil, fmt.Errorf("unknown %q builtin", n.name)
}

func (n functionNode) Eval(env interface{}) (interface{}, error) {
	fn, ok := extract(env, n.name)
	if !ok {
		return nil, fmt.Errorf("undefined: %v", n.name)
	}

	in := make([]reflect.Value, 0)

	for _, a := range n.arguments {
		i, err := a.Eval(env)
		if err != nil {
			return nil, err
		}
		in = append(in, reflect.ValueOf(i))
	}

	out := reflect.ValueOf(fn).Call(in)

	if len(out) == 0 {
		return nil, nil
	} else if len(out) > 1 {
		return nil, fmt.Errorf("func %q must return only one value", n.name)
	}

	if out[0].IsValid() && out[0].CanInterface() {
		return out[0].Interface(), nil
	}

	return nil, nil
}

func (n conditionalNode) Eval(env interface{}) (interface{}, error) {
	cond, err := n.cond.Eval(env)
	if err != nil {
		return nil, err
	}

	// If
	if toBool(n.cond, cond) {
		// Then
		a, err := n.exp1.Eval(env)
		if err != nil {
			return nil, err
		}
		return a, nil
	}
	// Else
	b, err := n.exp2.Eval(env)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (n arrayNode) Eval(env interface{}) (interface{}, error) {
	array := make([]interface{}, 0)
	for _, node := range n.nodes {
		val, err := node.Eval(env)
		if err != nil {
			return nil, err
		}
		array = append(array, val)
	}
	return array, nil
}

func (n mapNode) Eval(env interface{}) (interface{}, error) {
	m := make(map[interface{}]interface{})
	for _, pair := range n.pairs {
		key, err := pair.key.Eval(env)
		if err != nil {
			return nil, err
		}
		value, err := pair.value.Eval(env)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}
