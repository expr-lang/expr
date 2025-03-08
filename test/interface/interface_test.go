package interface_test

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/assert"
)

type StoreInterface interface {
	Get(string) int
}

type StoreImpt struct{}

func (f StoreImpt) Get(s string) int {
	return 42
}

func (f StoreImpt) Set(s string, i int) bool {
	return true
}

type Env struct {
	Store StoreInterface `expr:"store"`
}

func TestInterfaceHide(t *testing.T) {
	var env Env
	p, err := expr.Compile(`store.Get("foo")`, expr.Env(env))
	assert.NoError(t, err)

	out, err := expr.Run(p, Env{Store: StoreImpt{}})
	assert.NoError(t, err)
	assert.Equal(t, 42, out)

	_, err = expr.Compile(`store.Set("foo", 100)`, expr.Env(env))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type interface_test.StoreInterface has no method Set")
}
