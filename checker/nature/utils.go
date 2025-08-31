package nature

import (
	"reflect"
)

func derefTypeKind(t reflect.Type, k reflect.Kind) (_ reflect.Type, _ reflect.Kind, changed bool) {
	for k == reflect.Pointer {
		changed = true
		t = t.Elem()
		k = t.Kind()
	}
	return t, k, changed
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
	cache                     *Cache
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

func (s *structData) structField(parentEmbed *structData, name string) (structField, bool) {
	if s.fields == nil {
		if s.numField > 0 {
			s.fields = make(map[string]structField, s.numField)
		}
	} else if f, ok := s.fields[name]; ok {
		return f, true
	}
	if s.finished() {
		return structField{}, false
	}

	// Lookup own fields first.
	for ; s.ownIdx < s.numField; s.ownIdx++ {
		field := s.rType.Field(s.ownIdx)
		// BUG: we should skip if !field.IsExported() here

		if field.Anonymous && s.anonIdx < 0 {
			// start iterating anon fields on the first instead of zero
			s.anonIdx = s.ownIdx
		}
		fName, ok := fieldName(field.Name, field.Tag)
		if !ok || fName == "" {
			// name can still be empty for a type created at runtime with
			// reflect
			continue
		}
		nt := s.cache.FromType(field.Type)
		sf := structField{
			Nature: nt,
			Index:  field.Index,
		}
		s.fields[fName] = sf
		if parentEmbed != nil {
			parentEmbed.trySet(fName, sf)
		}
		if fName == name {
			return sf, true
		}
	}

	if s.curChild != nil {
		sf, ok := s.findInEmbedded(parentEmbed, s.curChild, s.curChildIndex, name)
		if ok {
			return sf, true
		}
	}

	// Lookup embedded fields through anon own fields
	for ; s.anonIdx >= 0 && s.anonIdx < s.numField; s.anonIdx++ {
		field := s.rType.Field(s.anonIdx)
		// we do enter embedded non-exported types because they could contain
		// exported fields
		if !field.Anonymous {
			continue
		}
		t, k, _ := derefTypeKind(field.Type, field.Type.Kind())
		if k != reflect.Struct {
			continue
		}

		childEmbed := s.cache.getStruct(t).structData
		sf, ok := s.findInEmbedded(parentEmbed, childEmbed, field.Index, name)
		if ok {
			return sf, true
		}
	}

	return structField{}, false
}

func (s *structData) findInEmbedded(
	parentEmbed, childEmbed *structData,
	childIndex []int,
	name string,
) (structField, bool) {
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
		s.trySet(name, sf)
	}

	// Recheck if we have what we needed from the above sync
	if sf, ok := s.fields[name]; ok {
		return sf, true
	}

	// Try finding in the child again in case it hasn't finished
	if !childEmbed.finished() {
		if _, ok := childEmbed.structField(s, name); ok {
			return s.fields[name], true
		}
	}

	return structField{}, false
}

func (s *structData) trySet(name string, sf structField) {
	if _, ok := s.fields[name]; ok {
		return
	}
	sf.Index = append(s.curChildIndex, sf.Index...)
	s.fields[name] = structField{
		Nature: sf.Nature,
		Index:  sf.Index,
	}
	if s.curParent != nil {
		s.curParent.trySet(name, sf)
	}
}

func StructFields(c *Cache, t reflect.Type) map[string]Nature {
	table := make(map[string]Nature)
	if t == nil {
		return table
	}
	t, k, _ := derefTypeKind(t, t.Kind())
	if k == reflect.Struct {
		// lookup for a field with an empty name, which will cause to never find a
		// match, meaning everything will have been cached.
		sd := c.getStruct(t).structData
		sd.structField(nil, "")
		for name, sf := range sd.fields {
			table[name] = sf.Nature
		}
	}
	return table
}

type methodset struct {
	cache          *Cache
	rType          reflect.Type
	kind           reflect.Kind
	methods        map[string]method
	numMethod, idx int
}

type method struct {
	reflect.Method
	nature Nature
}

func (s *methodset) method(name string) (method, bool) {
	if s.methods == nil {
		s.methods = make(map[string]method, s.numMethod)
	} else if m, ok := s.methods[name]; ok {
		return m, true
	}
	for ; s.idx < s.numMethod; s.idx++ {
		rm := s.rType.Method(s.idx)
		if !rm.IsExported() {
			continue
		}
		nt := s.cache.FromType(rm.Type)
		if s.rType.Kind() != reflect.Interface {
			nt.Method = true
			nt.MethodIndex = rm.Index
			// In case of interface type method will not have a receiver,
			// and to prevent checker decreasing numbers of in arguments
			// return method type as not method (second argument is false).

			// Also, we can not use m.Index here, because it will be
			// different indexes for different types which implement
			// the same interface.
		}
		m := method{
			Method: rm,
			nature: nt,
		}
		s.methods[rm.Name] = m
		if rm.Name == name {
			return m, true
		}
	}
	return method{}, false
}
