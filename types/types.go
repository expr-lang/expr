package types

import (
	"reflect"

	. "github.com/expr-lang/expr/checker/nature"
)

func TypeOf(v any) Type {
	return rtype{t: reflect.TypeOf(v)}
}

var (
	Int     = TypeOf(0)
	Int8    = TypeOf(int8(0))
	Int16   = TypeOf(int16(0))
	Int32   = TypeOf(int32(0))
	Int64   = TypeOf(int64(0))
	Uint    = TypeOf(uint(0))
	Uint8   = TypeOf(uint8(0))
	Uint16  = TypeOf(uint16(0))
	Uint32  = TypeOf(uint32(0))
	Uint64  = TypeOf(uint64(0))
	Float   = TypeOf(float32(0))
	Float64 = TypeOf(float64(0))
	String  = TypeOf("")
	Bool    = TypeOf(true)
	Nil     = nilType{}
)

// Type is a type that can be used to represent a value.
type Type interface {
	Nature() Nature
}

type nilType struct{}

func (nilType) Nature() Nature {
	return Nature{Nil: true}
}

type rtype struct {
	t reflect.Type
}

func (r rtype) Nature() Nature {
	return Nature{Type: r.t}
}

// Map returns a type that represents a map of the given type.
// The map is not strict, meaning that it can contain keys not defined in the map.
type Map map[string]Type

func (m Map) Nature() Nature {
	nt := Nature{
		Type:   reflect.TypeOf(map[string]any{}),
		Fields: make(map[string]Nature, len(m)),
	}
	for k, v := range m {
		nt.Fields[k] = v.Nature()
	}
	return nt
}

// StrictMap returns a type that represents a map of the given type.
// The map is strict, meaning that it can only contain keys defined in the map.
type StrictMap map[string]Type

func (m StrictMap) Nature() Nature {
	nt := Nature{
		Type:   reflect.TypeOf(map[string]any{}),
		Fields: make(map[string]Nature, len(m)),
		Strict: true,
	}
	for k, v := range m {
		nt.Fields[k] = v.Nature()
	}
	return nt
}

// Array returns a type that represents an array of the given type.
func Array(of Type) Type {
	return array{of}
}

type array struct {
	of Type
}

func (a array) Nature() Nature {
	of := a.of.Nature()
	return Nature{
		Type:    reflect.TypeOf([]any{}),
		Fields:  make(map[string]Nature, 1),
		ArrayOf: &of,
	}
}
