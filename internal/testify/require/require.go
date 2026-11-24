package require

import (
	"github.com/expr-lang/expr/internal/testify/assert"
)

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//	assert.Contains(t, "Hello World", "World")
//	assert.Contains(t, ["Hello", "World"], "World")
//	assert.Contains(t, {"Hello": "World"}, "Hello")
func Contains(t TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Contains(t, s, contains, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Equal asserts that two objects are equal.
//
//	assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func Equal(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Equal(t, expected, actual, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//	actualObj, err := SomeFunction()
//	assert.EqualError(t, err,  expectedErrorString)
func EqualError(t TestingT, theError error, errString string, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.EqualError(t, theError, errString, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// EqualValues asserts that two objects are equal or convertible to the same types
// and equal.
//
//	assert.EqualValues(t, uint32(123), int32(123))
func EqualValues(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.EqualValues(t, expected, actual, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//	  actualObj, err := SomeFunction()
//	  if assert.Error(t, err) {
//		   assert.Equal(t, expectedError, err)
//	  }
func Error(t TestingT, err error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Error(t, err, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// False asserts that the specified value is false.
//
//	assert.False(t, myBool)
func False(t TestingT, value bool, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.False(t, value, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Greater asserts that the first element is greater than the second
//
//	assert.Greater(t, 2, 1)
//	assert.Greater(t, float64(2), float64(1))
//	assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Greater(t, e1, e2, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Implements asserts that an object is implemented by the specified interface.
//
//	assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func Implements(t TestingT, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Implements(t, interfaceObject, object, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// IsType asserts that the specified objects are of the same type.
func IsType(t TestingT, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.IsType(t, expectedType, object, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//	assert.Len(t, mySlice, 3)
func Len(t TestingT, object interface{}, length int, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Len(t, object, length, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Less asserts that the first element is less than the second
//
//	assert.Less(t, 1, 2)
//	assert.Less(t, float64(1), float64(2))
//	assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Less(t, e1, e2, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Nil asserts that the specified object is nil.
//
//	assert.Nil(t, err)
func Nil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Nil(t, object, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// NoError asserts that a function returned no error (i.e. `nil`).
//
//	  actualObj, err := SomeFunction()
//	  if assert.NoError(t, err) {
//		   assert.Equal(t, expectedObj, actualObj)
//	  }
func NoError(t TestingT, err error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.NoError(t, err, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// NotNil asserts that the specified object is not nil.
//
//	assert.NotNil(t, err)
func NotNil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.NotNil(t, object, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//	assert.Panics(t, func(){ GoCrazy() })
func Panics(t TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.Panics(t, f, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// True asserts that the specified value is true.
//
//	assert.True(t, myBool)
func True(t TestingT, value bool, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.True(t, value, msgAndArgs...) {
		return
	}
	t.FailNow()
}
