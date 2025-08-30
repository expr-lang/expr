package conf

import (
	"fmt"
	"reflect"

	. "github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/internal/deref"
	"github.com/expr-lang/expr/types"
)

func Env(c *Cache, env any) Nature {
	if env == nil {
		n := NatureOf(c, map[string]any{})
		n.Strict = true
		return n
	}

	switch env := env.(type) {
	case types.Map:
		nt := env.Nature()
		nt.Cache = c
		return nt
	}

	v := reflect.ValueOf(env)
	t := v.Type()

	switch deref.Value(v).Kind() {
	case reflect.Struct:
		n := FromType(c, t)
		n.Strict = true
		return n

	case reflect.Map:
		n := FromType(c, v.Type())
		if n.Optional == nil {
			n.Optional = new(Optional)
		}
		n.Strict = true
		n.Fields = make(map[string]Nature, v.Len())

		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}

			face := elem.Interface()

			switch face := face.(type) {
			case types.Map:
				nt := face.Nature()
				nt.Cache = c
				n.Fields[key.String()] = nt

			default:
				if face == nil {
					n.Fields[key.String()] = NatureOf(c, nil)
					continue
				}
				n.Fields[key.String()] = NatureOf(c, face)
			}

		}

		return n
	}

	panic(fmt.Sprintf("unknown type %T", env))
}
