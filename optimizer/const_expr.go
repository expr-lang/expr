package optimizer

import (
	"fmt"
	"reflect"
	"strings"

	. "github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/file"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type constExpr struct {
	applied bool
	err     error
	fns     map[string]reflect.Value
}

func (*constExpr) Enter(*Node) {}
func (c *constExpr) Exit(node *Node) {
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

	patch := func(newNode Node) {
		c.applied = true
		Patch(node, newNode)
	}

	switch n := (*node).(type) {
	case *FunctionNode:
		fn, ok := c.fns[n.Name]
		if ok {
			in := make([]reflect.Value, len(n.Arguments))
			for i := 0; i < len(n.Arguments); i++ {
				arg := n.Arguments[i]
				var param interface{}

				switch a := arg.(type) {
				case *NilNode:
					param = nil
				case *IntegerNode:
					param = a.Value
				case *FloatNode:
					param = a.Value
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
			patch(constNode)
		}
	}
}
