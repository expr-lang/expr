package nature

import (
	"fmt"
	"reflect"
	"time"

	"github.com/expr-lang/expr/builtin"
	"github.com/expr-lang/expr/internal/deref"
)

var (
	intType      = reflect.TypeOf(0)
	floatType    = reflect.TypeOf(float64(0))
	arrayType    = reflect.TypeOf([]any{})
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
)

type NatureCheck int

const (
	_ NatureCheck = iota
	BoolCheck
	StringCheck
	IntegerCheck
	NumberCheck
	MapCheck
	ArrayCheck
	TimeCheck
	DurationCheck
)

type Nature struct {
	// The order of the fields matter, check alignment before making changes.

	Type reflect.Type // Type of the value. If nil, then value is unknown.
	Kind reflect.Kind // Kind of the value.

	*Cache
	*Optional
	Func *builtin.Function // Used to pass function type from callee to CallNode.

	// Ref is a reference used for multiple, disjoint purposes. When the Nature
	// is for a:
	//	- Predicate: then Ref is the nature of the Out of the predicate.
	//	- Array-like types: then Ref is the Elem nature of array type (usually Type is []any, but ArrayOf can be any nature).
	Ref *Nature

	Nil    bool // If value is nil.
	Strict bool // If map is types.StrictMap.
	Method bool // If value retrieved from method. Usually used to determine amount of in arguments.
}

type Optional struct {
	// struct-only data
	FieldIndex  []int // Index of field in type.
	MethodIndex int   // Index of method in type.

	// map-only data
	Fields          map[string]Nature // Fields of map type.
	DefaultMapValue *Nature           // Default value of map type.

	// func-only data
	inElem, outZero *Nature
}

// Cache is a shared cache of type information. It is only used in the stages
// where type information becomes relevant, so packages like ast, parser, types,
// and lexer do not need to use the cache because they don't need any service
// from the Nature type, they only describe. However, when receiving a Nature
// from one of those packages, the cache must be set immediately.
type Cache struct {
	methodByName map[rTypeWithKey]*Nature
	fieldByName  map[rTypeWithKey]*Nature
	get          map[rTypeWithKey]*Nature
}

type rTypeWithKey struct {
	t   reflect.Type
	key string
}

func NatureOf(c *Cache, i any) Nature {
	// reflect.TypeOf(nil) returns nil, but in FromType we want to differentiate
	// what nil means for us
	if i == nil {
		return Nature{Cache: c, Nil: true}
	}
	return FromType(c, reflect.TypeOf(i))
}

func FromType(c *Cache, t reflect.Type) Nature {
	if t != nil {
		k := t.Kind()
		var opt *Optional
		if k == reflect.Func {
			opt = new(Optional)
		}
		return Nature{Type: t, Kind: k, Optional: opt, Cache: c}
	}
	return Nature{Cache: c}
}

func ArrayFromType(c *Cache, t reflect.Type) Nature {
	elem := FromType(c, t)
	nt := FromType(c, arrayType)
	nt.Ref = &elem
	return nt
}

func (n *Nature) IsAny() bool {
	return n.Type != nil && n.Kind == reflect.Interface && n.NumMethods() == 0
}

func (n *Nature) IsUnknown() bool {
	return n.Type == nil && !n.Nil || n.IsAny()
}

func (n *Nature) String() string {
	if n.Type != nil {
		return n.Type.String()
	}
	return "unknown"
}

func (n *Nature) Deref() Nature {
	ret := *n
	if ret.Type != nil {
		ret.Type = deref.Type(ret.Type)
		ret.Kind = ret.Type.Kind()
	}
	return ret
}

func (n *Nature) Key() Nature {
	if n.Kind == reflect.Map {
		return FromType(n.Cache, n.Type.Key())
	}
	return FromType(n.Cache, nil)
}

func (n *Nature) Elem() Nature {
	switch n.Kind {
	case reflect.Ptr:
		return FromType(n.Cache, n.Type.Elem())
	case reflect.Map:
		if n.Optional != nil && n.DefaultMapValue != nil {
			return *n.DefaultMapValue
		}
		return FromType(n.Cache, n.Type.Elem())
	case reflect.Slice, reflect.Array:
		if n.Ref != nil {
			return *n.Ref
		}
		return FromType(n.Cache, n.Type.Elem())
	}
	return FromType(n.Cache, nil)
}

func (n *Nature) AssignableTo(nt Nature) bool {
	if n.Nil {
		// Untyped nil is assignable to any interface, but implements only the empty interface.
		if nt.IsAny() {
			return true
		}
	}
	if n.Type == nil || nt.Type == nil {
		return false
	}
	return n.Type.AssignableTo(nt.Type)
}

func (n *Nature) NumMethods() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumMethod()
}

func (n *Nature) MethodByName(name string) (Nature, bool) {
	if ntPtr := n.methodByNamePtr(name); ntPtr != nil {
		return *ntPtr, true
	}
	return FromType(n.Cache, nil), false
}

