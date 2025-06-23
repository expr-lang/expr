package runtime_test

import (
	"testing"

	"github.com/expr-lang/expr"
)

func TestCanSliceNonAddressableArrayType(t *testing.T) {
	prog, err := expr.Compile(`arr[1:2][0]`)
	if err != nil {
		t.Fatalf("error compiling program: %v", err)
	}
	val, err := expr.Run(prog, map[string]any{
		"arr": [5]int{0, 1, 2, 3, 4},
	})
	if err != nil {
		t.Fatalf("error running program: %v", err)
	}
	valInt, ok := val.(int)
	if !ok || valInt != 1 {
		t.Fatalf("invalid result, expected 1, got %v", val)
	}
}
