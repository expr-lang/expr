package docgen

import (
	"github.com/antonmedv/expr/internal/conf"
	"reflect"
	"regexp"
	"strings"
)

// Kind can be any of array, map, struct, func, string, int, float, bool or any.
type Kind string

// Identifier represents variable names and field names.
type Identifier string

// TypeName is a name of type in types map.
type TypeName string

type Context struct {
	Variables map[Identifier]*Type `json:"variables"`
	Types     map[TypeName]*Type   `json:"types"`
}

type Type struct {
	Name      TypeName             `json:"name,omitempty"`
	Kind      Kind                 `json:"kind,omitempty"`
	Type      *Type                `json:"type,omitempty"`
	Key       *Type                `json:"key_type,omitempty"`
	Fields    map[Identifier]*Type `json:"fields,omitempty"`
	Arguments []*Type              `json:"arguments,omitempty"`
	Return    *Type                `json:"return,omitempty"`
}

func CreateDoc(i interface{}) *Context {
	c := &Context{
		Variables: make(map[Identifier]*Type),
		Types:     make(map[TypeName]*Type),
	}

	for name, t := range conf.CreateTypesTable(i) {
		c.Variables[Identifier(name)] = c.use(t.Type, fromMethod(t.Method))
	}

	return c
}

type config struct {
	method bool
}

type option func(c *config)

func fromMethod(b bool) option {
	return func(c *config) {
		c.method = b
	}
}

func (c *Context) use(t reflect.Type, ops ...option) *Type {
	config := &config{}
	for _, op := range ops {
		op(config)
	}

	methods := make([]reflect.Method, 0)

	// Methods of struct should be gathered from original struct with pointer,
	// as methods maybe declared on pointer receiver. Also this method retrieves
	// all embedded structs methods as well, no need to recursion.
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if isPrivate(m.Name) || isProtobuf(m.Name) {
			continue
		}
		methods = append(methods, m)
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only named types will have methods defined on them.
	// It maybe not even struct, but we gonna call then
	// structs in appendix anyway.
	if len(methods) > 0 {
		goto appendix
	}

	// This switch only for "simple" types.
	switch t.Kind() {
	case reflect.Bool:
		return &Type{Kind: "bool"}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Type{Kind: "int"}

	case reflect.Float32, reflect.Float64:
		return &Type{Kind: "float"}

	case reflect.String:
		return &Type{Kind: "string"}

	case reflect.Interface:
		return &Type{Kind: "any"}

	case reflect.Array, reflect.Slice:
		return &Type{
			Kind: "array",
			Type: c.use(t.Elem()),
		}

	case reflect.Map:
		return &Type{
			Kind: "map",
			Key:  c.use(t.Key()),
			Type: c.use(t.Elem()),
		}

	case reflect.Struct:
		goto appendix

	case reflect.Func:
		arguments := make([]*Type, 0)
		start := 0
		if config.method {
			start = 1
		}
		for i := start; i < t.NumIn(); i++ {
			arguments = append(arguments, c.use(t.In(i)))
		}
		f := &Type{
			Kind:      "func",
			Arguments: arguments,
		}
		if t.NumOut() > 0 {
			f.Return = c.use(t.Out(0))
		}
		return f
	}

appendix:
	name := TypeName(t.Name())
	anonymous := name == ""

	a, ok := c.Types[name]

	if !ok {
		a = &Type{
			Kind:   "struct",
			Fields: make(map[Identifier]*Type),
		}

		// baseNode a should be saved before starting recursion, or it will never end.
		if !anonymous {
			c.Types[name] = a
		}

		for name, field := range conf.FieldsFromStruct(t) {
			if isPrivate(name) || isProtobuf(name) {
				continue
			}
			a.Fields[Identifier(name)] = c.use(field.Type)
		}

		for _, m := range methods {
			if isPrivate(m.Name) || isProtobuf(m.Name) {
				continue
			}
			a.Fields[Identifier(m.Name)] = c.use(m.Type, fromMethod(true))
		}
	}

	if anonymous {
		return a
	}

	return &Type{
		Kind: "struct",
		Name: name,
	}
}

var isCapital = regexp.MustCompile("^[A-Z]")

func isPrivate(s string) bool {
	return !isCapital.Match([]byte(s))
}

func isProtobuf(s string) bool {
	return strings.HasPrefix(s, "XXX_")
}
