package nature

import (
	"reflect"

	"github.com/expr-lang/expr/builtin"
	"github.com/expr-lang/expr/internal/deref"
)

var (
	unknown = Nature{}
)

type Nature struct {
	Type        reflect.Type
	SubType     SubType
	Func        *builtin.Function
	Method      bool
	MethodIndex int
	FieldIndex  []int
}

func (n Nature) String() string {
	if n.SubType != nil {
		return n.SubType.String()
	}
	if n.Type != nil {
		return n.Type.String()
	}
	return "unknown"
}

func (n Nature) Deref() Nature {
	if n.Type != nil {
		n.Type = deref.Type(n.Type)
	}
	return n
}

func (n Nature) Kind() reflect.Kind {
	if n.Type != nil {
		return n.Type.Kind()
	}
	return reflect.Invalid
}

func (n Nature) Key() Nature {
	if n.Kind() == reflect.Map {
		return Nature{Type: n.Type.Key()}
	}
	return unknown
}

func (n Nature) Elem() Nature {
	switch n.Kind() {
	case reflect.Map, reflect.Ptr:
		return Nature{Type: n.Type.Elem()}
	case reflect.Array, reflect.Slice:
		if array, ok := n.SubType.(Array); ok {
			return array.Elem
		}
		return Nature{Type: n.Type.Elem()}
	}
	return unknown
}

func (n Nature) AssignableTo(nt Nature) bool {
	if n.Type == nil || nt.Type == nil {
		return false
	}
	return n.Type.AssignableTo(nt.Type)
}

func (n Nature) MethodByName(name string) (Nature, bool) {
	if n.Type == nil {
		return unknown, false
	}
	method, ok := n.Type.MethodByName(name)
	if !ok {
		return unknown, false
	}

	if n.Type.Kind() == reflect.Interface {
		// In case of interface type method will not have a receiver,
		// and to prevent checker decreasing numbers of in arguments
		// return method type as not method (second argument is false).

		// Also, we can not use m.Index here, because it will be
		// different indexes for different types which implement
		// the same interface.
		return Nature{Type: method.Type}, true
	} else {
		return Nature{
			Type:        method.Type,
			Method:      true,
			MethodIndex: method.Index,
		}, true
	}
}

func (n Nature) NumIn() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumIn()
}

func (n Nature) In(i int) Nature {
	if n.Type == nil {
		return unknown
	}
	return Nature{Type: n.Type.In(i)}
}

func (n Nature) NumOut() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumOut()
}

func (n Nature) Out(i int) Nature {
	if n.Type == nil {
		return unknown
	}
	return Nature{Type: n.Type.Out(i)}
}

func (n Nature) IsVariadic() bool {
	if n.Type == nil {
		return false
	}
	return n.Type.IsVariadic()
}

func (n Nature) FieldByName(name string) (Nature, bool) {
	if n.Type == nil {
		return unknown, false
	}
	field, ok := fetchField(n.Type, name)
	return Nature{Type: field.Type, FieldIndex: field.Index}, ok
}

func (n Nature) IsFastMap() bool {
	if n.Type == nil {
		return false
	}
	if n.Type.Kind() == reflect.Map &&
		n.Type.Key().Kind() == reflect.String &&
		n.Type.Elem().Kind() == reflect.Interface {
		return true
	}
	return false
}

func (n Nature) Get(name string) (Nature, bool) {
	if n.Type == nil {
		return unknown, false
	}

	if m, ok := n.MethodByName(name); ok {
		return m, true
	}

	t := deref.Type(n.Type)

	switch t.Kind() {
	case reflect.Struct:
		if f, ok := fetchField(t, name); ok {
			return Nature{
				Type:       f.Type,
				FieldIndex: f.Index,
			}, true
		}
	case reflect.Map:
		if f, ok := n.SubType.Get(name); ok {
			return f, true
		}
	}
	return unknown, false
}

func (n Nature) List() map[string]Nature {
	table := make(map[string]Nature)

	if n.Type == nil {
		return table
	}

	for i := 0; i < n.Type.NumMethod(); i++ {
		method := n.Type.Method(i)
		table[method.Name] = Nature{
			Type:        method.Type,
			Method:      true,
			MethodIndex: method.Index,
		}
	}

	switch n.Type.Kind() {
	case reflect.Struct:
		for name, nt := range fields(n.Type) {
			if _, ok := table[name]; ok {
				continue
			}
			table[name] = nt
		}

	case reflect.Map:
		v := reflect.ValueOf(n.SubType)
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if key.Kind() == reflect.String && value.IsValid() && value.CanInterface() {
				table[key.String()] = Nature{Type: reflect.TypeOf(value.Interface())}
			}
		}
	}

	return table
}
