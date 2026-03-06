package expr_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// --- test types ---

type jsonTagged struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hidden    string `json:"hidden,omitempty"`
	Ignored   string `json:"-"`
	NoTag     string
	Embed     *jsonEmbedded `json:"embed,omitempty"`
}

type embeddedBase struct {
	BaseField string `json:"base_field"`
}

type jsonEmbedded struct {
	embeddedBase
	Own string `json:"own"`
}

// TestWithTag_BasicJSON checks that json-tagged fields are accessible by their tag name.
func TestWithTag_BasicJSON(t *testing.T) {
	env := jsonTagged{FirstName: "John", LastName: "Doe"}
	program, err := expr.Compile(`first_name + " " + last_name`,
		expr.Env(jsonTagged{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", out)
}

// TestWithTag_HideField checks that json:"-" causes a compile-time error when accessing the field.
func TestWithTag_HideField(t *testing.T) {
	_, err := expr.Compile(`Ignored`,
		expr.Env(jsonTagged{}),
		expr.WithTag("json"),
	)
	require.Error(t, err)
}

// TestWithTag_CommaStripped checks that "name,omitempty" is accessible as "name".
func TestWithTag_CommaStripped(t *testing.T) {
	env := jsonTagged{Hidden: "secret"}
	program, err := expr.Compile(`hidden`,
		expr.Env(jsonTagged{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "secret", out)
}

// TestWithTag_CommaStrippedEmbed checks that a field whose tag has comma options
// (e.g. json:"embed,omitempty") is accessible at runtime via its tag name.
func TestWithTag_CommaStrippedEmbed(t *testing.T) {
	env := jsonTagged{}
	program, err := expr.Compile(`(embed?.own ?? "") == "foo"`,
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.False(t, out.(bool))
}

// TestWithTag_ExprTagInactive verifies that the expr tag is NOT resolved when using WithTag("json").
func TestWithTag_ExprTagInactive(t *testing.T) {
	type withExprTag struct {
		Val string `expr:"renamed" json:"val"`
	}
	// With json tag, "val" should work, "renamed" should not
	_, err := expr.Compile(`renamed`,
		expr.Env(withExprTag{}),
		expr.WithTag("json"),
	)
	require.Error(t, err)

	program, err := expr.Compile(`val`,
		expr.Env(withExprTag{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, withExprTag{Val: "ok"})
	require.NoError(t, err)
	assert.Equal(t, "ok", out)
}

// TestWithTag_FallbackToFieldName checks that a field with no json tag is still accessible by its Go name.
func TestWithTag_FallbackToFieldName(t *testing.T) {
	env := jsonTagged{NoTag: "direct"}
	program, err := expr.Compile(`NoTag`,
		expr.Env(jsonTagged{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "direct", out)
}

// TestWithTag_OrderingBeforeEnv checks that WithTag before Env works.
func TestWithTag_OrderingBeforeEnv(t *testing.T) {
	env := jsonTagged{FirstName: "Jane"}
	program, err := expr.Compile(`first_name`,
		expr.WithTag("json"),
		expr.Env(jsonTagged{}),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "Jane", out)
}

// TestWithTag_OrderingAfterEnv checks that WithTag after Env works.
func TestWithTag_OrderingAfterEnv(t *testing.T) {
	env := jsonTagged{FirstName: "Jane"}
	program, err := expr.Compile(`first_name`,
		expr.Env(jsonTagged{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "Jane", out)
}

// TestWithTag_InOperator checks that the `in` operator respects the configured tag.
func TestWithTag_InOperator(t *testing.T) {
	env := jsonTagged{FirstName: "John"}

	// "first_name" in env should be true (field exists and is not hidden)
	program, err := expr.Compile(`"first_name" in env`,
		expr.Env(map[string]any{"env": jsonTagged{}}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, map[string]any{"env": env})
	require.NoError(t, err)
	assert.Equal(t, true, out)
}

// TestWithTag_Embedded checks that tag resolution works for embedded structs.
func TestWithTag_Embedded(t *testing.T) {
	env := jsonEmbedded{
		embeddedBase: embeddedBase{BaseField: "base"},
		Own:          "own",
	}
	program, err := expr.Compile(`base_field + "-" + own`,
		expr.Env(jsonEmbedded{}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "base-own", out)
}

// TestWithTag_GetBuiltin checks that get() uses the configured tag.
func TestWithTag_GetBuiltin(t *testing.T) {
	env := map[string]any{"s": jsonTagged{FirstName: "Alice"}}
	program, err := expr.Compile(`get(s, "first_name")`,
		expr.Env(map[string]any{"s": jsonTagged{}}),
		expr.WithTag("json"),
	)
	require.NoError(t, err)
	out, err := expr.Run(program, env)
	require.NoError(t, err)
	assert.Equal(t, "Alice", out)
}

// TestWithTag_DefaultUnaffected checks that Eval() without WithTag still uses the "expr" tag.
func TestWithTag_DefaultUnaffected(t *testing.T) {
	type withExprTag struct {
		Val string `expr:"my_val"`
	}
	out, err := expr.Eval(`my_val`, withExprTag{Val: "hello"})
	require.NoError(t, err)
	assert.Equal(t, "hello", out)
}
