// Package fields_test checks that struct fields are named the same way
// everywhere in expr: by the type checker (typed programs) and at runtime
// (untyped programs, the "in" operator, and the builtin get()).
//
// A field is named by its expr tag, or by its Go name if it has no tag.
// Unexported fields and fields tagged expr:"-" have no name. A struct's own
// fields come before the fields of its embedded structs, which are visited
// depth first in declaration order, and the first field found for a name wins.
package fields_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

// Shapes used by the tests. Every int field gets a distinct value, so the
// value read tells which field a name resolved to.

type Issue990 struct {
	Inner int `expr:"bar"`
}

type Hidden struct {
	Plain   int
	Hidden  int `expr:"-"`
	private int
}

type Swap struct {
	A int `expr:"B"`
	B int `expr:"A"`
}

type UnexportedClash struct {
	Inner int `expr:"bar"`
	bar   int
}

type TagTakesGoName struct {
	X int `expr:"Y"`
	Y int
}

type Emb struct {
	Deep int
	Own  int
}

type EmbTagged struct {
	Deep int `expr:"deep"`
}

type Shadow struct {
	Emb
	Own int
}

type EmbPtr struct {
	*Emb
}

type EmbFieldTagged struct {
	EmbTagged
}

type EmbStructTagged struct {
	Emb `expr:"e"`
}

type EmbStructHidden struct {
	Emb `expr:"-"`
}

type emb struct {
	Deep int
}

type EmbUnexported struct {
	emb
}

type C struct{ X int }
type DepthA struct{ C }
type DepthB struct{ X int }

// Depth first: A.C.X is found before B.X, although B.X is shallower.
type DepthConflict struct {
	DepthA
	DepthB
}

type SameA struct{ X int }
type SameB struct{ X int }

// Same depth: the first embedded struct wins.
type SameDepth struct {
	SameA
	SameB
}

type ShadowByTag struct {
	Emb
	Y int `expr:"Deep"`
}

type InnerByGoName struct {
	Inner int `expr:"bar"`
	InnerEmb
}

type InnerEmb struct {
	Inner int
}

type Recursive struct {
	*Recursive
	X int
}

type CycleA struct {
	*CycleB
	X int
}

type CycleB struct {
	*CycleA
	Y int `expr:"y"`
}

type DiamondA struct{ C }
type DiamondB struct{ C }

type Diamond struct {
	DiamondA
	DiamondB
}

type Pointers struct {
	P *Emb `expr:"p"`
	Q *Emb
}

var shapes = []any{
	Issue990{},
	Hidden{},
	Swap{},
	UnexportedClash{},
	TagTakesGoName{},
	Shadow{},
	EmbPtr{},
	EmbFieldTagged{},
	EmbStructTagged{},
	EmbStructHidden{},
	EmbUnexported{},
	DepthConflict{},
	SameDepth{},
	ShadowByTag{},
	InnerByGoName{},
	Recursive{},
	CycleA{},
	Diamond{},
	Pointers{},
}

// TestFieldNames pins the naming rules: for each shape, which names resolve
// and to which field's value.
func TestFieldNames(t *testing.T) {
	tests := []struct {
		shape any
		names map[string]any // name -> value read; nil means no such field
	}{
		{Issue990{}, map[string]any{"bar": 1, "Inner": nil}},
		{Hidden{}, map[string]any{"Plain": 1, "Hidden": nil, "private": nil}},
		{Swap{}, map[string]any{"A": 2, "B": 1}},
		{UnexportedClash{}, map[string]any{"bar": 1, "Inner": nil}},
		{TagTakesGoName{}, map[string]any{"Y": 1, "X": nil}},
		{Shadow{}, map[string]any{"Deep": 1, "Own": 3, "Emb": Emb{Deep: 1, Own: 2}}},
		{EmbPtr{}, map[string]any{"Deep": 1, "Own": 2}},
		{EmbFieldTagged{}, map[string]any{"deep": 1, "Deep": nil}},
		{EmbStructTagged{}, map[string]any{"e": Emb{Deep: 1, Own: 2}, "Emb": nil, "Deep": 1}},
		{EmbStructHidden{}, map[string]any{"Emb": nil, "Deep": 1}},
		{EmbUnexported{}, map[string]any{"Deep": 1, "emb": nil}},
		{DepthConflict{}, map[string]any{"X": 1}},
		{SameDepth{}, map[string]any{"X": 1}},
		{ShadowByTag{}, map[string]any{"Deep": 3, "Y": nil}},
		{InnerByGoName{}, map[string]any{"bar": 1, "Inner": 2}},
		{Recursive{}, map[string]any{"X": 2}},
		{CycleA{}, map[string]any{"X": 2, "y": 1, "Y": nil}},
		{Diamond{}, map[string]any{"X": 1}},
	}
	for _, tt := range tests {
		for name, want := range tt.names {
			t.Run(fmt.Sprintf("%T/%s", tt.shape, name), func(t *testing.T) {
				env := map[string]any{"s": fill(tt.shape, false)}
				code := fmt.Sprintf("s[%q]", name)

				out, err := expr.Eval(code, env)
				if want == nil {
					require.Error(t, err)
					require.Contains(t, err.Error(), "cannot fetch")
				} else {
					require.NoError(t, err)
					require.Equal(t, want, out)
				}

				in, err := expr.Eval(fmt.Sprintf("%q in s", name), env)
				require.NoError(t, err)
				require.Equal(t, want != nil, in)
			})
		}
	}
}

