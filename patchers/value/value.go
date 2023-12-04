// Package value provides a Patcher that uses interfaces to allow custom types that can be represented as standard go values to be used more easily in expressions.
//
// # Example Usage
//
//	import (
//		"fmt"
//		"github.com/antonmedv/expr/patchers/value"
//		"github.com/antonmedv/expr"
//	)
//
//	type customInt struct {
//	       	Int int
//	}
//
//	// Provides type checking at compile time
//	func (v *customInt) IntValue() int {
//	       	return v.Int
//	}
//
//	// Lets us return nil if we need to
//	func (v *customInt) ExprValue() any {
//	       	return v.Int
//	}
//
//	func main() {
//		env := make(map[string]any)
//		env["ValueOne"] = &customInt{1}
//		env["ValueTwo"] = &customInt{2}
//
//		program, err := expr.Compile("ValueOne + ValueTwo", expr.Env(env), value.Patcher)
//
//		if err != nil {
//			panic(err)
//		}
//
//		out, err := vm.Run(program, env)
//
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("Got %v", out)
//	}
package value

import (
	"reflect"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/conf"
)

// Patcher is an expr.Option that both patches the program and adds the `$patcher_value_getter` function.
// Use it directly as an Option to expr.Compile()
var Patcher = func() expr.Option {
	vPatcher := patcher{}
	return func(c *conf.Config) {
		c.Visitors = append(c.Visitors, vPatcher)
		vPatcher.ApplyOptions(c)
	}
}()

// A ExprValuer provides a generic function for a custom type to return standard go values.
// It allows for returning a `nil` value but does not provide any type checking at expression compile.
//
// A custom type may implement both ExprValuer and a type specific interface to enable both
// compile time checking and the ability to return a `nil` value.
type ExprValuer interface {
	ExprValue() any
}

type IntValuer interface {
	IntValue() int
}

type BoolValuer interface {
	BoolValue() bool
}

type Int8Valuer interface {
	Int8Value() int8
}

type Int16Valuer interface {
	Int16Value() int16
}

type Int32Valuer interface {
	Int32Value() int32
}

type Int64Valuer interface {
	Int64Value() int64
}

type UintValuer interface {
	UintValue() uint
}

type Uint8Valuer interface {
	Uint8Value() uint8
}

type Uint16Valuer interface {
	Uint16Value() uint16
}

type Uint32Valuer interface {
	Uint32Value() uint32
}

type Uint64Valuer interface {
	Uint64Value() uint64
}

type Float32Valuer interface {
	Float32Value() float32
}

type Float64Valuer interface {
	Float64Value() float64
}

type StringValuer interface {
	StringValue() string
}

type TimeValuer interface {
	TimeValue() time.Time
}

type DurationValuer interface {
	DurationValue() time.Duration
}

type ArrayValuer interface {
	ArrayValue() []any
}

type MapValuer interface {
	MapValue() map[string]any
}

var supportedInterfaces = []reflect.Type{
	reflect.TypeOf((*ExprValuer)(nil)).Elem(),
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
	id, ok := (*node).(*ast.IdentifierNode)
	if !ok {
		return
	}

	nodeType := id.Type()

	for _, t := range supportedInterfaces {
		if nodeType.Implements(t) {
			callnode := &ast.CallNode{
				Callee:    &ast.IdentifierNode{Value: "$patcher_value_getter"},
				Arguments: []ast.Node{id},
			}

			ast.Patch(node, callnode)
		}
	}
}

func (patcher) ApplyOptions(c *conf.Config) {
	getExprValueFunc(c)
}

func getExprValue(params ...any) (any, error) {
	switch v := params[0].(type) {
	case ExprValuer:
		return v.ExprValue(), nil
	case BoolValuer:
		return v.BoolValue(), nil
	case IntValuer:
		return v.IntValue(), nil
	case Int8Valuer:
		return v.Int8Value(), nil
	case Int16Valuer:
		return v.Int16Value(), nil
	case Int32Valuer:
		return v.Int32Value(), nil
	case Int64Valuer:
		return v.Int64Value(), nil
	case UintValuer:
		return v.UintValue(), nil
	case Uint8Valuer:
		return v.Uint8Value(), nil
	case Uint16Valuer:
		return v.Uint16Value(), nil
	case Uint32Valuer:
		return v.Uint32Value(), nil
	case Uint64Valuer:
		return v.Uint64Value(), nil
	case Float32Valuer:
		return v.Float32Value(), nil
	case Float64Valuer:
		return v.Float64Value(), nil
	case StringValuer:
		return v.StringValue(), nil
	case TimeValuer:
		return v.TimeValue(), nil
	case DurationValuer:
		return v.DurationValue(), nil
	case ArrayValuer:
		return v.ArrayValue(), nil
	case MapValuer:
		return v.MapValue(), nil
	}

	return params[0], nil
}

var getExprValueFunc = expr.Function("$patcher_value_getter", getExprValue,
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
