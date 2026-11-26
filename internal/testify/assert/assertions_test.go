package assert

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"time"
)

type assertCase struct {
	name string
	call func(t TestingT) bool
	ok   bool
	// msg must appear in the failure output; ignored when ok is true.
	msg string
}

// continuation matches the indentation labeledOutput adds to multi-line content.
var continuation = regexp.MustCompile(`\n\t +\t`)

func runAssertCases(t *testing.T, cases []assertCase) {
	t.Helper()
	for _, c := range cases {
		m := &mockT{}
		got := c.call(m)
		if got != c.ok {
			t.Errorf("%s: returned %v, want %v", c.name, got, c.ok)
		}
		if m.failed == c.ok {
			t.Errorf("%s: failed = %v, want %v\n%s", c.name, m.failed, !c.ok, m.message)
		}
		if !c.ok && !strings.Contains(continuation.ReplaceAllString(m.message, "\n"), c.msg) {
			t.Errorf("%s: failure output does not contain %q:\n%s", c.name, c.msg, m.message)
		}
		if !c.ok && m.helperCalls == 0 {
			t.Errorf("%s: Helper() was not called", c.name)
		}
	}
}

type shape interface{ Area() float64 }
type square struct{ side float64 }

func (s square) Area() float64 { return s.side * s.side }

func TestEqual(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"ints", func(t TestingT) bool { return Equal(t, 1, 1) }, true, ""},
		{"strings", func(t TestingT) bool { return Equal(t, "a", "a") }, true, ""},
		{"nils", func(t TestingT) bool { return Equal(t, nil, nil) }, true, ""},
		{"slices", func(t TestingT) bool { return Equal(t, []int{1}, []int{1}) }, true, ""},
		{"pointers by value", func(t TestingT) bool { return Equal(t, &square{1}, &square{1}) }, true, ""},
		{"bytes", func(t TestingT) bool { return Equal(t, []byte("a"), []byte("a")) }, true, ""},
		{"different ints", func(t TestingT) bool { return Equal(t, 1, 2) }, false, "Not equal: \nexpected: 1\nactual  : 2"},
		{"different types", func(t TestingT) bool { return Equal(t, int32(1), int64(1)) }, false, "expected: int32(1)\nactual  : int64(1)"},
		{"different strings", func(t TestingT) bool { return Equal(t, "a\nb", "a\nc") }, false, "Diff:"},
		{"nil vs value", func(t TestingT) bool { return Equal(t, nil, 0) }, false, "Not equal"},
		{"func", func(t TestingT) bool { return Equal(t, func() {}, func() {}) }, false, "cannot take func type as argument"},
		{"message", func(t TestingT) bool { return Equal(t, 1, 2, "case %d", 7) }, false, "Messages:   \tcase 7"},
	})
}

func TestEqualValues(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"same", func(t TestingT) bool { return EqualValues(t, 1, 1) }, true, ""},
		{"converted", func(t TestingT) bool { return EqualValues(t, uint32(123), int32(123)) }, true, ""},
		{"named type", func(t TestingT) bool { return EqualValues(t, myString("a"), "a") }, true, ""},
		{"different", func(t TestingT) bool { return EqualValues(t, 1, 2) }, false, "Not equal"},
		{"overflow", func(t TestingT) bool { return EqualValues(t, 270, uint8(14)) }, false, "expected: int(270)\nactual  : uint8(0xe)"},
		{"not convertible", func(t TestingT) bool { return EqualValues(t, "1", 1) }, false, "Not equal"},
	})
}

func TestNilAndNotNil(t *testing.T) {
	var nilPtr *int
	runAssertCases(t, []assertCase{
		{"Nil(nil)", func(t TestingT) bool { return Nil(t, nil) }, true, ""},
		{"Nil(typed nil)", func(t TestingT) bool { return Nil(t, nilPtr) }, true, ""},
		{"Nil(nil slice)", func(t TestingT) bool { return Nil(t, []int(nil)) }, true, ""},
		{"Nil(value)", func(t TestingT) bool { return Nil(t, 0) }, false, "Expected nil, but got: 0"},
		{"Nil(empty slice)", func(t TestingT) bool { return Nil(t, []int{}) }, false, "Expected nil, but got: []int{}"},
		{"NotNil(value)", func(t TestingT) bool { return NotNil(t, 0) }, true, ""},
		{"NotNil(empty slice)", func(t TestingT) bool { return NotNil(t, []int{}) }, true, ""},
		{"NotNil(nil)", func(t TestingT) bool { return NotNil(t, nil) }, false, "Expected value not to be nil."},
		{"NotNil(typed nil)", func(t TestingT) bool { return NotNil(t, nilPtr) }, false, "Expected value not to be nil."},
	})
}

