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
	Type    reflect.Type
	SubType SubType
	Func    *builtin.Function
	Method  bool
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
			return array.Of
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
		return Nature{Type: method.Type, Method: true}, true
	}
}

func (n Nature) NumField() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumField()
}

func (n Nature) Field(i int) reflect.StructField {
	if n.Type == nil {
		return reflect.StructField{}
	}
	return n.Type.Field(i)
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
