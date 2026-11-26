package assert

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// These tests use only the standard testing package: the code under test is
// our assertion library, so it must not be used to verify itself.

type myInt int
type myString string
type myBytes []byte
type myTime time.Time

func TestObjectsAreEqual(t *testing.T) {
	tests := []struct {
		expected, actual interface{}
		want             bool
	}{
		{nil, nil, true},
		{nil, 0, false},
		{0, nil, false},
		{1, 1, true},
		{1, int64(1), false},
		{"a", "a", true},
		{"a", "b", false},
		{[]int{1, 2}, []int{1, 2}, true},
		{[]int{1, 2}, []int{2, 1}, false},
		{map[string]int{"a": 1}, map[string]int{"a": 1}, true},
		{&struct{ A int }{1}, &struct{ A int }{1}, true},
		{[]byte("abc"), []byte("abc"), true},
		{[]byte("abc"), []byte("abd"), false},
		{[]byte{}, []byte{}, true},
		{[]byte(nil), []byte(nil), true},
		{[]byte(nil), []byte{}, false},
		{[]byte{}, []byte(nil), false},
		{[]byte("abc"), "abc", false},
		{"abc", []byte("abc"), false},
	}
	for i, tt := range tests {
		if got := ObjectsAreEqual(tt.expected, tt.actual); got != tt.want {
			t.Errorf("#%d ObjectsAreEqual(%#v, %#v) = %v, want %v", i, tt.expected, tt.actual, got, tt.want)
		}
	}
}

func TestObjectsAreEqualValues(t *testing.T) {
	tests := []struct {
		expected, actual interface{}
		want             bool
	}{
		{nil, nil, true},
		{nil, 0, false},
		{0, nil, false},
		{uint32(10), int32(10), true},
		{int32(10), uint32(10), true},
		{int64(10), uint8(10), true},
		{uint8(10), int64(10), true},
		{float64(10), 10, true},
		{10, float64(10), true},
		{10.5, 10, false},
		// Known limitation inherited from upstream: with equal sizes the actual
		// value is converted to the expected type, so 10.5 truncates to 10.
		{10, 10.5, true},
		// The smaller type is converted to the larger one, so overflow
		// must not create false positives.
		{int(270), uint8(14), false},
		{uint8(14), int(270), false},
		{int64(1 << 40), int32(0), false},
		{myInt(3), 3, true},
		{myString("a"), "a", true},
		{"a", myString("a"), true},
		{myString("a"), "b", false},
		{"1", 1, false},
		{[]int{1}, []int64{1}, false},
		{struct{ A int }{1}, struct{ A int }{1}, true},
	}
	for i, tt := range tests {
		if got := ObjectsAreEqualValues(tt.expected, tt.actual); got != tt.want {
			t.Errorf("#%d ObjectsAreEqualValues(%#v, %#v) = %v, want %v", i, tt.expected, tt.actual, got, tt.want)
		}
	}
}

func TestIsNumericType(t *testing.T) {
	numeric := []interface{}{
		int(0), int8(0), int16(0), int32(0), int64(0),
		uint(0), uint8(0), uint16(0), uint32(0), uint64(0), uintptr(0),
		float32(0), float64(0), complex64(0), complex128(0), myInt(0),
	}
	for _, v := range numeric {
		if !isNumericType(reflect.TypeOf(v)) {
			t.Errorf("isNumericType(%T) = false, want true", v)
		}
	}
	other := []interface{}{"", true, []int{}, struct{}{}, map[int]int{}, new(int)}
	for _, v := range other {
		if isNumericType(reflect.TypeOf(v)) {
			t.Errorf("isNumericType(%T) = true, want false", v)
		}
	}
}

func TestIsNil(t *testing.T) {
	var nilErr *strings.Builder
	var nilIface error
	nils := []interface{}{
		nil,
		(*int)(nil),
		[]int(nil),
		map[string]int(nil),
		(chan int)(nil),
		(func())(nil),
		nilErr,
		nilIface,
	}
	for _, v := range nils {
		if !isNil(v) {
			t.Errorf("isNil(%#v) = false, want true", v)
		}
	}
	notNils := []interface{}{0, "", false, struct{}{}, []int{}, map[string]int{}, new(int), make(chan int), func() {}, errors.New("x")}
	for _, v := range notNils {
		if isNil(v) {
			t.Errorf("isNil(%#v) = true, want false", v)
		}
	}
}