func (n *Nature) methodByNamePtr(name string) *Nature {
	return n.methodByNameSlow(name)
	var ntPtr *Nature
	var cacheHit bool
	if n.Cache.methodByName == nil {
		n.Cache.methodByName = map[rTypeWithKey]*Nature{}
	} else {
		ntPtr, cacheHit = n.Cache.methodByName[rTypeWithKey{n.Type, name}]
	}
	if !cacheHit {
		ntPtr = n.methodByNameSlow(name)
		n.Cache.methodByName[rTypeWithKey{n.Type, name}] = ntPtr
	}
	return ntPtr
}

func (n *Nature) methodByNameSlow(name string) *Nature {
	if n.Type == nil {
		return nil
	}
	method, ok := n.Type.MethodByName(name)
	if !ok {
		return nil
	}

	nt := FromType(n.Cache, method.Type)
	if n.Kind == reflect.Interface {
		// In case of interface type method will not have a receiver,
		// and to prevent checker decreasing numbers of in arguments
		// return method type as not method (second argument is false).

		// Also, we can not use m.Index here, because it will be
		// different indexes for different types which implement
		// the same interface.
		return &nt
	}
	if nt.Optional == nil {
		nt.Optional = new(Optional)
	}
	nt.Method = true
	nt.MethodIndex = method.Index
	return &nt
}

func (n *Nature) NumIn() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumIn()
}

func (n *Nature) InElem(i int) Nature {
	if n.inElem == nil {
		if n.Type == nil {
			n2 := FromType(n.Cache, nil)
			n.inElem = &n2
		} else {
			n2 := FromType(n.Cache, n.Type.In(i))
			n2 = n2.Elem()
			n.inElem = &n2
		}
	}
	return *n.inElem
}

func (n *Nature) In(i int) Nature {
	if n.Type == nil {
		return FromType(n.Cache, nil)
	}
	return FromType(n.Cache, n.Type.In(i))
}

func (n *Nature) IsFirstArgUnknown() bool {
	if n.Type != nil {
		n2 := FromType(n.Cache, n.Type.In(0))
		return n2.IsUnknown()
	}
	return false
}

func (n *Nature) NumOut() int {
	if n.Type == nil {
		return 0
	}
	return n.Type.NumOut()
}

func (n *Nature) Out(i int) Nature {
	if i != 0 {
		return n.out(i)
	}
	if n.outZero != nil {
		return *n.outZero
	}
	nt := n.out(0)
	n.outZero = &nt
	return nt
}

func (n *Nature) out(i int) Nature {
	if n.Type == nil {
		return FromType(n.Cache, nil)
	}
	return FromType(n.Cache, n.Type.Out(i))
}

func (n *Nature) IsVariadic() bool {
	if n.Type == nil {
		return false
	}
	return n.Type.IsVariadic()
}

func (n *Nature) FieldByName(name string) (Nature, bool) {
	var ntPtr *Nature
	var cacheHit bool
	if n.Cache.fieldByName == nil {
		n.Cache.fieldByName = map[rTypeWithKey]*Nature{}
	} else {
		ntPtr, cacheHit = n.Cache.fieldByName[rTypeWithKey{n.Type, name}]
	}
	if !cacheHit {
		ntPtr = n.fieldByNameSlow(name)
		n.Cache.fieldByName[rTypeWithKey{n.Type, name}] = ntPtr
	}
	if ntPtr != nil {
		return *ntPtr, true
	}
	return FromType(n.Cache, nil), false
}

func (n *Nature) fieldByNameSlow(name string) *Nature {
	if n.Type == nil {
		return nil
	}
	if field, ok := fetchField(n.Type, name); ok {
		nt := FromType(n.Cache, field.Type)
		if nt.Optional == nil {
			nt.Optional = new(Optional)
		}
		nt.FieldIndex = field.Index
		return &nt
	}
	return nil
}

func (n *Nature) PkgPath() string {
	if n.Type == nil {
		return ""
	}
	return n.Type.PkgPath()
}

func (n *Nature) IsFastMap() bool {
	if n.Type == nil {
		return false
	}
	if n.Type.Kind() == reflect.Map &&
		n.Type.Key().Kind() == reflect.String &&
		n.Type.Elem().Kind() == reflect.Interface {
		return true
	}
	return false
}

func (n *Nature) Get(name string) (Nature, bool) {
	var ntPtr *Nature
	var cacheHit bool
	if n.Cache.get == nil {
		n.Cache.get = map[rTypeWithKey]*Nature{}
	} else {
		ntPtr, cacheHit = n.Cache.get[rTypeWithKey{n.Type, name}]
	}
	if !cacheHit {
		ntPtr = n.getSlow(name)
		n.Cache.get[rTypeWithKey{n.Type, name}] = ntPtr
	}
	if ntPtr != nil {
		return *ntPtr, true
	}
	return FromType(n.Cache, nil), false
}

