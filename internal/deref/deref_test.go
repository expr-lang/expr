package deref_test

import (
	"reflect"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"

	"github.com/expr-lang/expr/internal/deref"
)

func TestDeref(t *testing.T) {
	a := uint(42)
	b := &a
	c := &b
	d := &c

	got := deref.Deref(d)
	assert.Equal(t, uint(42), got)
}

func TestDeref_mix_ptr_with_interface(t *testing.T) {
	a := uint(42)
	var b any = &a
	var c any = &b
	d := &c

	got := deref.Deref(d)
	assert.Equal(t, uint(42), got)
}

func TestDeref_nil(t *testing.T) {
	var a *int
	assert.Nil(t, deref.Deref(a))
	assert.Nil(t, deref.Deref(nil))
}

func TestType(t *testing.T) {
	a := uint(42)
	b := &a
	c := &b
	d := &c

	dt := deref.Type(reflect.TypeOf(d))
	assert.Equal(t, reflect.Uint, dt.Kind())
}

func TestType_two_ptr_with_interface(t *testing.T) {
	a := uint(42)
	var b any = &a

	dt := deref.Type(reflect.TypeOf(b))
	assert.Equal(t, reflect.Uint, dt.Kind())

}

func TestType_three_ptr_with_interface(t *testing.T) {
	a := uint(42)
	var b any = &a
	var c any = &b

	dt := deref.Type(reflect.TypeOf(c))
	assert.Equal(t, reflect.Interface, dt.Kind())
}

func TestType_nil(t *testing.T) {
	assert.Nil(t, deref.Type(nil))
}

func TestValue(t *testing.T) {
	a := uint(42)
	b := &a
	c := &b
	d := &c

	got := deref.Value(reflect.ValueOf(d))
	assert.Equal(t, uint(42), got.Interface())
}

func TestValue_two_ptr_with_interface(t *testing.T) {
	a := uint(42)
	var b any = &a

	got := deref.Value(reflect.ValueOf(b))
	assert.Equal(t, uint(42), got.Interface())
}

func TestValue_three_ptr_with_interface(t *testing.T) {
	a := uint(42)
	var b any = &a
	c := &b

	got := deref.Value(reflect.ValueOf(c))
	assert.Equal(t, uint(42), got.Interface())
}

func TestValue_nil(t *testing.T) {
	got := deref.Value(reflect.ValueOf(nil))
	assert.False(t, got.IsValid())
}

func TestValue_nil_in_chain(t *testing.T) {
	var a any = nil
	var b any = &a
	c := &b

	got := deref.Value(reflect.ValueOf(c))
	assert.True(t, got.IsValid())
	assert.True(t, got.IsNil())
	assert.Nil(t, got.Interface())
}