func TestGetLen(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	tests := []struct {
		value  interface{}
		length int
		ok     bool
	}{
		{[]int{1, 2, 3}, 3, true},
		{[2]int{}, 2, true},
		{"abcd", 4, true},
		{map[int]int{1: 1}, 1, true},
		{ch, 2, true},
		{[]int(nil), 0, true},
		{nil, 0, false},
		{5, 0, false},
		{struct{}{}, 0, false},
	}
	for i, tt := range tests {
		l, ok := getLen(tt.value)
		if l != tt.length || ok != tt.ok {
			t.Errorf("#%d getLen(%#v) = (%d, %v), want (%d, %v)", i, tt.value, l, ok, tt.length, tt.ok)
		}
	}
}

func TestContainsElement(t *testing.T) {
	tests := []struct {
		list, element interface{}
		ok, found     bool
	}{
		{"Hello World", "World", true, true},
		{"Hello World", "world", true, false},
		{myString("Hello"), "ell", true, true},
		{[]string{"a", "b"}, "b", true, true},
		{[]string{"a", "b"}, "c", true, false},
		{[2]int{1, 2}, 2, true, true},
		{[]int{1, 2}, int64(2), true, false},
		{[][]byte{[]byte("x")}, []byte("x"), true, true},
		{map[string]int{"a": 1}, "a", true, true},
		{map[string]int{"a": 1}, 1, true, false},
		{[]int(nil), 1, true, false},
		{nil, 1, false, false},
		{1, 1, false, false},
		{struct{}{}, 1, false, false},
	}
	for i, tt := range tests {
		ok, found := containsElement(tt.list, tt.element)
		if ok != tt.ok || found != tt.found {
			t.Errorf("#%d containsElement(%#v, %#v) = (%v, %v), want (%v, %v)", i, tt.list, tt.element, ok, found, tt.ok, tt.found)
		}
	}
}

func TestMessageFromMsgAndArgs(t *testing.T) {
	tests := []struct {
		args []interface{}
		want string
	}{
		{nil, ""},
		{[]interface{}{}, ""},
		{[]interface{}{"msg"}, "msg"},
		{[]interface{}{"100%"}, "100%"},
		{[]interface{}{"got %d and %s", 1, "x"}, "got 1 and x"},
		{[]interface{}{struct{ A int }{5}}, "{A:5}"},
		{[]interface{}{42}, "42"},
	}
	for i, tt := range tests {
		if got := messageFromMsgAndArgs(tt.args...); got != tt.want {
			t.Errorf("#%d messageFromMsgAndArgs(%#v) = %q, want %q", i, tt.args, got, tt.want)
		}
	}
}

func TestIsTest(t *testing.T) {
	tests := []struct {
		name, prefix string
		want         bool
	}{
		{"Test", "Test", true},
		{"TestFoo", "Test", true},
		{"Test_foo", "Test", true},
		{"Test1", "Test", true},
		{"Testfoo", "Test", false},
		{"Foo", "Test", false},
		{"BenchmarkFoo", "Benchmark", true},
		{"Benchmarkfoo", "Benchmark", false},
		{"ExampleFoo", "Example", true},
	}
	for _, tt := range tests {
		if got := isTest(tt.name, tt.prefix); got != tt.want {
			t.Errorf("isTest(%q, %q) = %v, want %v", tt.name, tt.prefix, got, tt.want)
		}
	}
}

func TestIndentMessageLines(t *testing.T) {
	tests := []struct {
		message string
		longest int
		want    string
	}{
		{"", 3, ""},
		{"one", 3, "one"},
		{"one\ntwo", 3, "one\n\t    \ttwo"},
		{"one\ntwo\nthree", 1, "one\n\t  \ttwo\n\t  \tthree"},
	}
	for _, tt := range tests {
		if got := indentMessageLines(tt.message, tt.longest); got != tt.want {
			t.Errorf("indentMessageLines(%q, %d) = %q, want %q", tt.message, tt.longest, got, tt.want)
		}
	}
}

func TestLabeledOutput(t *testing.T) {
	got := labeledOutput(
		labeledContent{"Error", "line1\nline2"},
		labeledContent{"Test", "TestX"},
		labeledContent{"Messages", "m"},
	)
	want := "\tError:   \tline1\n\t         \tline2\n" +
		"\tTest:    \tTestX\n" +
		"\tMessages:\tm\n"
	if got != want {
		t.Errorf("labeledOutput() =\n%q\nwant\n%q", got, want)
	}
	if got := labeledOutput(); got != "" {
		t.Errorf("labeledOutput() = %q, want empty", got)
	}
}

func TestCallerInfo(t *testing.T) {
	// Frames from files in the assert, require and mock directories are
	// skipped, and the walk stops at the test function. Everything on this
	// stack lives in the assert directory, so nothing is reported. Real
	// callers are covered by internal/testify/caller_test.go.
	var info []string
	func() { info = CallerInfo() }()
	if len(info) != 0 {
		t.Errorf("CallerInfo() = %q, want empty", info)
	}
}

