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

	cache *Cache
	*Optional
	*FuncData

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
	pkgPath   string
	methodset *methodset // optional to avoid the map in *Cache

	*structData

	// map-only data
	Fields          map[string]Nature // Fields of map type.
	DefaultMapValue *Nature           // Default value of map type.

	pkgPathSet bool
}

type FuncData struct {
	Func        *builtin.Function // Used to pass function type from callee to CallNode.
	MethodIndex int               // Index of method in type.

	inElem, outZero *Nature
	numIn, numOut   int

	isVariadic    bool
	isVariadicSet bool
	numInSet      bool
	numOutSet     bool
}

// Cache is a shared cache of type information. It is only used in the stages
// where type information becomes relevant, so packages like ast, parser, types,
// and lexer do not need to use the cache because they don't need any service
// from the Nature type, they only describe. However, when receiving a Nature
// from one of those packages, the cache must be set immediately.
type Cache struct {
	methods map[reflect.Type]*methodset
	structs map[reflect.Type]Nature
}

// NatureOf returns a Nature describing "i". If "i" is nil then it returns a
// Nature describing the value "nil".
func (c *Cache) NatureOf(i any) Nature {
	// reflect.TypeOf(nil) returns nil, but in FromType we want to differentiate
	// what nil means for us
	if i == nil {
		return Nature{cache: c, Nil: true}
	}
	return c.FromType(reflect.TypeOf(i))
}

// FromType returns a Nature describing a value of type "t". If "t" is nil then
// it returns a Nature describing an unknown value.
func (c *Cache) FromType(t reflect.Type) Nature {
	if t == nil {
		return Nature{}
	}
	var fd *FuncData
	k := t.Kind()
	switch k {
	case reflect.Struct:
		return c.getStruct(t)
	case reflect.Func:
		fd = new(FuncData)
	}
	return Nature{Type: t, Kind: k, FuncData: fd, cache: c}
}

func (c *Cache) getStruct(t reflect.Type) Nature {
	if c != nil {
		if c.structs == nil {
			c.structs = map[reflect.Type]Nature{}
		} else if nt, ok := c.structs[t]; ok {
			return nt
		}
	}
	nt := Nature{
		Type: t,
		Kind: reflect.Struct,
		Optional: &Optional{
			structData: &structData{
				cache:    c,
				rType:    t,
				numField: t.NumField(),
				anonIdx:  -1, // do not lookup embedded fields yet
			},
		},
	}
	if c != nil {
		nt.SetCache(c)
	}
	return nt
}

func (c *Cache) getMethodset(t reflect.Type, k reflect.Kind) *methodset {
	if t == nil || c == nil {
		return nil
	}
	if c.methods == nil {
		c.methods = map[reflect.Type]*methodset{
			t: nil,
		}
	} else if s, ok := c.methods[t]; ok {
		return s
	}
	numMethod := t.NumMethod()
	if numMethod < 1 {
		c.methods[t] = nil // negative cache
		return nil
	}
	s := &methodset{
		cache:     c,
		rType:     t,
		kind:      k,
		numMethod: numMethod,
	}
	c.methods[t] = s
	return s
}

// NatureOf calls NatureOf on a nil *Cache. See the comment on Cache.
func NatureOf(i any) Nature {
	var c *Cache
	return c.NatureOf(i)
}

// FromType calls FromType on a nil *Cache. See the comment on Cache.
func FromType(t reflect.Type) Nature {
	var c *Cache
	return c.FromType(t)
}

func ArrayFromType(c *Cache, t reflect.Type) Nature {
	elem := c.FromType(t)
	nt := c.FromType(arrayType)
	nt.Ref = &elem
	return nt
}

func (n *Nature) SetCache(c *Cache) {
	n.cache = c
	if n.Kind == reflect.Struct {
		n.structData.cache = c
		if c.structs == nil {
			c.structs = map[reflect.Type]Nature{
				n.Type: *n,
			}
		} else if nt, ok := c.structs[n.Type]; ok {
			// invalidate local, use shared from cache
			n.Optional.structData = nt.Optional.structData
		} else {
			c.structs[n.Type] = *n
		}
	}
	if n.Optional != nil {
		if s, ok := c.methods[n.Type]; ok {
			// invalidate local if set, use shared from cache
			n.Optional.methodset = s
		} else if n.Optional.methodset != nil {
			c.methods[n.Type] = n.Optional.methodset
		}
	}
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
	t, _, changed := derefTypeKind(n.Type, n.Kind)
	if !changed {
		return *n
	}
	return n.cache.FromType(t)
}

func (n *Nature) Key() Nature {
	if n.Kind == reflect.Map {
		return n.cache.FromType(n.Type.Key())
	}
	return Nature{}
}

func (n *Nature) Elem() Nature {
	switch n.Kind {
	case reflect.Ptr:
		return n.cache.FromType(n.Type.Elem())
	case reflect.Map:
		if n.Optional != nil && n.DefaultMapValue != nil {
			return *n.DefaultMapValue
		}
		return n.cache.FromType(n.Type.Elem())
	case reflect.Slice, reflect.Array:
		if n.Ref != nil {
			return *n.Ref
		}
		return n.cache.FromType(n.Type.Elem())
	}
	return Nature{}
}

