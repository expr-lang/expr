package conf

import (
	"fmt"
	"reflect"

	. "expr/checker/nature"
	"expr/internal/deref"
	"expr/types"
)

func Env(env any) Nature {
	if env == nil {
		return Nature{
			Type:   reflect.TypeOf(map[string]any{}),
			Strict: true,
		}
	}

	switch env := env.(type) {
	case types.Map:
		return env.Nature()
	}

	v := reflect.ValueOf(env)
	d := deref.Value(v)

	switch d.Kind() {
	case reflect.Struct:
		return Nature{
			Type:   v.Type(),
			Strict: true,
		}

	case reflect.Map:
		n := Nature{
			Type:   v.Type(),
			Fields: make(map[string]Nature, v.Len()),
		}

		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}

			face := elem.Interface()

			switch face.(type) {
			case types.Map:
				n.Fields[key.String()] = face.(types.Map).Nature()

			default:
				if face == nil {
					n.Fields[key.String()] = Nature{Nil: true}
					continue
				}
				n.Fields[key.String()] = Nature{Type: reflect.TypeOf(face)}
			}

		}

		return n
	}

	panic(fmt.Sprintf("unknown type %T", env))
}
