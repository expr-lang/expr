package checker

import "reflect"

type TypesTable map[string]reflect.Type

// OptionFn for configuring checker.
type OptionFn func(v *visitor)

// Define sets variable for type checks during parsing.
func Define(name string, t interface{}) OptionFn {
	return func(v *visitor) {
		v.types[name] = reflect.TypeOf(t)
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

func createTypesTable(i interface{}) TypesTable {
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
			types[m.Name] = m.Type
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if key.Kind() == reflect.String && value.IsValid() && value.CanInterface() {
				types[key.String()] = reflect.TypeOf(value.Interface())
			}
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

			types[f.Name] = f.Type
		}
	}

	return types
}