func (n *Nature) getSlow(name string) *Nature {
	if n.Type == nil {
		return nil
	}

	if m := n.methodByNamePtr(name); m != nil {
		return m
	}

	t := deref.Type(n.Type)
	switch t.Kind() {
	case reflect.Struct:
		if f, ok := fetchField(t, name); ok {
			nt := FromType(n.Cache, f.Type)
			if nt.Optional == nil {
				nt.Optional = new(Optional)
			}
			nt.FieldIndex = f.Index
			return &nt
		}
	case reflect.Map:
		if n.Optional != nil {
			if f, ok := n.Fields[name]; ok {
				return &f
			}
		}
	}
	return nil
}

func (n *Nature) All() map[string]Nature {
	table := make(map[string]Nature)

	if n.Type == nil {
		return table
	}

	for i := 0; i < n.NumMethods(); i++ {
		method := n.Type.Method(i)
		nt := FromType(n.Cache, method.Type)
		if nt.Optional == nil {
			nt.Optional = new(Optional)
		}
		nt.Method = true
		nt.MethodIndex = method.Index
		table[method.Name] = nt
	}

	t := deref.Type(n.Type)

	switch t.Kind() {
	case reflect.Struct:
		for name, nt := range StructFields(n.Cache, t) {
			if _, ok := table[name]; ok {
				continue
			}
			table[name] = nt
		}

	case reflect.Map:
		if n.Optional != nil {
			for key, nt := range n.Fields {
				if _, ok := table[key]; ok {
					continue
				}
				table[key] = nt
			}
		}
	}

	return table
}

func (n *Nature) IsNumber() bool {
	return n.IsInteger() || n.IsFloat()
}

func (n *Nature) IsInteger() bool {
	switch n.Kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return n.PkgPath() == ""
	}
	return false
}

func (n *Nature) IsFloat() bool {
	switch n.Kind {
	case reflect.Float32, reflect.Float64:
		return n.PkgPath() == ""
	}
	return false
}

func (n *Nature) PromoteNumericNature(rhs Nature) Nature {
	if n.IsUnknown() || rhs.IsUnknown() {
		return FromType(n.Cache, nil)
	}
	if n.IsFloat() || rhs.IsFloat() {
		return FromType(n.Cache, floatType)
	}
	return FromType(n.Cache, intType)
}

func (n *Nature) IsTime() bool {
	return n.Type == timeType
}

func (n *Nature) IsDuration() bool {
	return n.Type == durationType
}

func (n *Nature) IsBool() bool {
	return n.Kind == reflect.Bool
}

func (n *Nature) IsString() bool {
	return n.Kind == reflect.String
}

func (n *Nature) IsArray() bool {
	k := n.Kind
	return k == reflect.Slice || k == reflect.Array
}

func (n *Nature) IsMap() bool {
	return n.Kind == reflect.Map
}

func (n *Nature) IsStruct() bool {
	return n.Kind == reflect.Struct
}

func (n *Nature) IsFunc() bool {
	return n.Kind == reflect.Func
}

func (n *Nature) IsPointer() bool {
	return n.Kind == reflect.Ptr
}

func (n *Nature) IsAnyOf(cs ...NatureCheck) bool {
	var result bool
	for i := 0; i < len(cs) && !result; i++ {
		switch cs[i] {
		case BoolCheck:
			result = n.IsBool()
		case StringCheck:
			result = n.IsString()
		case IntegerCheck:
			result = n.IsInteger()
		case NumberCheck:
			result = n.IsNumber()
		case MapCheck:
			result = n.IsMap()
		case ArrayCheck:
			result = n.IsArray()
		case TimeCheck:
			result = n.IsTime()
		case DurationCheck:
			result = n.IsDuration()
		default:
			panic(fmt.Sprintf("unknown check value %d", cs[i]))
		}
	}
	return result
}

func (n *Nature) ComparableTo(rhs Nature) bool {
	return n.IsUnknown() || rhs.IsUnknown() ||
		n.Nil || rhs.Nil ||
		n.IsNumber() && rhs.IsNumber() ||
		n.IsDuration() && rhs.IsDuration() ||
		n.IsTime() && rhs.IsTime() ||
		n.IsArray() && rhs.IsArray() ||
		n.AssignableTo(rhs)
}

func (n *Nature) MaybeCompatible(rhs Nature, cs ...NatureCheck) bool {
	nIsUnknown := n.IsUnknown()
	rshIsUnknown := rhs.IsUnknown()
	return nIsUnknown && rshIsUnknown ||
		nIsUnknown && rhs.IsAnyOf(cs...) ||
		rshIsUnknown && n.IsAnyOf(cs...)
}

func (n *Nature) MakeArrayOf() Nature {
	nt := FromType(n.Cache, arrayType)
	nt.Ref = n
	return nt
}
