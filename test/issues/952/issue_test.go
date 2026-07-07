package issue952

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

// Node is an exported interface embedded into Wrapper below. Field access on the
// concrete type behind it is resolved by runtime.Fetch since #952.
type Node interface {
	ID() string
}

type Leaf struct {
	Name  string
	Value string
}

func (l Leaf) ID() string { return l.Name }

// Wrapper embeds the Node interface; the concrete value behind it (a Leaf) has
// a Value field that is not promoted by Go's standard field promotion.
type Wrapper struct {
	Node
}

// TestGetEmbeddedInterfaceField guards against a regression where the builtin
// get() resolved struct fields inconsistently with member access (a.field).
//
// #952 taught runtime.Fetch to resolve fields on the concrete type held by an
// embedded interface, but left the sibling builtin get() behind, so the same
// field access returned the value through member access yet nil through get().
// get() is documented to differ from runtime.Fetch only in returning nil
// instead of panicking, so the two must agree on which field is found.
func TestGetEmbeddedInterfaceField(t *testing.T) {
	w := Wrapper{Node: Leaf{Name: "n1", Value: "hello"}}
	env := map[string]any{"w": w}

	// Member access resolves the field via the concrete type behind the
	// embedded interface. The ternary keeps the type unknown at compile time
	// so the access goes through dynamic OpFetch (runtime.Fetch).
	out, err := expr.Eval(`(true ? w : "x").Value`, env)
	require.NoError(t, err)
	require.Equal(t, "hello", out)

	// get() must resolve the same field to the same value rather than nil.
	out, err = expr.Eval(`get((true ? w : "x"), "Value")`, env)
	require.NoError(t, err)
	require.Equal(t, "hello", out)

	// A genuinely missing field still yields nil from get() (not an error).
	out, err = expr.Eval(`get((true ? w : "x"), "Missing")`, env)
	require.NoError(t, err)
	require.Nil(t, out)
}
