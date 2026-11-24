// Package testify_test checks the "Error Trace" of failure messages. It must
// live outside the assert and require directories: CallerInfo skips frames
// from files there, so tests inside those packages never appear in a trace.
package testify_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/expr-lang/expr/internal/testify/assert"
	"github.com/expr-lang/expr/internal/testify/require"
)

// Uses only the standard testing package to check the assertion library.

type mockT struct{ message string }

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.message = fmt.Sprintf(format, args...)
}

func (m *mockT) FailNow() {}

func line() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file, line+1)
}

func TestErrorTracePointsAtCaller(t *testing.T) {
	m := &mockT{}
	want := line()
	assert.Equal(m, 1, 2)
	if !strings.Contains(m.message, "Error Trace:\t"+want+"\n") {
		t.Errorf("assert trace does not point at %s:\n%s", want, m.message)
	}

	m = &mockT{}
	want = line()
	require.Equal(m, 1, 2)
	if !strings.Contains(m.message, "Error Trace:\t"+want+"\n") {
		t.Errorf("require trace does not point at %s:\n%s", want, m.message)
	}
	if strings.Contains(m.message, "/require/") || strings.Contains(m.message, "/assert/") {
		t.Errorf("trace includes assertion library frames:\n%s", m.message)
	}
}

func TestErrorTraceIncludesHelpers(t *testing.T) {
	m := &mockT{}
	var helperLine string
	check := func() {
		helperLine = line()
		assert.True(m, false)
	}
	callLine := line()
	check()

	for _, want := range []string{helperLine, callLine} {
		if !strings.Contains(m.message, want) {
			t.Errorf("trace does not contain %s:\n%s", want, m.message)
		}
	}
	// The walk stops at the test function, so testing internals are not listed.
	if strings.Contains(m.message, "testing.go") {
		t.Errorf("trace includes the testing package:\n%s", m.message)
	}
}

func TestErrorTraceInSubtest(t *testing.T) {
	t.Run("sub", func(t *testing.T) {
		m := &mockT{}
		want := line()
		assert.NoError(m, fmt.Errorf("x"))
		if !strings.Contains(m.message, want) {
			t.Errorf("trace does not point at %s:\n%s", want, m.message)
		}
		if strings.Contains(m.message, "testing.go") {
			t.Errorf("trace includes the testing package:\n%s", m.message)
		}
	})
}