func TestLen(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"slice", func(t TestingT) bool { return Len(t, []int{1, 2}, 2) }, true, ""},
		{"map", func(t TestingT) bool { return Len(t, map[int]int{1: 1}, 1) }, true, ""},
		{"string", func(t TestingT) bool { return Len(t, "abc", 3) }, true, ""},
		{"nil slice", func(t TestingT) bool { return Len(t, []int(nil), 0) }, true, ""},
		{"wrong length", func(t TestingT) bool { return Len(t, []int{1, 2}, 3) }, false, `"[1 2]" should have 3 item(s), but has 2`},
		{"no len", func(t TestingT) bool { return Len(t, 5, 1) }, false, `"5" could not be applied builtin len()`},
		{"nil", func(t TestingT) bool { return Len(t, nil, 0) }, false, "could not be applied builtin len()"},
	})
}

func TestTrueAndFalse(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"True(true)", func(t TestingT) bool { return True(t, true) }, true, ""},
		{"True(false)", func(t TestingT) bool { return True(t, false) }, false, "Should be true"},
		{"True(false, msg)", func(t TestingT) bool { return True(t, false, "want %v, got %v", 1, 2) }, false, "want 1, got 2"},
		{"False(false)", func(t TestingT) bool { return False(t, false) }, true, ""},
		{"False(true)", func(t TestingT) bool { return False(t, true) }, false, "Should be false"},
	})
}

func TestContains(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"substring", func(t TestingT) bool { return Contains(t, "Hello World", "World") }, true, ""},
		{"slice", func(t TestingT) bool { return Contains(t, []string{"a", "b"}, "b") }, true, ""},
		{"map key", func(t TestingT) bool { return Contains(t, map[string]int{"a": 1}, "a") }, true, ""},
		{"missing substring", func(t TestingT) bool { return Contains(t, "Hello", "World") }, false, `"Hello" does not contain "World"`},
		{"missing element", func(t TestingT) bool { return Contains(t, []int{1}, 2) }, false, "[]int{1} does not contain 2"},
		{"map value", func(t TestingT) bool { return Contains(t, map[string]int{"a": 1}, 1) }, false, "does not contain 1"},
		{"not a list", func(t TestingT) bool { return Contains(t, 1, 1) }, false, "1 could not be applied builtin len()"},
		{"nil", func(t TestingT) bool { return Contains(t, nil, 1) }, false, "<nil> could not be applied builtin len()"},
	})
}

func TestPanics(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"panics", func(t TestingT) bool { return Panics(t, func() { panic("x") }) }, true, ""},
		{"no panic", func(t TestingT) bool { return Panics(t, func() {}) }, false, "should panic\n\tPanic value:\t<nil>"},
	})
}

func TestPanicsWithValue(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"value", func(t TestingT) bool { return PanicsWithValue(t, "x", func() { panic("x") }) }, true, ""},
		{"int value", func(t TestingT) bool { return PanicsWithValue(t, 1, func() { panic(1) }) }, true, ""},
		{"wrong value", func(t TestingT) bool { return PanicsWithValue(t, "x", func() { panic("y") }) }, false, "should panic with value:\t\"x\"\n\tPanic value:\t\"y\"\n\tPanic stack:\t"},
		{"wrong type", func(t TestingT) bool { return PanicsWithValue(t, 1, func() { panic(int64(1)) }) }, false, "Panic value:\t1"},
		{"no panic", func(t TestingT) bool { return PanicsWithValue(t, "x", func() {}) }, false, "should panic\n\tPanic value:\t<nil>"},
	})
}

func TestPanicsWithError(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"error", func(t TestingT) bool { return PanicsWithError(t, "boom", func() { panic(errors.New("boom")) }) }, true, ""},
		{"wrapped", func(t TestingT) bool {
			return PanicsWithError(t, "ctx: boom", func() { panic(fmt.Errorf("ctx: %w", errors.New("boom"))) })
		}, true, ""},
		{"wrong message", func(t TestingT) bool { return PanicsWithError(t, "boom", func() { panic(errors.New("bang")) }) }, false, "should panic with error message:\t\"boom\""},
		{"not an error", func(t TestingT) bool { return PanicsWithError(t, "boom", func() { panic("boom") }) }, false, "Panic value:\t\"boom\""},
		{"no panic", func(t TestingT) bool { return PanicsWithError(t, "boom", func() {}) }, false, "should panic\n\tPanic value:\t<nil>"},
	})
}