func (n *Nature) AssignableTo(nt Nature) bool {
	if n.Nil {
		switch nt.Kind {
		case reflect.Pointer, reflect.Interface:
			return true
		}
	}
	if n.Type == nil || nt.Type == nil ||
		n.Kind != nt.Kind && nt.Kind != reflect.Interface {
		return false
	}
	return n.Type.AssignableTo(nt.Type)
}

func (n *Nature) getMethodset() *methodset {
	if n.Optional != nil && n.Optional.methodset != nil {
		return n.Optional.methodset
	}
	s := n.cache.getMethodset(n.Type, n.Kind)
	if n.Optional != nil {
		n.Optional.methodset = s // cache locally if possible
	}
	return s
}

func (n *Nature) NumMethods() int {
	if s := n.getMethodset(); s != nil {
		return s.numMethod
	}
	return 0
}

func (n *Nature) MethodByName(name string) (Nature, bool) {
	if s := n.getMethodset(); s != nil {
		if m, ok := s.method(name); ok {
			return m.nature, true
		}
	}
	return Nature{}, false
}

func (n *Nature) NumIn() int {
	if n.numInSet {
		return n.numIn
	}
	n.numInSet = true
	n.numIn = n.Type.NumIn()
	return n.numIn
}

func (n *Nature) InElem(i int) Nature {
	if n.inElem == nil {
		n2 := n.cache.FromType(n.Type.In(i))
		n2 = n2.Elem()
		n.inElem = &n2
	}
	return *n.inElem
}

func (n *Nature) In(i int) Nature {
	return n.cache.FromType(n.Type.In(i))
}

func (n *Nature) IsFirstArgUnknown() bool {
	if n.Type != nil {
		n2 := n.cache.FromType(n.Type.In(0))
		return n2.IsUnknown()
	}
	return false
}

func (n *Nature) NumOut() int {
	if n.numOutSet {
		return n.numOut
	}
	n.numOutSet = true
	n.numOut = n.Type.NumOut()
	return n.numOut
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
		return Nature{}
	}
	return n.cache.FromType(n.Type.Out(i))
}

func (n *Nature) IsVariadic() bool {
	if n.isVariadicSet {
		return n.isVariadic
	}
	n.isVariadicSet = true
	n.isVariadic = n.Type.IsVariadic()
	return n.isVariadic
}

func (n *Nature) FieldByName(name string) (Nature, bool) {
	if n.Kind != reflect.Struct {
		return Nature{}, false
	}
	var sd *structData
	if n.Optional != nil && n.structData != nil {
		sd = n.structData
	} else {
		sd = n.cache.getStruct(n.Type).structData
	}
	if sf, ok := sd.structField(nil, name); ok {
		return sf.Nature, true
	}
	return Nature{}, false
}

func (n *Nature) PkgPath() string {
	if n.Type == nil {
		return ""
	}
	if n.Optional != nil && n.Optional.pkgPathSet {
		return n.Optional.pkgPath
	}
	p := n.Type.PkgPath()
	if n.Optional != nil {
		n.Optional.pkgPathSet = true
		n.Optional.pkgPath = p
	}
	return p
}

func (n *Nature) IsFastMap() bool {
	return n.Type != nil &&
		n.Type.Kind() == reflect.Map &&
		n.Type.Key().Kind() == reflect.String &&
		n.Type.Elem().Kind() == reflect.Interface
}

func (n *Nature) Get(name string) (Nature, bool) {
	if n.Kind == reflect.Map && n.Optional != nil {
		f, ok := n.Fields[name]
		return f, ok
	}
	return n.getSlow(name)
}

func (n *Nature) getSlow(name string) (Nature, bool) {
	if nt, ok := n.MethodByName(name); ok {
		return nt, true
	}
	if n.Kind == reflect.Struct {
		if sf, ok := n.structField(nil, name); ok {
			return sf.Nature, true
		}
	}
	return Nature{}, false
}

func (n *Nature) FieldIndex(name string) ([]int, bool) {
	if n.Kind != reflect.Struct {
		return nil, false
	}
	if sf, ok := n.structField(nil, name); ok {
		return sf.Index, true
	}
	return nil, false
}

func (n *Nature) All() map[string]Nature {
	table := make(map[string]Nature)

	if n.Type == nil {
		return table
	}

	for i := 0; i < n.NumMethods(); i++ {
		method := n.Type.Method(i)
		nt := n.cache.FromType(method.Type)
		if nt.Optional == nil {
			nt.FuncData = new(FuncData)
		}
		nt.Method = true
		nt.MethodIndex = method.Index
		table[method.Name] = nt
	}

	t := deref.Type(n.Type)

	switch t.Kind() {
	case reflect.Struct:
		for name, nt := range StructFields(n.cache, t) {
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
		return Nature{}
	}
	if n.IsFloat() || rhs.IsFloat() {
		return n.cache.FromType(floatType)
	}
	return n.cache.FromType(intType)
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
	nt := n.cache.FromType(arrayType)
	nt.Ref = n
	return nt
}
