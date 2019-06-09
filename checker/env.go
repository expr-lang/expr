package checker

import "reflect"

type Tag struct {
	Type   reflect.Type
	method bool
}

type TypesTable map[string]Tag

// Env creates types table for type checks during parsing.
// If struct is passed, all fields will be treated as variables,
// as well as all fields of embedded structs and struct itself.
//
// If map is passed, all items will be treated as variables
// (key as name, value as type).
func Env(i interface{}) TypesTable {
	types := make(TypesTable)
	for k, v := range CreateTypesTable(i) {
		types[k] = v
	}
	return types
}

func CreateTypesTable(i interface{}) TypesTable {
	if i == nil {
		return nil
	}

	types := make(TypesTable)
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)

	d := t
	if t.Kind() == reflect.Ptr {
		d = t.Elem()
	}

	switch d.Kind() {
	case reflect.Struct:
		types = fieldsFromStruct(d)

		// Methods of struct should be gathered from original struct with pointer,
		// as methods maybe declared on pointer receiver. Also this method retrieves
		// all embedded structs methods as well, no need to recursion.
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			types[m.Name] = Tag{Type: m.Type, method: true}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if key.Kind() == reflect.String && value.IsValid() && value.CanInterface() {
				types[key.String()] = Tag{Type: reflect.TypeOf(value.Interface())}
			}
		}

		// A map may have method too.
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			types[m.Name] = Tag{Type: m.Type, method: true}
		}
	}

	return types
}

func fieldsFromStruct(t reflect.Type) TypesTable {
	types := make(TypesTable)
	t = dereference(t)
	if t == nil {
		return types
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if f.Anonymous {
				for name, typ := range fieldsFromStruct(f.Type) {
					types[name] = typ
				}
			}

			types[f.Name] = Tag{Type: f.Type}
		}
	}

	return types
}
