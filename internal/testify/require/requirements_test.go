package require

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

// These tests use only the standard testing package: the code under test is
// our assertion library, so it must not be used to verify itself.

type mockT struct {
	errorf      int
	failNow     int
	helperCalls int
	message     string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.errorf++
	m.message = fmt.Sprintf(format, args...)
}

func (m *mockT) FailNow() { m.failNow++ }

func (m *mockT) Helper() { m.helperCalls++ }

type shape interface{ Area() float64 }
type square struct{}

func (square) Area() float64 { return 0 }

func TestRequire(t *testing.T) {
	err := errors.New("some error")
	tests := []struct {
		name string
		pass func(t TestingT)
		fail func(t TestingT)
		msg  string
	}{
		{"Contains", func(t TestingT) { Contains(t, "abc", "b") }, func(t TestingT) { Contains(t, "abc", "x") }, `"abc" does not contain "x"`},
		{"Equal", func(t TestingT) { Equal(t, 1, 1) }, func(t TestingT) { Equal(t, 1, 2) }, "Not equal"},
		{"EqualError", func(t TestingT) { EqualError(t, io.EOF, "EOF") }, func(t TestingT) { EqualError(t, io.EOF, "x") }, "Error message not equal"},
		{"EqualValues", func(t TestingT) { EqualValues(t, int32(1), int64(1)) }, func(t TestingT) { EqualValues(t, 1, 2) }, "Not equal"},
		{"Error", func(t TestingT) { Error(t, err) }, func(t TestingT) { Error(t, nil) }, "An error is expected but got nil."},
		{"False", func(t TestingT) { False(t, false) }, func(t TestingT) { False(t, true) }, "Should be false"},
		{"Greater", func(t TestingT) { Greater(t, 2, 1) }, func(t TestingT) { Greater(t, 1, 2) }, `"1" is not greater than "2"`},
		{"Implements", func(t TestingT) { Implements(t, (*shape)(nil), square{}) }, func(t TestingT) { Implements(t, (*shape)(nil), 1) }, "int must implement require.shape"},
		{"IsType", func(t TestingT) { IsType(t, 1, 2) }, func(t TestingT) { IsType(t, 1, "a") }, "Object expected to be of type int, but was string"},
		{"Len", func(t TestingT) { Len(t, []int{1}, 1) }, func(t TestingT) { Len(t, []int{1}, 2) }, "should have 2 item(s), but has 1"},
		{"Less", func(t TestingT) { Less(t, 1, 2) }, func(t TestingT) { Less(t, 2, 1) }, `"2" is not less than "1"`},
		{"Nil", func(t TestingT) { Nil(t, nil) }, func(t TestingT) { Nil(t, 1) }, "Expected nil, but got: 1"},
		{"NoError", func(t TestingT) { NoError(t, nil) }, func(t TestingT) { NoError(t, err) }, "Received unexpected error:"},
		{"NotNil", func(t TestingT) { NotNil(t, 1) }, func(t TestingT) { NotNil(t, nil) }, "Expected value not to be nil."},
		{"Panics", func(t TestingT) { Panics(t, func() { panic("x") }) }, func(t TestingT) { Panics(t, func() {}) }, "should panic"},
		{"True", func(t TestingT) { True(t, true) }, func(t TestingT) { True(t, false) }, "Should be true"},
		{"message", func(t TestingT) { True(t, true, "want %d", 1) }, func(t TestingT) { True(t, false, "want %d", 1) }, "want 1"},
	}
	for _, tt := range tests {
		m := &mockT{}
		tt.pass(m)
		if m.errorf != 0 || m.failNow != 0 {
			t.Errorf("%s: passing case reported failure (Errorf=%d, FailNow=%d): %s", tt.name, m.errorf, m.failNow, m.message)
		}

		m = &mockT{}
		tt.fail(m)
		if m.errorf != 1 {
			t.Errorf("%s: failing case called Errorf %d times, want 1", tt.name, m.errorf)
		}
		if m.failNow != 1 {
			t.Errorf("%s: failing case called FailNow %d times, want 1", tt.name, m.failNow)
		}
		if m.helperCalls == 0 {
			t.Errorf("%s: failing case did not call Helper()", tt.name)
		}
		if !strings.Contains(m.message, tt.msg) {
			t.Errorf("%s: failure output does not contain %q:\n%s", tt.name, tt.msg, m.message)
		}
	}
}
