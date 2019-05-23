package checker

import "reflect"

var (
	nilType       = reflect.TypeOf(nil)
	boolType      = reflect.TypeOf(true)
	integerType   = reflect.TypeOf(int64(0))
	floatType     = reflect.TypeOf(float64(0))
	stringType    = reflect.TypeOf("")
	arrayType     = reflect.TypeOf([]interface{}{})
	mapType       = reflect.TypeOf(map[interface{}]interface{}{})
	interfaceType = reflect.TypeOf(new(interface{})).Elem()
)

func dereference(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Ptr {
		t = dereference(t.Elem())
	}
	return t
}

func isComparable(l, r reflect.Type) bool {
	l = dereference(l)
	r = dereference(r)

	if l == nil || r == nil {
		return true // It is possible to compare with nil.
	}

	if isInteger(l) && isInteger(r) {
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

func isInterface(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Interface:
			return true
		}
	}
	return false
}

func isInteger(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		}
	}
	return false
}

func isFloat(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Float32, reflect.Float64:
			return true
		}
	}
	return false
}

func isBool(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Bool:
			return true
		}
	}
	return false
}

func isString(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.String:
			return true
		}
	}
	return false
}

func isArray(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			return true
		}
	}
	return false
}

func isMap(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Map:
			return true
		}
	}
	return false
}

func isStruct(t reflect.Type) bool {
	t = dereference(t)
	if t != nil {
		switch t.Kind() {
		case reflect.Struct:
			return true
		}
	}
	return false
}

func fieldType(ntype reflect.Type, name string) (reflect.Type, bool) {
	ntype = dereference(ntype)
	if ntype != nil {
		switch ntype.Kind() {
		case reflect.Interface:
			return interfaceType, true
		case reflect.Struct:
			// First check all struct's fields.
			for i := 0; i < ntype.NumField(); i++ {
				f := ntype.Field(i)
				if !f.Anonymous && f.Name == name {
					return f.Type, true
				}
			}

			// Second check fields of embedded structs.
			for i := 0; i < ntype.NumField(); i++ {
				f := ntype.Field(i)
				if f.Anonymous {
					if t, ok := fieldType(f.Type, name); ok {
						return t, true
					}
				}
			}
		case reflect.Map:
			return ntype.Elem(), true
		}
	}

	return nil, false
}

func methodType(ntype reflect.Type, name string) (reflect.Type, bool) {
	ntype = dereference(ntype)
	if ntype != nil {
		switch ntype.Kind() {
		case reflect.Interface:
			return interfaceType, true
		case reflect.Struct:
			// First check all struct's methods.
			for i := 0; i < ntype.NumMethod(); i++ {
				m := ntype.Method(i)
				if m.Name == name {
					return m.Type, true
				}
			}

			// Second check all struct's fields.
			for i := 0; i < ntype.NumField(); i++ {
				f := ntype.Field(i)
				if !f.Anonymous && f.Name == name {
					return f.Type, true
				}
			}

			// Third check fields of embedded structs.
			for i := 0; i < ntype.NumField(); i++ {
				f := ntype.Field(i)
				if f.Anonymous {
					if t, ok := methodType(f.Type, name); ok {
						return t, true
					}
				}
			}
		case reflect.Map:
			return ntype.Elem(), true
		}
	}

	return nil, false
}

func indexType(ntype reflect.Type) (reflect.Type, bool) {
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

func funcType(ntype reflect.Type) (reflect.Type, bool) {
	ntype = dereference(ntype)
	if ntype == nil {
		return nil, false
	}

	switch ntype.Kind() {
	case reflect.Interface:
		return interfaceType, true
	case reflect.Func:
		if ntype.NumOut() > 0 {
			return ntype.Out(0), true
		}
		return nilType, true
	}

	return nil, false
}
