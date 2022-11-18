package runtime

import (
	"time"
)

// tryToInt tries converting the given value to int without loosing precision.
func tryToInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case time.Duration:
		return int(v), true
	default:
		return 0, false
	}
}

// tryBothToInt tries converting both values to int, using tryToInt.
func tryBothToInt(a, b interface{}) (int, int, bool) {
	aInt, ok := tryToInt(a)
	if !ok {
		return 0, 0, false
	}
	bInt, ok := tryToInt(b)
	if !ok {
		return 0, 0, false
	}
	return aInt, bInt, true
}

// tryToFloat tries converting numeric types to float64.
func tryToFloat(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case time.Duration:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// tryBothToFloat tries converting both values to float64, using tryToFloat.
func tryBothToFloat(a, b interface{}) (float64, float64, bool) {
	aFloat, ok := tryToFloat(a)
	if !ok {
		return 0, 0, false
	}
	bFloat, ok := tryToFloat(b)
	if !ok {
		return 0, 0, false
	}
	return aFloat, bFloat, true
}

// tryToString converts the given value to string if being of string type.
func tryToString(v interface{}) (string, bool) {
	switch v := v.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

// tryBothToString tries converting both values to string, using tryToString.
func tryBothToString(a, b interface{}) (string, string, bool) {
	aStr, ok := tryToString(a)
	if !ok {
		return "", "", false
	}
	bStr, ok := tryToString(b)
	if !ok {
		return "", "", false
	}
	return aStr, bStr, true
}

// tryToTime converts the given value to time.Time if it is of time.Time type.
func tryToTime(v interface{}) (time.Time, bool) {
	switch v := v.(type) {
	case time.Time:
		return v, true
	default:
		return time.Time{}, false
	}
}

// tryBothToTime tries converting both values to time.Time, using tryToTime.
func tryBothToTime(a, b interface{}) (time.Time, time.Time, bool) {
	aTime, ok := tryToTime(a)
	if !ok {
		return time.Time{}, time.Time{}, false
	}
	bTime, ok := tryToTime(b)
	if !ok {
		return time.Time{}, time.Time{}, false
	}
	return aTime, bTime, true
}

// tryToDuration converts the given value to time.Duration if being of
// time.Duration type.
func tryToDuration(v interface{}) (time.Duration, bool) {
	switch v := v.(type) {
	case time.Duration:
		return v, true
	default:
		return 0, false
	}
}

// tryBothToDuration tries converting both values to time.Duration, using
// tryToDuration.
func tryBothToDuration(a, b interface{}) (time.Duration, time.Duration, bool) {
	aDur, ok := tryToDuration(a)
	if !ok {
		return 0, 0, false
	}
	bDur, ok := tryToDuration(b)
	if !ok {
		return 0, 0, false
	}
	return aDur, bDur, true
}

// typeName is a type representation, avoiding plain reflect-types.
type typeName int

// Supported type names.
const (
	typeUInt typeName = iota
	typeUInt8
	typeUInt16
	typeUInt32
	typeUInt64
	typeInt
	typeInt8
	typeInt16
	typeInt32
	typeInt64
	typeFloat32
	typeFloat64
	typeDuration
	typeString
	typeTime
)

// typeOf returns the typeName of the given value.
func typeOf(v interface{}) (typeName, bool) {
	switch v.(type) {
	case uint:
		return typeUInt, true
	case uint8:
		return typeUInt8, true
	case uint16:
		return typeUInt16, true
	case uint32:
		return typeUInt32, true
	case uint64:
		return typeUInt64, true
	case int:
		return typeInt, true
	case int8:
		return typeInt8, true
	case int16:
		return typeInt16, true
	case int32:
		return typeInt32, true
	case int64:
		return typeInt64, true
	case time.Duration:
		return typeDuration, true
	case float32:
		return typeFloat32, true
	case float64:
		return typeFloat64, true
	case time.Time:
		return typeTime, true
	case string:
		return typeString, true
	default:
		return 0, false
	}
}

// reorderByType takes two values and a typeName-order. It sorts both values
// according to the given order or returns false and the original one, if no
// matches where found. If only a single value was not found in the given list,
// it will be treated with less priority than the found one and therefore
// returned as second value. The third return value indicates, if the values have
// been swapped. The fourth one, if reordering was successful.
func reorderByType(a, b interface{}, order []typeName) (interface{}, interface{}, bool, bool) {
	aType, ok := typeOf(a)
	if !ok {
		return a, b, false, false
	}
	bType, ok := typeOf(b)
	if !ok {
		return a, b, false, false
	}
	// Search for a and b, remember if found at any time and if a nd b should be
	// swapped.
	search := true
	swap := false
	aFound := false
	bFound := false
	for _, name := range order {
		if aType == name {
			aFound = true
			if search {
				search = false
				swap = false
			}
		}
		if bType == name {
			bFound = true
			if search {
				search = false
				swap = true
			}
		}
	}
	if search || !aFound || !bFound {
		return a, b, false, false
	}
	if swap {
		a, b = b, a
	}
	return a, b, swap, true
}
