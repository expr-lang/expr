package runtime_test

import (
	"github.com/expr-lang/expr/vm/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEquals(t *testing.T) {
	t.Run("same type", func(t *testing.T) {
		assert.True(t, runtime.Equal(12, 12))
		assert.True(t, runtime.Equal(12.34, 12.34))
	})

	t.Run("different int type", func(t *testing.T) {
		assert.True(t, runtime.Equal(int64(12), int32(12)))
	})

	t.Run("different float type", func(t *testing.T) {
		assert.True(t, runtime.Equal(float64(12.34), float32(12.34)))
	})
}
