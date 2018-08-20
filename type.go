package expr

import (
	"fmt"
	"reflect"
)

// Type is a reflect.Type alias.
type Type = reflect.Type

type typesTable map[string]Type

var (
	boolType      = reflect.TypeOf(true)
	numberType    = reflect.TypeOf(float64(0))
	textType      = reflect.TypeOf("")
	arrayType     = reflect.TypeOf([]interface{}{})
	mapType       = reflect.TypeOf(map[interface{}]interface{}{})
	interfaceType = reflect.TypeOf(new(interface{})).Elem()
)

func (n nilNode) Type(table typesTable) (Type, error) {
	return nil, nil
}

func (n identifierNode) Type(table typesTable) (Type, error) {
	return textType, nil
}

func (n numberNode) Type(table typesTable) (Type, error) {
	return numberType, nil
}

func (n boolNode) Type(table typesTable) (Type, error) {
	return boolType, nil
}

func (n textNode) Type(table typesTable) (Type, error) {
	return textType, nil
}

func (n nameNode) Type(table typesTable) (Type, error) {
	if t, ok := table[n.name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("unknown name %v", n)
}

func (n unaryNode) Type(table typesTable) (Type, error) {
	ntype, err := n.node.Type(table)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "!", "not":
		if isBoolType(ntype) || isInterfaceType(ntype) {
			return boolType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched type %v)`, n, ntype)
	}

	return interfaceType, nil
}

func (n binaryNode) Type(table typesTable) (Type, error) {
	var err error
	ltype, err := n.left.Type(table)
	if err != nil {
		return nil, err
	}
	rtype, err := n.right.Type(table)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case "==", "!=":
		if isComparable(ltype, rtype) {
			return boolType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)

	case "or", "||", "and", "&&":
		if (isBoolType(ltype) || isInterfaceType(ltype)) && (isBoolType(rtype) || isInterfaceType(rtype)) {
			return boolType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)

	case "<", ">", ">=", "<=":
		if (isNumberType(ltype) || isInterfaceType(ltype)) && (isNumberType(rtype) || isInterfaceType(rtype)) {
			return boolType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)

	case "/", "+", "-", "*", "**", "|", "^", "&", "%":
		if (isNumberType(ltype) || isInterfaceType(ltype)) && (isNumberType(rtype) || isInterfaceType(rtype)) {
			return numberType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)

	case "..":
		if (isNumberType(ltype) || isInterfaceType(ltype)) && (isNumberType(rtype) || isInterfaceType(rtype)) {
			return arrayType, nil
		}
		return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)

	}

	return interfaceType, nil
}

func (n matchesNode) Type(table typesTable) (Type, error) {
	var err error
	ltype, err := n.left.Type(table)
	if err != nil {
		return nil, err
	}
	rtype, err := n.right.Type(table)
	if err != nil {
		return nil, err
	}
	if (isStringType(ltype) || isInterfaceType(ltype)) && (isStringType(rtype) || isInterfaceType(rtype)) {
		return boolType, nil
	}
	return nil, fmt.Errorf(`invalid operation: %v (mismatched types %v and %v)`, n, ltype, rtype)
}

func (n propertyNode) Type(table typesTable) (Type, error) {
	ntype, err := n.node.Type(table)
	if err != nil {
		return nil, err
	}
	if t, ok := fieldType(ntype, n.property); ok {
		return t, nil
	}
	return nil, fmt.Errorf("%v undefined (type %v has no field %v)", n, ntype, n.property)
}

func (n indexNode) Type(table typesTable) (Type, error) {
	ntype, err := n.node.Type(table)
	if err != nil {
		return nil, err
	}
	_, err = n.index.Type(table)
	if err != nil {
		return nil, err
	}
	if t, ok := indexType(ntype); ok {
		return t, nil
	}
	return nil, fmt.Errorf("invalid operation: %v (type %v does not support indexing)", n, ntype)
}

func (n methodNode) Type(table typesTable) (Type, error) {
	ntype, err := n.node.Type(table)
	if err != nil {
		return nil, err
	}
	for _, node := range n.arguments {
		_, err := node.Type(table)
		if err != nil {
			return nil, err
		}
	}
	if t, ok := fieldType(ntype, n.method); ok {
		if f, ok := funcType(t); ok {
			return f, nil
		}
	}

	return nil, fmt.Errorf("%v undefined (type %v has no method %v)", n, ntype, n.method)
}

func (n builtinNode) Type(table typesTable) (Type, error) {
	for _, node := range n.arguments {
		_, err := node.Type(table)
		if err != nil {
			return nil, err
		}
	}
	switch n.name {
	case "len":
		// TODO: Add arguments type checks.
		return numberType, nil
	}
	return nil, fmt.Errorf("%v undefined", n)
}

func (n functionNode) Type(table typesTable) (Type, error) {
	for _, node := range n.arguments {
		_, err := node.Type(table)
		if err != nil {
			return nil, err
		}
	}
	if t, ok := table[n.name]; ok {
		if f, ok := funcType(t); ok {
			return f, nil
		}
	}
	return nil, fmt.Errorf("unknown func %v", n)
}

func (n conditionalNode) Type(table typesTable) (Type, error) {
	ctype, err := n.cond.Type(table)
	if err != nil {
		return nil, err
	}
	if !isBoolType(ctype) && !isInterfaceType(ctype) {
		return nil, fmt.Errorf("non-bool %v (type %v) used as condition", n.cond, ctype)
	}
	_, err = n.exp1.Type(table)
	if err != nil {
		return nil, err
	}
	_, err = n.exp2.Type(table)
	if err != nil {
		return nil, err
	}
	return boolType, nil
}

func (n arrayNode) Type(table typesTable) (Type, error) {
	for _, node := range n.nodes {
		_, err := node.Type(table)
		if err != nil {
			return nil, err
		}
	}
	return arrayType, nil
}

func (n mapNode) Type(table typesTable) (Type, error) {
	for _, node := range n.pairs {
		_, err := node.Type(table)
		if err != nil {
			return nil, err
		}
	}
	return mapType, nil
}

func (n pairNode) Type(table typesTable) (Type, error) {
	var err error
	_, err = n.key.Type(table)
	if err != nil {
		return nil, err
	}
	_, err = n.value.Type(table)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// helper funcs for reflect

func isComparable(l Type, r Type) bool {
	l = dereference(l)
	r = dereference(r)

	if l == nil || r == nil {
		return true // It is possible to compare with nil.
	}

	if isNumberType(l) && isNumberType(r) {
		return true
	} else if l.Kind() == reflect.Interface {
		return true
	} else if r.Kind() == reflect.Interface {
		return true
	} else if l == r {
		return true
	}
	return false
}

func isInterfaceType(t Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Interface:
			return true
		}
	}
	return false
}

func isNumberType(t Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		}
	}
	return false
}

func isBoolType(t Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Bool:
			return true
		}
	}
	return false
}

func isStringType(t Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.String:
			return true
		}
	}
	return false
}

func fieldType(ntype Type, name string) (Type, bool) {
	ntype = dereference(ntype)
	if ntype != nil {
		switch ntype.Kind() {
		case reflect.Interface:
			return interfaceType, true
		case reflect.Struct:
			if t, ok := ntype.FieldByName(name); ok {
				return t.Type, true
			}
		case reflect.Map:
			return ntype.Elem(), true
		}
	}
	return nil, false
}

func indexType(ntype Type) (Type, bool) {
	ntype = dereference(ntype)
	if ntype == nil {
		return nil, false
	}

	switch ntype.Kind() {
	case reflect.Interface:
		return interfaceType, true
	case reflect.Map, reflect.Array, reflect.Slice:
		return ntype.Elem(), true
	}

	return nil, false
}

func funcType(ntype Type) (Type, bool) {
	ntype = dereference(ntype)
	if ntype == nil {
		return nil, false
	}

	switch ntype.Kind() {
	case reflect.Interface:
		return interfaceType, true
	case reflect.Func:
		return ntype, true
	}

	return nil, false
}

func dereference(ntype Type) Type {
	if ntype == nil {
		return nil
	}
	if ntype.Kind() == reflect.Ptr {
		ntype = dereference(ntype.Elem())
	}
	return ntype
}