func TestValidateEqualArgs(t *testing.T) {
	if err := validateEqualArgs(nil, nil); err != nil {
		t.Errorf("validateEqualArgs(nil, nil) = %v, want nil", err)
	}
	if err := validateEqualArgs(1, "a"); err != nil {
		t.Errorf("validateEqualArgs(1, \"a\") = %v, want nil", err)
	}
	if err := validateEqualArgs(func() {}, nil); err == nil {
		t.Error("validateEqualArgs(func, nil) = nil, want error")
	}
	if err := validateEqualArgs(nil, func() {}); err == nil {
		t.Error("validateEqualArgs(nil, func) = nil, want error")
	}
}

func TestIsFunction(t *testing.T) {
	if isFunction(nil) {
		t.Error("isFunction(nil) = true")
	}
	if isFunction(1) {
		t.Error("isFunction(1) = true")
	}
	if !isFunction(func() {}) {
		t.Error("isFunction(func) = false")
	}
	if !isFunction((func())(nil)) {
		t.Error("isFunction(nil func) = false")
	}
}

func TestFormatUnequalValues(t *testing.T) {
	tests := []struct {
		expected, actual interface{}
		e, a             string
	}{
		{1, 2, "1", "2"},
		{"a", "b", `"a"`, `"b"`},
		{int64(1), int32(1), "int64(1)", "int32(1)"},
		{1, "1", "int(1)", `string("1")`},
		{time.Second, 2 * time.Second, "1s", "2s"},
		{[]int{1}, []int{2}, "[]int{1}", "[]int{2}"},
	}
	for _, tt := range tests {
		e, a := formatUnequalValues(tt.expected, tt.actual)
		if e != tt.e || a != tt.a {
			t.Errorf("formatUnequalValues(%#v, %#v) = (%q, %q), want (%q, %q)", tt.expected, tt.actual, e, a, tt.e, tt.a)
		}
	}
}

func TestTruncatingFormat(t *testing.T) {
	if got := truncatingFormat("abc"); got != `"abc"` {
		t.Errorf("truncatingFormat(short) = %q", got)
	}
	long := strings.Repeat("x", bufio.MaxScanTokenSize)
	got := truncatingFormat(long)
	if !strings.HasSuffix(got, "<... truncated>") {
		t.Errorf("truncatingFormat(long) not truncated, len %d", len(got))
	}
	if len(got) >= bufio.MaxScanTokenSize {
		t.Errorf("truncatingFormat(long) too long: %d", len(got))
	}
}

func TestDiff(t *testing.T) {
	empty := []struct{ expected, actual interface{} }{
		{nil, nil},
		{nil, 1},
		{1, nil},
		{1, 2},
		{1.0, 2.0},
		{true, false},
		{1, "a"},
		{[]int{1}, []string{"1"}},
	}
	for _, tt := range empty {
		if got := diff(tt.expected, tt.actual); got != "" {
			t.Errorf("diff(%#v, %#v) = %q, want empty", tt.expected, tt.actual, got)
		}
	}

	type S struct {
		A int
		B string
	}
	withDiff := []struct {
		expected, actual interface{}
		contains         []string
	}{
		{"foo\nbar", "foo\nbaz", []string{"-bar", "+baz"}},
		{S{1, "x"}, S{1, "y"}, []string{`- B: (string) (len=1) "x"`, `+ B: (string) (len=1) "y"`}},
		{&S{1, "x"}, &S{2, "x"}, []string{"- A: (int) 1", "+ A: (int) 2"}},
		{[]int{1, 2}, []int{1, 3}, []string{"- (int) 2", "+ (int) 3"}},
		{[2]int{1, 2}, [2]int{1, 3}, []string{"- (int) 2", "+ (int) 3"}},
		{map[string]int{"a": 1}, map[string]int{"a": 2}, []string{`- (string) (len=1) "a": (int) 1`, `+ (string) (len=1) "a": (int) 2`}},
		{
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			[]string{"-(time.Time) 2020-01-01 00:00:00 +0000 UTC", "+(time.Time) 2021-01-01 00:00:00 +0000 UTC"},
		},
	}
	for _, tt := range withDiff {
		got := diff(tt.expected, tt.actual)
		if !strings.HasPrefix(got, "\n\nDiff:\n--- Expected\n+++ Actual\n") {
			t.Errorf("diff(%#v, %#v) has no diff header:\n%s", tt.expected, tt.actual, got)
		}
		for _, s := range tt.contains {
			if !strings.Contains(got, s) {
				t.Errorf("diff(%#v, %#v) does not contain %q:\n%s", tt.expected, tt.actual, s, got)
			}
		}
	}
}

