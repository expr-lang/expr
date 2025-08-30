package nature

import (
	"reflect"

	"github.com/expr-lang/expr/internal/deref"
)

func derefTypeKind(t reflect.Type, k reflect.Kind) (reflect.Type, reflect.Kind) {
	for k == reflect.Pointer {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}

func fieldName(fieldName string, tag reflect.StructTag) (string, bool) {
	switch taggedName := tag.Get("expr"); taggedName {
	case "-":
		return "", false
	case "":
		return fieldName, true
	default:
		return taggedName, true
	}
}

func (c *Cache) fetchField(info map[string]Nature, t reflect.Type, name string) (Nature, bool) {
	numField := t.NumField()
	switch {
	case info != nil:
	case c.structs == nil:
		c.structs = map[reflect.Type]map[string]Nature{}
		fallthrough
	case c.structs[t] == nil:
		info = make(map[string]Nature, numField)
		c.structs[t] = info
	default:
		info = c.structs[t]
	}

	// Lookup own fields first. Cache all that is possible
	for i := 0; i < numField; i++ {
		sf := t.Field(i)
		// BUG: we should skip if !sf.IsExported()
		fName, ok := fieldName(sf.Name, sf.Tag)
		if !ok || fName == "" {
			// name can still be empty for a type created at runtime with
			// reflect
			continue
		}
		nt := c.FromType(sf.Type)
		if nt.Optional == nil {
			nt.Optional = new(Optional)
		}
		nt.FieldIndex = sf.Index
		if _, ok := info[fName]; !ok {
			// avoid overwriting fields that could potentially be own fields of
			// a parent struct
			info[fName] = nt
		}
		if fName == name {
			return nt, true
		}
	}

	// Lookup embedded fields
	for i := 0; i < numField; i++ {
		sf := t.Field(i)
		// we do enter embedded non-exported types because they could contain
		// exported fields
		if !sf.Anonymous {
			continue
		}
		t, k := derefTypeKind(sf.Type, sf.Type.Kind())
		if k != reflect.Struct {
			continue
		}
		nt, ok := c.fetchField(info, t, name)
		if ok {
			nt.FieldIndex = append(sf.Index, nt.FieldIndex...)
			return nt, true
		}
	}

	return c.FromType(nil), false
}

func StructFields(c *Cache, t reflect.Type) map[string]Nature {
	table := make(map[string]Nature)

	t = deref.Type(t)
	if t == nil {
		return table
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if f.Anonymous {
				for name, typ := range StructFields(c, f.Type) {
					if _, ok := table[name]; ok {
						continue
					}
					if typ.Optional == nil {
						typ.Optional = new(Optional)
					}
					typ.FieldIndex = append(f.Index, typ.FieldIndex...)
					table[name] = typ
				}
			}

			name, ok := fieldName(f.Name, f.Tag)
			if !ok {
				continue
			}
			nt := c.FromType(f.Type)
			if nt.Optional == nil {
				nt.Optional = new(Optional)
			}
			nt.FieldIndex = f.Index
			table[name] = nt

		}
	}

	return table
}
