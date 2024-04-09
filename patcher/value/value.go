// Package value provides a Patcher that uses interfaces to allow custom types that can be represented as standard go values to be used more easily in expressions.
package value

import (
	"reflect"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/conf"
)

// ValueGetter is a Patcher that allows custom types to be represented as standard go values for use with expr.
// It also adds the `$patcher_value_getter` function to the program for efficiently calling matching interfaces.
//
// The purpose of this Patcher is to make it seamless to use custom types in expressions without the need to
// first convert them to standard go values. It may also facilitate using already existing structs or maps as
// environments when they contain compatible types.
//
// An example usage may be modeling a database record with columns that have varying data types and constraints.
// In such an example you may have custom types that, beyond storing a simple value, such as an integer, may
// contain metadata such as column type and if a value is specifically a NULL value.
//
// Use it directly as an Option to expr.Compile()
var ValueGetter = expr.Option(func(c *conf.Config) {
	c.Visitors = append(c.Visitors, patcher{})
	getValueFunc(c)
})

// A AnyValuer provides a generic function for a custom type to return standard go values.
// It allows for returning a `nil` value but does not provide any type checking at expression compile.
//
// A custom type may implement both AnyValuer and a type specific interface to enable both
// compile time checking and the ability to return a `nil` value.
type AnyValuer interface {
	AsAny() any
}

type IntValuer interface {
	AsInt() int
}

type BoolValuer interface {
	AsBool() bool
}

type Int8Valuer interface {
	AsInt8() int8
}

type Int16Valuer interface {
	AsInt16() int16
}

type Int32Valuer interface {
	AsInt32() int32
}

type Int64Valuer interface {
	AsInt64() int64
}

type UintValuer interface {
	AsUint() uint
}

type Uint8Valuer interface {
	AsUint8() uint8
}

type Uint16Valuer interface {
	AsUint16() uint16
}

type Uint32Valuer interface {
	AsUint32() uint32
}

type Uint64Valuer interface {
	AsUint64() uint64
}

type Float32Valuer interface {
	AsFloat32() float32
}

type Float64Valuer interface {
	AsFloat64() float64
}

type StringValuer interface {
	AsString() string
}

type TimeValuer interface {
	AsTime() time.Time
}

type DurationValuer interface {
	AsDuration() time.Duration
}

type ArrayValuer interface {
	AsArray() []any
}

type MapValuer interface {
	AsMap() map[string]any
}

var supportedInterfaces = []reflect.Type{
	reflect.TypeOf((*AnyValuer)(nil)).Elem(),
	reflect.TypeOf((*BoolValuer)(nil)).Elem(),
	reflect.TypeOf((*IntValuer)(nil)).Elem(),
	reflect.TypeOf((*Int8Valuer)(nil)).Elem(),
	reflect.TypeOf((*Int16Valuer)(nil)).Elem(),
	reflect.TypeOf((*Int32Valuer)(nil)).Elem(),
	reflect.TypeOf((*Int64Valuer)(nil)).Elem(),
	reflect.TypeOf((*UintValuer)(nil)).Elem(),
	reflect.TypeOf((*Uint8Valuer)(nil)).Elem(),
	reflect.TypeOf((*Uint16Valuer)(nil)).Elem(),
	reflect.TypeOf((*Uint32Valuer)(nil)).Elem(),
	reflect.TypeOf((*Uint64Valuer)(nil)).Elem(),
	reflect.TypeOf((*Float32Valuer)(nil)).Elem(),
	reflect.TypeOf((*Float64Valuer)(nil)).Elem(),
	reflect.TypeOf((*StringValuer)(nil)).Elem(),
	reflect.TypeOf((*TimeValuer)(nil)).Elem(),
	reflect.TypeOf((*DurationValuer)(nil)).Elem(),
	reflect.TypeOf((*ArrayValuer)(nil)).Elem(),
	reflect.TypeOf((*MapValuer)(nil)).Elem(),
}

type patcher struct{}

func (patcher) Visit(node *ast.Node) {
	switch id := (*node).(type) {
	case *ast.IdentifierNode, *ast.MemberNode:
		nodeType := id.Type()
		for _, t := range supportedInterfaces {
			if nodeType.Implements(t) {
				ast.Patch(node, &ast.CallNode{
					Callee:    &ast.IdentifierNode{Value: "$patcher_value_getter"},
					Arguments: []ast.Node{id},
				})
				return
			}
		}
	}
}

func getValue(params ...any) (any, error) {
	switch v := params[0].(type) {
	case AnyValuer:
		return v.AsAny(), nil
	case BoolValuer:
		return v.AsBool(), nil
	case IntValuer:
		return v.AsInt(), nil
	case Int8Valuer:
		return v.AsInt8(), nil
	case Int16Valuer:
		return v.AsInt16(), nil
	case Int32Valuer:
		return v.AsInt32(), nil
	case Int64Valuer:
		return v.AsInt64(), nil
	case UintValuer:
		return v.AsUint(), nil
	case Uint8Valuer:
		return v.AsUint8(), nil
	case Uint16Valuer:
		return v.AsUint16(), nil
	case Uint32Valuer:
		return v.AsUint32(), nil
	case Uint64Valuer:
		return v.AsUint64(), nil
	case Float32Valuer:
		return v.AsFloat32(), nil
	case Float64Valuer:
		return v.AsFloat64(), nil
	case StringValuer:
		return v.AsString(), nil
	case TimeValuer:
		return v.AsTime(), nil
	case DurationValuer:
		return v.AsDuration(), nil
	case ArrayValuer:
		return v.AsArray(), nil
	case MapValuer:
		return v.AsMap(), nil
	}

	return params[0], nil
}

var getValueFunc = expr.Function("$patcher_value_getter", getValue,
	new(func(BoolValuer) bool),
	new(func(IntValuer) int),
	new(func(Int8Valuer) int8),
	new(func(Int16Valuer) int16),
	new(func(Int32Valuer) int32),
	new(func(Int64Valuer) int64),
	new(func(UintValuer) uint),
	new(func(Uint8Valuer) uint8),
	new(func(Uint16Valuer) uint16),
	new(func(Uint32Valuer) uint32),
	new(func(Uint64Valuer) uint64),
	new(func(Float32Valuer) float32),
	new(func(Float64Valuer) float64),
	new(func(StringValuer) string),
	new(func(TimeValuer) time.Time),
	new(func(DurationValuer) time.Duration),
	new(func(ArrayValuer) []any),
	new(func(MapValuer) map[string]any),
	new(func(any) any),
)
