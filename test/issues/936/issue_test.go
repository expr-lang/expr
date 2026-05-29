package issue936

import (
	"reflect"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

// TestIssue936 tests that dynamic struct types created with reflect.StructOf
// compile and evaluate correctly even when fields have lowercase names (which
// require PkgPath to be set, making them appear "unexported" to reflect).
func TestIssue936(t *testing.T) {
	dynType := reflect.StructOf([]reflect.StructField{
		{
			Name:    "value",
			Type:    reflect.TypeFor[bool](),
			PkgPath: "github.com/some/package",
		},
	})
	env := reflect.New(dynType).Elem().Interface()

	// Compilation should succeed.
	program, err := expr.Compile("value", expr.Env(env))
	require.NoError(t, err)

	// Evaluation should also succeed and return the zero value (false).
	result, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, false, result)
}

// TestIssue936MultipleFields tests a dynamic struct with multiple field types.
func TestIssue936MultipleFields(t *testing.T) {
	dynType := reflect.StructOf([]reflect.StructField{
		{
			Name:    "name",
			Type:    reflect.TypeFor[string](),
			PkgPath: "github.com/some/package",
		},
		{
			Name:    "count",
			Type:    reflect.TypeFor[int](),
			PkgPath: "github.com/some/package",
		},
		{
			Name:    "active",
			Type:    reflect.TypeFor[bool](),
			PkgPath: "github.com/some/package",
		},
	})
	env := reflect.New(dynType).Elem().Interface()

	for _, tc := range []struct {
		expr string
	}{
		{`name == ""`},
		{`count == 0`},
		{`active == false`},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			program, err := expr.Compile(tc.expr, expr.Env(env))
			require.NoError(t, err)

			result, err := expr.Run(program, env)
			require.NoError(t, err)
			require.Equal(t, true, result)
		})
	}
}
