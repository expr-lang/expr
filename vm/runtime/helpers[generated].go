package runtime

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func Equal(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) == int(y)
		case uint8:
			return int(x) == int(y)
		case uint16:
			return int(x) == int(y)
		case uint32:
			return int(x) == int(y)
		case uint64:
			return int(x) == int(y)
		case int:
			return int(x) == int(y)
		case int8:
			return int(x) == int(y)
		case int16:
			return int(x) == int(y)
		case int32:
			return int(x) == int(y)
		case int64:
			return int(x) == int(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) == float64(y)
		case uint8:
			return float64(x) == float64(y)
		case uint16:
			return float64(x) == float64(y)
		case uint32:
			return float64(x) == float64(y)
		case uint64:
			return float64(x) == float64(y)
		case int:
			return float64(x) == float64(y)
		case int8:
			return float64(x) == float64(y)
		case int16:
			return float64(x) == float64(y)
		case int32:
			return float64(x) == float64(y)
		case int64:
			return float64(x) == float64(y)
		case float32:
			return float64(x) == float64(y)
		case float64:
			return float64(x) == float64(y)
		}
	case []any:
		switch y := b.(type) {
		case []string:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []uint64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []int64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []float32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []float64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		case []any:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if !Equal(x[i], y[i]) {
					return false
				}
			}
			return true
		}
	case []string:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []string:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint8:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint16:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []uint64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []uint64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int8:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int8:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int16:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int16:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []int64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []int64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []float32:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []float32:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case []float64:
		switch y := b.(type) {
		case []any:
			return Equal(y, x)
		case []float64:
			if len(x) != len(y) {
				return false
			}
			for i := range x {
				if x[i] != y[i] {
					return false
				}
			}
			return true
		}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x == y
		}
	case bool:
		switch y := b.(type) {
		case bool:
			return x == y
		}
	}
	if IsNil(a) && IsNil(b) {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func Less(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return false // nil is not less than nil
	}
	if a == nil {
		return true // nil is less than any non-nil value
	}
	if b == nil {
		return false // non-nil is not less than nil
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) < int(y)
		case uint8:
			return int(x) < int(y)
		case uint16:
			return int(x) < int(y)
		case uint32:
			return int(x) < int(y)
		case uint64:
			return int(x) < int(y)
		case int:
			return int(x) < int(y)
		case int8:
			return int(x) < int(y)
		case int16:
			return int(x) < int(y)
		case int32:
			return int(x) < int(y)
		case int64:
			return int(x) < int(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) < float64(y)
		case uint8:
			return float64(x) < float64(y)
		case uint16:
			return float64(x) < float64(y)
		case uint32:
			return float64(x) < float64(y)
		case uint64:
			return float64(x) < float64(y)
		case int:
			return float64(x) < float64(y)
		case int8:
			return float64(x) < float64(y)
		case int16:
			return float64(x) < float64(y)
		case int32:
			return float64(x) < float64(y)
		case int64:
			return float64(x) < float64(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) < float64(y)
		case uint8:
			return float64(x) < float64(y)
		case uint16:
			return float64(x) < float64(y)
		case uint32:
			return float64(x) < float64(y)
		case uint64:
			return float64(x) < float64(y)
		case int:
			return float64(x) < float64(y)
		case int8:
			return float64(x) < float64(y)
		case int16:
			return float64(x) < float64(y)
		case int32:
			return float64(x) < float64(y)
		case int64:
			return float64(x) < float64(y)
		case float32:
			return float64(x) < float64(y)
		case float64:
			return float64(x) < float64(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x < y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x < y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T < %T", a, b))
}

func More(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return false // nil is not more than nil
	}
	if a == nil {
		return false // nil is not more than any non-nil value
	}
	if b == nil {
		return true // non-nil is more than nil
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) > int(y)
		case uint8:
			return int(x) > int(y)
		case uint16:
			return int(x) > int(y)
		case uint32:
			return int(x) > int(y)
		case uint64:
			return int(x) > int(y)
		case int:
			return int(x) > int(y)
		case int8:
			return int(x) > int(y)
		case int16:
			return int(x) > int(y)
		case int32:
			return int(x) > int(y)
		case int64:
			return int(x) > int(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) > float64(y)
		case uint8:
			return float64(x) > float64(y)
		case uint16:
			return float64(x) > float64(y)
		case uint32:
			return float64(x) > float64(y)
		case uint64:
			return float64(x) > float64(y)
		case int:
			return float64(x) > float64(y)
		case int8:
			return float64(x) > float64(y)
		case int16:
			return float64(x) > float64(y)
		case int32:
			return float64(x) > float64(y)
		case int64:
			return float64(x) > float64(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) > float64(y)
		case uint8:
			return float64(x) > float64(y)
		case uint16:
			return float64(x) > float64(y)
		case uint32:
			return float64(x) > float64(y)
		case uint64:
			return float64(x) > float64(y)
		case int:
			return float64(x) > float64(y)
		case int8:
			return float64(x) > float64(y)
		case int16:
			return float64(x) > float64(y)
		case int32:
			return float64(x) > float64(y)
		case int64:
			return float64(x) > float64(y)
		case float32:
			return float64(x) > float64(y)
		case float64:
			return float64(x) > float64(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x > y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x > y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T > %T", a, b))
}

func LessOrEqual(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true // nil is equal to nil
	}
	if a == nil {
		return true // nil is less than or equal to any non-nil value
	}
	if b == nil {
		return false // non-nil is not less than or equal to nil
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) <= int(y)
		case uint8:
			return int(x) <= int(y)
		case uint16:
			return int(x) <= int(y)
		case uint32:
			return int(x) <= int(y)
		case uint64:
			return int(x) <= int(y)
		case int:
			return int(x) <= int(y)
		case int8:
			return int(x) <= int(y)
		case int16:
			return int(x) <= int(y)
		case int32:
			return int(x) <= int(y)
		case int64:
			return int(x) <= int(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) <= float64(y)
		case uint8:
			return float64(x) <= float64(y)
		case uint16:
			return float64(x) <= float64(y)
		case uint32:
			return float64(x) <= float64(y)
		case uint64:
			return float64(x) <= float64(y)
		case int:
			return float64(x) <= float64(y)
		case int8:
			return float64(x) <= float64(y)
		case int16:
			return float64(x) <= float64(y)
		case int32:
			return float64(x) <= float64(y)
		case int64:
			return float64(x) <= float64(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) <= float64(y)
		case uint8:
			return float64(x) <= float64(y)
		case uint16:
			return float64(x) <= float64(y)
		case uint32:
			return float64(x) <= float64(y)
		case uint64:
			return float64(x) <= float64(y)
		case int:
			return float64(x) <= float64(y)
		case int8:
			return float64(x) <= float64(y)
		case int16:
			return float64(x) <= float64(y)
		case int32:
			return float64(x) <= float64(y)
		case int64:
			return float64(x) <= float64(y)
		case float32:
			return float64(x) <= float64(y)
		case float64:
			return float64(x) <= float64(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x <= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x <= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T <= %T", a, b))
}

func MoreOrEqual(a, b interface{}) bool {
	// Handle nil values first
	if IsNil(a) && IsNil(b) {
		return true // nil is equal to nil
	}
	if a == nil {
		return false // nil is not more than or equal to any non-nil value
	}
	if b == nil {
		return true // non-nil is more than or equal to nil
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) >= int(y)
		case uint8:
			return int(x) >= int(y)
		case uint16:
			return int(x) >= int(y)
		case uint32:
			return int(x) >= int(y)
		case uint64:
			return int(x) >= int(y)
		case int:
			return int(x) >= int(y)
		case int8:
			return int(x) >= int(y)
		case int16:
			return int(x) >= int(y)
		case int32:
			return int(x) >= int(y)
		case int64:
			return int(x) >= int(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) >= float64(y)
		case uint8:
			return float64(x) >= float64(y)
		case uint16:
			return float64(x) >= float64(y)
		case uint32:
			return float64(x) >= float64(y)
		case uint64:
			return float64(x) >= float64(y)
		case int:
			return float64(x) >= float64(y)
		case int8:
			return float64(x) >= float64(y)
		case int16:
			return float64(x) >= float64(y)
		case int32:
			return float64(x) >= float64(y)
		case int64:
			return float64(x) >= float64(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) >= float64(y)
		case uint8:
			return float64(x) >= float64(y)
		case uint16:
			return float64(x) >= float64(y)
		case uint32:
			return float64(x) >= float64(y)
		case uint64:
			return float64(x) >= float64(y)
		case int:
			return float64(x) >= float64(y)
		case int8:
			return float64(x) >= float64(y)
		case int16:
			return float64(x) >= float64(y)
		case int32:
			return float64(x) >= float64(y)
		case int64:
			return float64(x) >= float64(y)
		case float32:
			return float64(x) >= float64(y)
		case float64:
			return float64(x) >= float64(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x >= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y) || x.Equal(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x >= y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T >= %T", a, b))
}

func Add(a, b interface{}) interface{} {
	// Handle nil values first
	if a == nil {
		switch y := b.(type) {
		case string:
			return "" + y
		case uint:
			return 0 + int(y)
		case uint8:
			return 0 + int(y)
		case uint16:
			return 0 + int(y)
		case uint32:
			return 0 + int(y)
		case uint64:
			return 0 + int(y)
		case int:
			return 0 + y
		case int8:
			return 0 + int(y)
		case int16:
			return 0 + int(y)
		case int32:
			return 0 + int(y)
		case int64:
			return 0 + int(y)
		case float32:
			return 0.0 + float64(y)
		case float64:
			return 0.0 + y
		default:
			return fmt.Sprint(nil) + fmt.Sprint(b)
		}
	}
	if b == nil {
		switch x := a.(type) {
		case string:
			return x + ""
		case uint:
			return int(x) + 0
		case uint8:
			return int(x) + 0
		case uint16:
			return int(x) + 0
		case uint32:
			return int(x) + 0
		case uint64:
			return int(x) + 0
		case int:
			return x + 0
		case int8:
			return int(x) + 0
		case int16:
			return int(x) + 0
		case int32:
			return int(x) + 0
		case int64:
			return int(x) + 0
		case float32:
			return float64(x) + 0.0
		case float64:
			return x + 0.0
		default:
			return fmt.Sprint(a) + fmt.Sprint(nil)
		}
	}

	// Handle string concatenation
	if str, ok := a.(string); ok {
		return str + fmt.Sprint(b)
	}
	if str, ok := b.(string); ok {
		return fmt.Sprint(a) + str
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) + int(y)
		case uint8:
			return int(x) + int(y)
		case uint16:
			return int(x) + int(y)
		case uint32:
			return int(x) + int(y)
		case uint64:
			return int(x) + int(y)
		case int:
			return int(x) + int(y)
		case int8:
			return int(x) + int(y)
		case int16:
			return int(x) + int(y)
		case int32:
			return int(x) + int(y)
		case int64:
			return int(x) + int(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) + float64(y)
		case uint8:
			return float64(x) + float64(y)
		case uint16:
			return float64(x) + float64(y)
		case uint32:
			return float64(x) + float64(y)
		case uint64:
			return float64(x) + float64(y)
		case int:
			return float64(x) + float64(y)
		case int8:
			return float64(x) + float64(y)
		case int16:
			return float64(x) + float64(y)
		case int32:
			return float64(x) + float64(y)
		case int64:
			return float64(x) + float64(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) + float64(y)
		case uint8:
			return float64(x) + float64(y)
		case uint16:
			return float64(x) + float64(y)
		case uint32:
			return float64(x) + float64(y)
		case uint64:
			return float64(x) + float64(y)
		case int:
			return float64(x) + float64(y)
		case int8:
			return float64(x) + float64(y)
		case int16:
			return float64(x) + float64(y)
		case int32:
			return float64(x) + float64(y)
		case int64:
			return float64(x) + float64(y)
		case float32:
			return float64(x) + float64(y)
		case float64:
			return float64(x) + float64(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x + y
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return x + fmt.Sprintf("%v", y)
		}
	case time.Time:
		switch y := b.(type) {
		case time.Duration:
			return x.Add(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Time:
			return y.Add(x)
		case time.Duration:
			return x + y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T + %T", a, b))
}

func Subtract(a, b interface{}) interface{} {
	// Handle nil values first
	if a == nil {
		switch y := b.(type) {
		case string:
			return 0 - len(y) // Convert string length to numeric subtraction
		case uint:
			return 0 - int(y)
		case uint8:
			return 0 - int(y)
		case uint16:
			return 0 - int(y)
		case uint32:
			return 0 - int(y)
		case uint64:
			return 0 - int(y)
		case int:
			return 0 - y
		case int8:
			return 0 - int(y)
		case int16:
			return 0 - int(y)
		case int32:
			return 0 - int(y)
		case int64:
			return 0 - int(y)
		case float32:
			return 0.0 - float64(y)
		case float64:
			return 0.0 - y
		default:
			return 0 // Default to 0 for unknown types
		}
	}
	if b == nil {
		switch x := a.(type) {
		case string:
			return len(x) - 0 // Convert string length to numeric subtraction
		case uint:
			return int(x) - 0
		case uint8:
			return int(x) - 0
		case uint16:
			return int(x) - 0
		case uint32:
			return int(x) - 0
		case uint64:
			return int(x) - 0
		case int:
			return x - 0
		case int8:
			return int(x) - 0
		case int16:
			return int(x) - 0
		case int32:
			return int(x) - 0
		case int64:
			return int(x) - 0
		case float32:
			return float64(x) - 0.0
		case float64:
			return x - 0.0
		default:
			return 0 // Default to 0 for unknown types
		}
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) - int(y)
		case uint8:
			return int(x) - int(y)
		case uint16:
			return int(x) - int(y)
		case uint32:
			return int(x) - int(y)
		case uint64:
			return int(x) - int(y)
		case int:
			return int(x) - int(y)
		case int8:
			return int(x) - int(y)
		case int16:
			return int(x) - int(y)
		case int32:
			return int(x) - int(y)
		case int64:
			return int(x) - int(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) - float64(y)
		case uint8:
			return float64(x) - float64(y)
		case uint16:
			return float64(x) - float64(y)
		case uint32:
			return float64(x) - float64(y)
		case uint64:
			return float64(x) - float64(y)
		case int:
			return float64(x) - float64(y)
		case int8:
			return float64(x) - float64(y)
		case int16:
			return float64(x) - float64(y)
		case int32:
			return float64(x) - float64(y)
		case int64:
			return float64(x) - float64(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) - float64(y)
		case uint8:
			return float64(x) - float64(y)
		case uint16:
			return float64(x) - float64(y)
		case uint32:
			return float64(x) - float64(y)
		case uint64:
			return float64(x) - float64(y)
		case int:
			return float64(x) - float64(y)
		case int8:
			return float64(x) - float64(y)
		case int16:
			return float64(x) - float64(y)
		case int32:
			return float64(x) - float64(y)
		case int64:
			return float64(x) - float64(y)
		case float32:
			return float64(x) - float64(y)
		case float64:
			return float64(x) - float64(y)
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Sub(y)
		case time.Duration:
			return x.Add(-y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Duration:
			return x - y
		}
	}
	panic(fmt.Sprintf("invalid operation: %T - %T", a, b))
}

func Multiply(a, b interface{}) interface{} {
	// Handle nil values first
	if a == nil || b == nil {
		return 0 // Any multiplication with nil results in 0
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint16:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			return int(x) * int(y)
		case uint8:
			return int(x) * int(y)
		case uint16:
			return int(x) * int(y)
		case uint32:
			return int(x) * int(y)
		case uint64:
			return int(x) * int(y)
		case int:
			return int(x) * int(y)
		case int8:
			return int(x) * int(y)
		case int16:
			return int(x) * int(y)
		case int32:
			return int(x) * int(y)
		case int64:
			return int(x) * int(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			return float64(x) * float64(y)
		case uint8:
			return float64(x) * float64(y)
		case uint16:
			return float64(x) * float64(y)
		case uint32:
			return float64(x) * float64(y)
		case uint64:
			return float64(x) * float64(y)
		case int:
			return float64(x) * float64(y)
		case int8:
			return float64(x) * float64(y)
		case int16:
			return float64(x) * float64(y)
		case int32:
			return float64(x) * float64(y)
		case int64:
			return float64(x) * float64(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return float64(x) * float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			return float64(x) * float64(y)
		case uint8:
			return float64(x) * float64(y)
		case uint16:
			return float64(x) * float64(y)
		case uint32:
			return float64(x) * float64(y)
		case uint64:
			return float64(x) * float64(y)
		case int:
			return float64(x) * float64(y)
		case int8:
			return float64(x) * float64(y)
		case int16:
			return float64(x) * float64(y)
		case int32:
			return float64(x) * float64(y)
		case int64:
			return float64(x) * float64(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return float64(x) * float64(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case uint:
			return time.Duration(x) * time.Duration(y)
		case uint8:
			return time.Duration(x) * time.Duration(y)
		case uint16:
			return time.Duration(x) * time.Duration(y)
		case uint32:
			return time.Duration(x) * time.Duration(y)
		case uint64:
			return time.Duration(x) * time.Duration(y)
		case int:
			return time.Duration(x) * time.Duration(y)
		case int8:
			return time.Duration(x) * time.Duration(y)
		case int16:
			return time.Duration(x) * time.Duration(y)
		case int32:
			return time.Duration(x) * time.Duration(y)
		case int64:
			return time.Duration(x) * time.Duration(y)
		case float32:
			return float64(x) * float64(y)
		case float64:
			return float64(x) * float64(y)
		case time.Duration:
			return time.Duration(x) * time.Duration(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T * %T", a, b))
}

func Divide(a, b interface{}) float64 {
	// Handle nil values first
	if a == nil {
		return 0.0 // 0 divided by anything is 0
	}
	if b == nil {
		return 0.0 // Division by nil is treated as division by 0, which results in 0
	}

	// Handle numeric operations
	switch x := a.(type) {
	case uint:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint8:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	// Continue similar pattern for all other numeric types...
	case uint16:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case uint64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int8:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int16:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case int64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case float32:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	case float64:
		switch y := b.(type) {
		case uint:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case uint64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float32:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		case float64:
			if float64(y) == 0.0 {
				panic("integer divide by zero")
			}
			return float64(x) / float64(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
}

func Modulo(a, b interface{}) interface{} {
	// Handle nil values first
	if a == nil {
		return 0 // 0 modulo anything is 0
	}
	if b == nil {
		return 0 // Modulo by nil is treated as modulo by 0, which results in 0
	}

	// Handle numeric operations
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int8:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int16:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int32:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case int64:
		switch y := b.(type) {
		case int:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int8:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int16:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int32:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case int64:
			if int(y) == 0 {
				panic("integer divide by zero")
			}
			return int(x) % int(y)
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case float32:
		switch y := b.(type) {
		case int:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int8:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int16:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}
	case float64:
		switch y := b.(type) {
		case int:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int8:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int16:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case int64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float32:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		case float64:
			if float64(y) == 0.0 {
				panic("float modulo by zero")
			}
			return math.Mod(float64(x), float64(y))
		}

	}
	panic(fmt.Sprintf("invalid operation: %T %% %T", a, b))
}
