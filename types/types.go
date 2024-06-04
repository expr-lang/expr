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