// TestTypedMatchesUntyped checks, for every shape and every name that occurs
// in it, that typed and untyped programs agree on field access, "in", get()
// and $env, on values and on pointers to them.
func TestTypedMatchesUntyped(t *testing.T) {
	for _, shape := range shapes {
		for _, nilEmbedded := range []bool{false, true} {
			for _, ptr := range []bool{false, true} {
				value := fill(shape, nilEmbedded)
				if ptr {
					p := reflect.New(reflect.TypeOf(value))
					p.Elem().Set(reflect.ValueOf(value))
					value = p.Interface()
				}
				for _, name := range names(reflect.TypeOf(shape)) {
					t.Run(fmt.Sprintf("%T/nil=%v/ptr=%v/%q", shape, nilEmbedded, ptr, name), func(t *testing.T) {
						checkName(t, value, name)
					})
				}
			}
		}
	}
}

func checkName(t *testing.T, value any, name string) {
	quoted := strconv.Quote(name)
	env := map[string]any{"s": value}

	// Typed program: the checker decides whether the field exists.
	program, err := expr.Compile("s["+quoted+"]", expr.Env(env))
	exists := err == nil
	var typedOut any
	var typedErr error
	if exists {
		typedOut, typedErr = expr.Run(program, env)
	}

	// Field access, typed and untyped, as member and as top-level name.
	runs := []func() (any, error){
		func() (any, error) { return expr.Eval("s["+quoted+"]", env) },
		func() (any, error) { return expr.Eval("$env["+quoted+"]", value) },
		func() (any, error) { return evalTyped("$env["+quoted+"]", value) },
	}
	if isIdent(name) {
		runs = append(runs,
			func() (any, error) { return expr.Eval("s."+name, env) },
			func() (any, error) { return evalTyped("s."+name, env) },
			func() (any, error) { return expr.Eval(name, value) },
			func() (any, error) { return evalTyped(name, value) },
		)
	}
	for i, run := range runs {
		out, err := run()
		t.Logf("run %d: %v, %v", i, out, err)
		if typedErr != nil || !exists {
			require.Error(t, err)
			require.False(t, strings.Contains(err.Error(), "reflect"), err.Error())
			require.False(t, strings.Contains(err.Error(), "index out of range"), err.Error())
			continue
		}
		require.NoError(t, err)
		require.Equal(t, typedOut, out)
	}

	// "in" asks whether the field exists, typed or untyped.
	for _, code := range []string{quoted + " in s"} {
		out, err := expr.Eval(code, env)
		require.NoError(t, err)
		require.Equal(t, exists, out, code)

		out, err = evalTyped(code, env)
		require.NoError(t, err)
		require.Equal(t, exists, out, code)
	}
	out, err := expr.Eval(quoted+" in $env", value)
	require.NoError(t, err)
	require.Equal(t, exists, out)

	// get() returns the field, or nil if it can't be read.
	out, err = expr.Eval("get(s, "+quoted+")", env)
	require.NoError(t, err)
	if typedErr == nil && exists {
		require.Equal(t, typedOut, out)
	} else {
		require.Nil(t, out)
	}
}

func evalTyped(code string, env any) (any, error) {
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		return nil, err
	}
	return expr.Run(program, env)
}

var identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func isIdent(name string) bool {
	switch name {
	case "in", "not", "and", "or", "let", "nil", "true", "false", "if", "else", "matches", "contains":
		return false
	}
	return identRe.MatchString(name)
}

// names lists every Go name and expr tag in t and its embedded and nested
// structs, plus names that never resolve.
func names(t reflect.Type) []string {
	seen := map[string]bool{"missing": true, "": true, "-": true}
	visited := map[reflect.Type]bool{}
	var walk func(t reflect.Type)
	walk = func(t reflect.Type) {
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct || visited[t] {
			return
		}
		visited[t] = true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			seen[f.Name] = true
			if tag := f.Tag.Get("expr"); tag != "" {
				seen[tag] = true
			}
			walk(f.Type)
		}
	}
	walk(t)
	var list []string
	for name := range seen {
		list = append(list, name)
	}
	return list
}

// fill returns a copy of shape with a distinct value in every exported int
// field, in declaration order (depth first), and with embedded pointers
// allocated one level deep unless nilEmbedded is set.
func fill(shape any, nilEmbedded bool) any {
	v := reflect.New(reflect.TypeOf(shape)).Elem()
	n := 0
	var walk func(v reflect.Value, depth int)
	walk = func(v reflect.Value, depth int) {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fv := v.Field(i)
			switch {
			case f.Type.Kind() == reflect.Int:
				n++
				if fv.CanSet() {
					fv.SetInt(int64(n))
				}
			case f.Type.Kind() == reflect.Struct:
				walk(fv, depth)
			case f.Type.Kind() == reflect.Ptr && f.Anonymous && fv.CanSet() && !nilEmbedded && depth < 1:
				fv.Set(reflect.New(f.Type.Elem()))
				walk(fv.Elem(), depth+1)
			}
		}
	}
	walk(v, 0)
	return v.Interface()
}
