package expr

import (
	"fmt"
	"math"
	"regexp"
)

type evaluator interface {
	eval(env interface{}) (interface{}, error)
}

// Eval parses and evaluates given input.
func Eval(input string, env interface{}) (node interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	node, err = Parse(input)

	if err != nil {
		return nil, err
	}
	return Run(node, env)
}

// Run evaluates given ast.
func Run(node Node, env interface{}) (interface{}, error) {
	if e, ok := node.(evaluator); ok {
		return e.eval(env)
	}
	return nil, fmt.Errorf("implement evaluator for %T", node)
}

// eval functions

func (n nilNode) eval(env interface{}) (interface{}, error) {
	return nil, nil
}

func (n identifierNode) eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n numberNode) eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n boolNode) eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n textNode) eval(env interface{}) (interface{}, error) {
	return n.value, nil
}

func (n nameNode) eval(env interface{}) (interface{}, error) {
	return extract(env, n.name)
}

func (n unaryNode) eval(env interface{}) (interface{}, error) {
	val, err := Run(n.node, env)
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

func (n binaryNode) eval(env interface{}) (interface{}, error) {
	left, err := Run(n.left, env)
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

		right, err := Run(n.right, env)
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
			right, err := Run(n.right, env)
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

	right, err := Run(n.right, env)
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

	case "matches":
		if isText(left) && isText(right) {
			matched, err := regexp.MatchString(toText(right), toText(left))
			if err != nil {
				return nil, err
			}
			return matched, nil
		}
		return nil, fmt.Errorf("operator matches not defined on (%T, %T)", left, right)

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

func (n propertyNode) eval(env interface{}) (interface{}, error) {
	v, err := Run(n.node, env)
	if err != nil {
		return nil, err
	}
	p, err := Run(n.property, env)
	if err != nil {
		return nil, err
	}

	return extract(v, p)
}

func (n methodNode) eval(env interface{}) (interface{}, error) {
	v, err := Run(n.node, env)
	if err != nil {
		return nil, err
	}

	method, err := extract(v, n.property.value)
	if err != nil {
		return nil, err
	}

	return call(n.property.value, method, n.arguments, env)
}

func (n functionNode) eval(env interface{}) (interface{}, error) {
	fn, err := extract(env, n.name)
	if err != nil {
		return nil, err
	}

	return call(n.name, fn, n.arguments, env)
}

func (n conditionalNode) eval(env interface{}) (interface{}, error) {
	cond, err := Run(n.cond, env)
	if err != nil {
		return nil, err
	}

	// If
	if !isBool(cond) {
		return nil, fmt.Errorf("non-bool value used in cond (%T)", cond)
	}
	// Then
	if toBool(cond) {
		a, err := Run(n.exp1, env)
		if err != nil {
			return nil, err
		}
		return a, nil
	}
	// Else
	b, err := Run(n.exp2, env)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (n arrayNode) eval(env interface{}) (interface{}, error) {
	array := make([]interface{}, 0)
	for _, node := range n.nodes {
		val, err := Run(node, env)
		if err != nil {
			return nil, err
		}
		array = append(array, val)
	}
	return array, nil
}

func (n mapNode) eval(env interface{}) (interface{}, error) {
	m := make(map[interface{}]interface{})
	for _, pair := range n.pairs {
		key, err := Run(pair.key, env)
		if err != nil {
			return nil, err
		}
		value, err := Run(pair.value, env)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}
