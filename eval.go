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
	return extract(env, n.name)
}

func (n unaryNode) Eval(env interface{}) (interface{}, error) {
	val, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "not", "!":
		if !isBool(val) {
			return nil, fmt.Errorf("negation of non-bool type %T", val)
		}
		return !toBool(val), nil
	}

	v, err := cast(val)
	if err != nil {
		return nil, err
	}
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
		if !isBool(left) {
			return nil, fmt.Errorf("non-bool value in cond (%T)", left)
		}

		if toBool(left) {
			return true, nil
		}

		right, err := n.right.Eval(env)
		if err != nil {
			return nil, err
		}
		if !isBool(right) {
			return nil, fmt.Errorf("non-bool value in cond (%T)", right)
		}

		return toBool(right), nil

	case "and", "&&":
		if !isBool(left) {
			return nil, fmt.Errorf("non-bool value in cond (%T)", left)
		}

		if toBool(left) {
			right, err := n.right.Eval(env)
			if err != nil {
				return nil, err
			}
			if !isBool(right) {
				return nil, fmt.Errorf("non-bool value in cond (%T)", right)
			}
			return toBool(right), nil
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
		if isText(left) && isText(right) {
			return toText(left) + toText(right), nil
		}
		return nil, fmt.Errorf("operator ~ not defined on (%T, %T)", left, right)
	}

	// Next goes operators on numbers

	l, err := cast(left)
	if err != nil {
		return nil, err
	}

	r, err := cast(right)
	if err != nil {
		return nil, err
	}

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
		if isText(left) {
			return n.r.MatchString(toText(left)), nil
		}
	}

	right, err := n.right.Eval(env)
	if err != nil {
		return nil, err
	}

	if isText(left) && isText(right) {
		matched, err := regexp.MatchString(toText(right), toText(left))
		if err != nil {
			return nil, err
		}
		return matched, nil
	}

	return nil, fmt.Errorf("operator matches doesn't defined on (%T, %T): %v", left, right, n)
}

func (n propertyNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}
	return extract(v, n.property)
}

func (n indexNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}
	p, err := n.index.Eval(env)
	if err != nil {
		return nil, err
	}
	return extract(v, p)
}

func (n methodNode) Eval(env interface{}) (interface{}, error) {
	v, err := n.node.Eval(env)
	if err != nil {
		return nil, err
	}

	method, err := extract(v, n.method)
	if err != nil {
		return nil, err
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
	fn, err := extract(env, n.name)
	if err != nil {
		return nil, err
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
	if !isBool(cond) {
		return nil, fmt.Errorf("non-bool value used in cond (%T)", cond)
	}
	// Then
	if toBool(cond) {
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
