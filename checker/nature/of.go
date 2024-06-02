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
		fields := make(map[string]Nature, v.Len())
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}
			face := elem.Interface()
			switch face.(type) {
			case types.Map, types.StrictMap:
				fields[key.String()] = Of(face)
			default:
				if face == nil {
					fields[key.String()] = Nature{Nil: true}
					continue
				}
				fields[key.String()] = Nature{Type: reflect.TypeOf(face)}

			}
		}
		return Nature{
			Type:   v.Type(),
			Fields: fields,
			Strict: strict,
		}
	}

	return Nature{Type: v.Type()}
}
