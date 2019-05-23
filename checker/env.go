package checker

import "reflect"

type tag struct {
	reflect.Type
	method bool
}

type typesTable map[string]tag

// OptionFn for configuring checker.
type OptionFn func(v *visitor)

// Define sets variable for type checks during parsing.
func Define(name string, t interface{}) OptionFn {
	return func(v *visitor) {
		v.types[name] = tag{Type: reflect.TypeOf(t)}
	}
}

// Env sets variables for type checks during parsing.
// If struct is passed, all fields will be treated as variables,
// as well as all fields of embedded structs and struct itself.
//
// If map is passed, all items will be treated as variables
// (key as name, value as type).
func Env(i interface{}) OptionFn {
	return func(visitor *visitor) {
		for k, v := range createTypesTable(i) {
			visitor.types[k] = v
		}
	}
}

func createTypesTable(i interface{}) typesTable {
	types := make(typesTable)
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
			types[m.Name] = tag{Type: m.Type, method: true}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if key.Kind() == reflect.String && value.IsValid() && value.CanInterface() {
				types[key.String()] = tag{Type: reflect.TypeOf(value.Interface())}
			}
		}
	}

	return types
}

func fieldsFromStruct(t reflect.Type) typesTable {
	types := make(typesTable)
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

			types[f.Name] = tag{Type: f.Type}
		}
	}

	return types
}
