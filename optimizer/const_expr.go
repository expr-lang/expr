package optimizer

import (
	"fmt"
	"reflect"
	"strings"

	. "github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/file"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type constExpr struct {
	applied bool
	err     error
	fns     map[string]reflect.Value
}

func (c *constExpr) Visit(node *Node) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			// Make message more actual, it's a runtime error, but at compile step.
			msg = strings.Replace(msg, "runtime error:", "compile error:", 1)
			c.err = &file.Error{
				Location: (*node).Location(),
				Message:  msg,
			}
		}
	}()

	if call, ok := (*node).(*CallNode); ok {
		if name, ok := call.Callee.(*IdentifierNode); ok {
			fn, ok := c.fns[name.Value]
			if ok {
				in := make([]reflect.Value, len(call.Arguments))
				for i := 0; i < len(call.Arguments); i++ {
					arg := call.Arguments[i]
					var param any

					switch a := arg.(type) {
					case *NilNode:
						param = nil
					case *IntegerNode:
						if param, ok = literal(a.Value, a.Type()); !ok {
							return
						}
					case *FloatNode:
						if param, ok = literal(a.Value, a.Type()); !ok {
							return
						}
					case *BoolNode:
						param = a.Value
					case *StringNode:
						param = a.Value
					case *ConstantNode:
						param = a.Value

					default:
						return // Const expr optimization not applicable.
					}

					if param == nil && reflect.TypeOf(param) == nil {
						// In case of nil value and nil type use this hack,
						// otherwise reflect.Call will panic on zero value.
						in[i] = reflect.ValueOf(&param).Elem()
					} else {
						in[i] = reflect.ValueOf(param)
					}
				}

				out := fn.Call(in)
				value := out[0].Interface()
				if len(out) == 2 && out[1].Type() == errorType && !out[1].IsNil() {
					c.err = out[1].Interface().(error)
					return
				}
				constNode := &ConstantNode{Value: value}
				patchWithType(node, constNode)
				c.applied = true
			}
		}
	}
}

var literalTypes = map[reflect.Kind]reflect.Type{
	reflect.Int:     reflect.TypeOf(int(0)),
	reflect.Int8:    reflect.TypeOf(int8(0)),
	reflect.Int16:   reflect.TypeOf(int16(0)),
	reflect.Int32:   reflect.TypeOf(int32(0)),
	reflect.Int64:   reflect.TypeOf(int64(0)),
	reflect.Uint:    reflect.TypeOf(uint(0)),
	reflect.Uint8:   reflect.TypeOf(uint8(0)),
	reflect.Uint16:  reflect.TypeOf(uint16(0)),
	reflect.Uint32:  reflect.TypeOf(uint32(0)),
	reflect.Uint64:  reflect.TypeOf(uint64(0)),
	reflect.Float32: reflect.TypeOf(float32(0)),
	reflect.Float64: reflect.TypeOf(float64(0)),
}

// literal returns a number literal as the compiler pushes it: converted to the
// numeric type the checker gave the node, e.g. int64 for a func(int64)
// parameter. It returns false if the value does not fit, leaving the call to
// the compiler.
func literal(value any, typ reflect.Type) (any, bool) {
	if typ == nil {
		return value, true
	}
	t, ok := literalTypes[typ.Kind()]
	if !ok {
		return value, true
	}
	v := reflect.ValueOf(value)
	out := reflect.New(t).Elem()
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Kind() != reflect.Int || out.OverflowInt(v.Int()) {
			return nil, false
		}
		out.SetInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Kind() != reflect.Int || v.Int() < 0 || out.OverflowUint(uint64(v.Int())) {
			return nil, false
		}
		out.SetUint(uint64(v.Int()))
	case reflect.Float32, reflect.Float64:
		out.SetFloat(v.Convert(literalTypes[reflect.Float64]).Float())
	}
	return out.Interface(), true
}