func TestDidPanic(t *testing.T) {
	panicked, value, stack := didPanic(func() {})
	if panicked || value != nil || stack != "" {
		t.Errorf("didPanic(no panic) = (%v, %v, %q)", panicked, value, stack)
	}

	panicked, value, stack = didPanic(func() { panic("boom") })
	if !panicked || value != "boom" {
		t.Errorf("didPanic(panic) = (%v, %v)", panicked, value)
	}
	if !strings.Contains(stack, "TestDidPanic") {
		t.Errorf("didPanic stack does not mention the test:\n%s", stack)
	}

	err := errors.New("e")
	_, value, _ = didPanic(func() { panic(err) })
	if value != err {
		t.Errorf("didPanic(panic(err)) value = %v, want %v", value, err)
	}
}

func TestCompare(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Second)
	tests := []struct {
		less, greater interface{}
	}{
		{int(1), int(2)},
		{int(-2), int(-1)},
		{int8(1), int8(2)},
		{int16(1), int16(2)},
		{int32(1), int32(2)},
		{int64(1), int64(2)},
		{uint(1), uint(2)},
		{uint8(1), uint8(2)},
		{uint16(1), uint16(2)},
		{uint32(1), uint32(2)},
		{uint64(1), uint64(2)},
		{uintptr(1), uintptr(2)},
		{float32(1.5), float32(2.5)},
		{float64(1.5), float64(2.5)},
		{"a", "b"},
		{"", "a"},
		{t1, t2},
		{myTime(t1), myTime(t2)},
		{[]byte("a"), []byte("b")},
		{[]byte(""), []byte("a")},
		{myBytes("a"), myBytes("b")},
		{myInt(1), myInt(2)},
		{myString("a"), myString("b")},
	}
	for _, tt := range tests {
		kind := reflect.ValueOf(tt.less).Kind()
		if r, ok := compare(tt.less, tt.greater, kind); !ok || r != compareLess {
			t.Errorf("compare(%#v, %#v) = (%v, %v), want (less, true)", tt.less, tt.greater, r, ok)
		}
		if r, ok := compare(tt.greater, tt.less, kind); !ok || r != compareGreater {
			t.Errorf("compare(%#v, %#v) = (%v, %v), want (greater, true)", tt.greater, tt.less, r, ok)
		}
		if r, ok := compare(tt.less, tt.less, kind); !ok || r != compareEqual {
			t.Errorf("compare(%#v, %#v) = (%v, %v), want (equal, true)", tt.less, tt.less, r, ok)
		}
	}

	// Same instant in different locations is equal.
	if r, ok := compare(t1, t1.In(time.FixedZone("X", 3600)), reflect.Struct); !ok || r != compareEqual {
		t.Errorf("compare(time in different zones) = (%v, %v), want (equal, true)", r, ok)
	}

	notComparable := []interface{}{
		struct{ A int }{1},
		[]int{1},
		map[int]int{},
		true,
		new(int),
		complex(1, 1),
	}
	for _, v := range notComparable {
		if _, ok := compare(v, v, reflect.ValueOf(v).Kind()); ok {
			t.Errorf("compare(%#v) reported comparable", v)
		}
	}
}

func TestContainsValue(t *testing.T) {
	values := []compareResult{compareLess, compareEqual}
	if !containsValue(values, compareLess) || !containsValue(values, compareEqual) {
		t.Error("containsValue missed a present value")
	}
	if containsValue(values, compareGreater) {
		t.Error("containsValue found an absent value")
	}
	if containsValue(nil, compareEqual) {
		t.Error("containsValue(nil) = true")
	}
}

func TestFailMessageFormat(t *testing.T) {
	m := &mockT{name: "TestSomething"}
	if Fail(m, "it broke", "with %s", "context") {
		t.Fatal("Fail returned true")
	}
	got := m.message
	for _, want := range []string{
		"\n\tError Trace:\t\n",
		"\n\tError:      \tit broke\n",
		"\n\tTest:       \tTestSomething\n",
		"\n\tMessages:   \twith context\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Fail output does not contain %q:\n%s", want, got)
		}
	}
	if m.helperCalls == 0 {
		t.Error("Fail did not call Helper()")
	}

	// Without Name() and without a message, those labels are omitted.
	b := &bareT{}
	Fail(b, "x")
	if strings.Contains(b.message, "Test:") || strings.Contains(b.message, "Messages:") {
		t.Errorf("unexpected labels in output:\n%s", b.message)
	}
}

// mockT implements TestingT, tHelper and Name() and records failures.
type mockT struct {
	name        string
	failed      bool
	message     string
	helperCalls int
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.message = fmt.Sprintf(format, args...)
}

func (m *mockT) Helper() { m.helperCalls++ }

func (m *mockT) Name() string { return m.name }

// bareT implements only TestingT.
type bareT struct {
	message string
}

func (b *bareT) Errorf(format string, args ...interface{}) {
	b.message = fmt.Sprintf(format, args...)
}
