package nature

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr/types"
)

func Of(value any) Nature {
	if value == nil {
		return Nature{Nil: true}
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Map:
		_, strict := value.(types.StrictMap)
		subMap := Map{
			Elem:   Nature{Type: v.Type().Elem()},
			Fields: make(map[string]Nature, v.Len()),
			Strict: strict,
		}
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}
			face := elem.Interface()
			switch face.(type) {
			case types.Map, types.StrictMap:
				subMap.Fields[key.String()] = Of(face)
			default:
				if face == nil {
					subMap.Fields[key.String()] = Nature{Nil: true}
					continue
				}
				subMap.Fields[key.String()] = Nature{Type: reflect.TypeOf(face)}

			}
		}
		return Nature{
			Type:    v.Type(),
			SubType: subMap,
		}
	}

	return Nature{Type: v.Type()}
}

type SubType interface {
	String() string
	Get(name string) (Nature, bool)
}

type Array struct {
	Elem Nature
}

func (a Array) String() string {
	return "[]" + a.Elem.String()
}

func (a Array) Get(name string) (Nature, bool) {
	return unknown, false
}

type Map struct {
	Elem   Nature
	Fields map[string]Nature
	Strict bool
}

func (m Map) String() string {
	return "map[string]" + m.Elem.String()
}

func (m Map) Get(name string) (Nature, bool) {
	if n, ok := m.Fields[name]; ok {
		return n, true
	}
	return unknown, false
}