func TestNoErrorAndError(t *testing.T) {
	err := errors.New("some error")
	runAssertCases(t, []assertCase{
		{"NoError(nil)", func(t TestingT) bool { return NoError(t, nil) }, true, ""},
		{"NoError(err)", func(t TestingT) bool { return NoError(t, err) }, false, "Received unexpected error:\nsome error"},
		{"Error(err)", func(t TestingT) bool { return Error(t, err) }, true, ""},
		{"Error(nil)", func(t TestingT) bool { return Error(t, nil) }, false, "An error is expected but got nil."},
	})
}

func TestEqualError(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"match", func(t TestingT) bool { return EqualError(t, io.EOF, "EOF") }, true, ""},
		{"mismatch", func(t TestingT) bool { return EqualError(t, io.EOF, "eof") }, false, "Error message not equal:\nexpected: \"eof\"\nactual  : \"EOF\""},
		{"nil", func(t TestingT) bool { return EqualError(t, nil, "EOF") }, false, "An error is expected but got nil."},
	})
}

func TestImplements(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"implements", func(t TestingT) bool { return Implements(t, (*shape)(nil), square{}) }, true, ""},
		{"pointer implements", func(t TestingT) bool { return Implements(t, (*shape)(nil), &square{}) }, true, ""},
		{"does not", func(t TestingT) bool { return Implements(t, (*shape)(nil), 1) }, false, "int must implement assert.shape"},
		{"nil", func(t TestingT) bool { return Implements(t, (*shape)(nil), nil) }, false, "Cannot check if nil implements assert.shape"},
	})
}

func TestIsType(t *testing.T) {
	runAssertCases(t, []assertCase{
		{"same", func(t TestingT) bool { return IsType(t, square{}, square{side: 2}) }, true, ""},
		{"pointer", func(t TestingT) bool { return IsType(t, &square{}, &square{}) }, true, ""},
		{"nil and nil", func(t TestingT) bool { return IsType(t, nil, nil) }, true, ""},
		{"different", func(t TestingT) bool { return IsType(t, int64(0), 0) }, false, "Object expected to be of type int64, but was int"},
		{"pointer vs value", func(t TestingT) bool { return IsType(t, &square{}, square{}) }, false, "Object expected to be of type *assert.square, but was assert.square"},
	})
}

func TestGreaterAndLess(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	runAssertCases(t, []assertCase{
		{"Greater ints", func(t TestingT) bool { return Greater(t, 2, 1) }, true, ""},
		{"Greater floats", func(t TestingT) bool { return Greater(t, 2.5, 2.4) }, true, ""},
		{"Greater strings", func(t TestingT) bool { return Greater(t, "b", "a") }, true, ""},
		{"Greater times", func(t TestingT) bool { return Greater(t, t1.Add(1), t1) }, true, ""},
		{"Greater equal", func(t TestingT) bool { return Greater(t, 1, 1) }, false, `"1" is not greater than "1"`},
		{"Greater less", func(t TestingT) bool { return Greater(t, 1, 2) }, false, `"1" is not greater than "2"`},
		{"Less ints", func(t TestingT) bool { return Less(t, 1, 2) }, true, ""},
		{"Less bytes", func(t TestingT) bool { return Less(t, []byte("a"), []byte("b")) }, true, ""},
		{"Less equal", func(t TestingT) bool { return Less(t, "a", "a") }, false, `"a" is not less than "a"`},
		{"Less greater", func(t TestingT) bool { return Less(t, 2, 1) }, false, `"2" is not less than "1"`},
		{"different kinds", func(t TestingT) bool { return Less(t, 1, int64(2)) }, false, "Elements should be the same type"},
		{"not comparable", func(t TestingT) bool { return Less(t, struct{}{}, struct{}{}) }, false, `Can not compare type "struct {}"`},
		{"bool", func(t TestingT) bool { return Greater(t, true, false) }, false, `Can not compare type "bool"`},
	})
}

func TestNoHelperOnSuccess(t *testing.T) {
	// Assertions that pass must not report a failure, with or without Helper().
	b := &bareT{}
	if !Equal(b, 1, 1) || !True(b, true) || !NoError(b, nil) || !Contains(b, "ab", "a") {
		t.Fatal("passing assertion returned false")
	}
	if b.message != "" {
		t.Errorf("passing assertion reported: %s", b.message)
	}
}
