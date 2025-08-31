package nature

import (
	"fmt"
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

type structData struct {
	*Cache
	rType                     reflect.Type
	fields                    map[string]structField
	numField, ownIdx, anonIdx int

	curParent, curChild *structData
	curChildIndex       []int
}

type structField struct {
	Nature
	Index []int
}

func (s *structData) finished() bool {
	return s.ownIdx >= s.numField && // no own fields left to visit
		s.anonIdx >= s.numField && // no embedded fields to visit
		s.curChild == nil // no child in process of visiting
}

func (s *structData) structDebug(prefix string) { // TODO: DEBUG
	if len(s.fields) == 0 {
		fmt.Printf("%s[%s]\n", prefix, s.rType)
		return
	}
	fmt.Printf("%s[%s] fields:\n", prefix, s.rType)
	prefix += "  "
	for k, v := range s.fields {
		fmt.Printf("%s%s at index %v:", prefix, k, v.Index)
		if v.Nature.Kind == reflect.Struct {
			v.Nature.structDebug(prefix)
		} else {
			fmt.Printf("%s[%s]\n", prefix, v.Nature.Type)
		}
	}
}

func (s *structData) structField(parentEmbed *structData, name string) (Nature, bool) {
	if f, ok := s.fields[name]; ok {
		return f.Nature, true
	}
	if s.finished() {
		return Nature{}, false
	}

	// Lookup own fields first.
	for ; s.ownIdx < s.numField; s.ownIdx++ {
		sf := s.rType.Field(s.ownIdx)
		// BUG: we should skip if !sf.IsExported() here

		if sf.Anonymous && s.anonIdx < 0 {
			// start iterating anon fields on the first instead of zero
			s.anonIdx = s.ownIdx
		}
		fName, ok := fieldName(sf.Name, sf.Tag)
		if !ok || fName == "" {
			// name can still be empty for a type created at runtime with
			// reflect
			continue
		}
		nt := s.FromType(sf.Type)

		// TODO BEGIN: deprecate
		opt := new(Optional)
		if nt.Optional != nil {
			*opt = *nt.Optional
		}
		nt.Optional = opt
		nt.FieldIndex = sf.Index
		// TODO END: deprecate

		s.fields[fName] = structField{
			Nature: nt,
			Index:  sf.Index,
		}
		if parentEmbed != nil {
			parentEmbed.trySet(fName, nt, sf.Index)
		}
		if fName == name {
			return nt, true
		}
	}

	if s.curChild != nil {
		nt, ok := s.findInEmbedded(parentEmbed, s.curChild, s.curChildIndex, name)
		if ok {
			return nt, true
		}
	}

	// Lookup embedded fields through anon own fields
	for ; s.anonIdx >= 0 && s.anonIdx < s.numField; s.anonIdx++ {
		sf := s.rType.Field(s.anonIdx)
		// we do enter embedded non-exported types because they could contain
		// exported fields
		if !sf.Anonymous {
			continue
		}
		t, k := derefTypeKind(sf.Type, sf.Type.Kind())
		if k != reflect.Struct {
			continue
		}

		childEmbed := s.Cache.getStruct(t).structData
		nt, ok := s.findInEmbedded(parentEmbed, childEmbed, sf.Index, name)
		if ok {
			return nt, true
		}
	}

	return Nature{}, false
}

func (s *structData) findInEmbedded(
	parentEmbed, childEmbed *structData,
	childIndex []int,
	name string,
) (Nature, bool) {
	// Set current parent/child data. This allows trySet to handle child fields
	// and add them to our struct and to the parent as well if needed
	s.curParent = parentEmbed
	s.curChild = childEmbed
	s.curChildIndex = childIndex
	defer func() {
		// Ensure to cleanup references
		s.curParent = nil
		if childEmbed.finished() {
			// If the child can still have more fields to explore then keep it
			// referened to look it up again if we need to
			s.curChild = nil
			s.curChildIndex = nil
		}
	}()

	// See if the child has already cached its fields. This is still important
	// to check even if it's the s.unfinishedEmbedded because it may have
	// explored new fields since the last time we visited it
	for name, sf := range childEmbed.fields {
		s.trySet(name, sf.Nature, sf.Index)
	}

	// Recheck if we have what we needed from the above sync
	if sf, ok := s.fields[name]; ok {
		return sf.Nature, true
	}

	// Try finding in the child again in case it hasn't finished
	if !childEmbed.finished() {
		if _, ok := childEmbed.structField(s, name); ok {
			return s.fields[name].Nature, true
		}
	}

	return Nature{}, false
}

func (s *structData) trySet(name string, nt Nature, idx []int) {
	if _, ok := s.fields[name]; ok {
		return
	}
	idx = append(s.curChildIndex, idx...)

	// TODO BEGIN: deprecate
	opt := new(Optional)
	if nt.Optional != nil {
		*opt = *nt.Optional
	}
	nt.Optional = opt
	nt.FieldIndex = idx
	// TODO END: deprecate

	s.fields[name] = structField{
		Nature: nt,
		Index:  idx,
	}
	if s.curParent != nil {
		s.curParent.trySet(name, nt, idx)
	}
}

// TODO: deprecate
func (c *Cache) fetchField(info map[string]Nature, t reflect.Type, name string) (Nature, bool) {
	numField := t.NumField()
	switch {
	case info != nil:
	case c.xxxStructs == nil:
		c.xxxStructs = map[reflect.Type]map[string]Nature{}
		fallthrough
	case c.xxxStructs[t] == nil:
		info = make(map[string]Nature, numField)
		c.xxxStructs[t] = info
	default:
		info = c.xxxStructs[t]
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
