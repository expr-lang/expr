package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func Test_tryToInt(t *testing.T) {
	tests := []struct {
		v       interface{}
		wantInt int
		wantOK  bool
	}{
		{
			v:       0,
			wantInt: 0,
			wantOK:  true,
		},
		{
			v:       uint(368),
			wantInt: 368,
			wantOK:  true,
		},
		{
			v:       uint8(25),
			wantInt: 25,
			wantOK:  true,
		},
		{
			v:       uint16(487),
			wantInt: 487,
			wantOK:  true,
		},
		{
			v:       uint32(489),
			wantInt: 489,
			wantOK:  true,
		}, {
			v:       uint64(445),
			wantInt: 445,
			wantOK:  true,
		}, {
			v:       int(-934),
			wantInt: -934,
			wantOK:  true,
		},
		{
			v:       int8(-3),
			wantInt: -3,
			wantOK:  true,
		},
		{
			v:       int16(24),
			wantInt: 24,
			wantOK:  true,
		},
		{
			v:       int32(-849),
			wantInt: -849,
			wantOK:  true,
		},
		{
			v:       int64(479),
			wantInt: 479,
			wantOK:  true,
		},
		{
			v:       12 * time.Hour,
			wantInt: int(12 * time.Hour),
			wantOK:  true,
		},
		{
			v:      "Hello World!",
			wantOK: false,
		},
		{
			v:      1.3,
			wantOK: false,
		},
		{
			v:      time.Now(),
			wantOK: false,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.v), tt.v), func(t *testing.T) {
			t.Parallel()
			got, ok := tryToInt(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantInt, got, "should return correct value")
			}
		})
	}
}

func Test_tryBothToInt(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		wantIntA int
		wantIntB int
		wantOK   bool
	}{
		{
			name:     "both ints",
			a:        194,
			b:        150,
			wantIntA: 194,
			wantIntB: 150,
			wantOK:   true,
		},
		{
			name:     "both convertable",
			a:        uint8(23),
			b:        int32(-344),
			wantIntA: 23,
			wantIntB: -344,
			wantOK:   true,
		},
		{
			name:   "both not convertable",
			a:      "Hello World!",
			b:      45.3,
			wantOK: false,
		},
		{
			name:   "one not convertable",
			a:      23,
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotA, gotB, ok := tryBothToInt(tt.a, tt.b)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantIntA, gotA, "should return correct value for a")
				assert.Equal(t, tt.wantIntB, gotB, "should return correct value for b")
			}
		})
	}
}

func Test_tryToFloat(t *testing.T) {
	tests := []struct {
		v           interface{}
		wantFloat64 float64
		wantOK      bool
	}{
		{
			v:           0,
			wantFloat64: 0,
			wantOK:      true,
		},
		{
			v:           uint(368),
			wantFloat64: 368,
			wantOK:      true,
		},
		{
			v:           uint8(25),
			wantFloat64: 25,
			wantOK:      true,
		},
		{
			v:           uint16(487),
			wantFloat64: 487,
			wantOK:      true,
		},
		{
			v:           uint32(489),
			wantFloat64: 489,
			wantOK:      true,
		}, {
			v:           uint64(445),
			wantFloat64: 445,
			wantOK:      true,
		}, {
			v:           int(-934),
			wantFloat64: -934,
			wantOK:      true,
		},
		{
			v:           int8(-3),
			wantFloat64: -3,
			wantOK:      true,
		},
		{
			v:           int16(24),
			wantFloat64: 24,
			wantOK:      true,
		},
		{
			v:           int32(-849),
			wantFloat64: -849,
			wantOK:      true,
		},
		{
			v:           int64(479),
			wantFloat64: 479,
			wantOK:      true,
		},
		{
			v:           12 * time.Hour,
			wantFloat64: float64(12 * time.Hour),
			wantOK:      true,
		},
		{
			v:      "Hello World!",
			wantOK: false,
		},
		{
			v:           float32(1.5),
			wantFloat64: 1.5,
			wantOK:      true,
		},
		{
			v:           -444.2,
			wantFloat64: -444.2,
			wantOK:      true,
		},
		{
			v:      time.Now(),
			wantOK: false,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.v), tt.v), func(t *testing.T) {
			t.Parallel()
			got, ok := tryToFloat(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantFloat64, got, "should return correct value")
			}
		})
	}
}

func Test_tryBothToFloat(t *testing.T) {
	tests := []struct {
		name         string
		a            interface{}
		b            interface{}
		wantFloat64A float64
		wantFloat64B float64
		wantOK       bool
	}{
		{
			name:         "both floats",
			a:            194.23,
			b:            150.8,
			wantFloat64A: 194.23,
			wantFloat64B: 150.8,
			wantOK:       true,
		},
		{
			name:         "both convertable",
			a:            uint8(23),
			b:            float32(-344),
			wantFloat64A: 23,
			wantFloat64B: -344,
			wantOK:       true,
		},
		{
			name:   "both not convertable",
			a:      "Hello World!",
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "one not convertable",
			a:      23.1,
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotA, gotB, ok := tryBothToFloat(tt.a, tt.b)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantFloat64A, gotA, "should return correct value for a")
				assert.Equal(t, tt.wantFloat64B, gotB, "should return correct value for b")
			}
		})
	}
}

func Test_tryToString(t *testing.T) {
	tests := []struct {
		v       interface{}
		wantStr string
		wantOK  bool
	}{
		{
			v:      uint(368),
			wantOK: false,
		},
		{
			v:      int(-934),
			wantOK: false,
		},
		{
			v:      12 * time.Hour,
			wantOK: false,
		},
		{
			v:       "Hello World!",
			wantStr: "Hello World!",
			wantOK:  true,
		},
		{
			v:       "",
			wantStr: "",
			wantOK:  true,
		},
		{
			v:      time.Now(),
			wantOK: false,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.v), tt.v), func(t *testing.T) {
			t.Parallel()
			got, ok := tryToString(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantStr, got, "should return correct value")
			}
		})
	}
}

func Test_tryBothToString(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		wantStrA string
		wantStrB string
		wantOK   bool
	}{
		{
			name:     "both string",
			a:        "husband",
			b:        "form",
			wantStrA: "husband",
			wantStrB: "form",
			wantOK:   true,
		},
		{
			name:   "both not convertable",
			a:      12,
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "one not convertable",
			a:      "change",
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotA, gotB, ok := tryBothToString(tt.a, tt.b)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantStrA, gotA, "should return correct value for a")
				assert.Equal(t, tt.wantStrB, gotB, "should return correct value for b")
			}
		})
	}
}

func Test_tryToTime(t *testing.T) {
	tests := []struct {
		v        interface{}
		wantTime time.Time
		wantOK   bool
	}{
		{
			v:      uint(368),
			wantOK: false,
		},
		{
			v:      int(-934),
			wantOK: false,
		},
		{
			v:      12 * time.Hour,
			wantOK: false,
		},
		{
			v:      "Hello World!",
			wantOK: false,
		},
		{
			v:        time.Time{},
			wantTime: time.Time{},
			wantOK:   true,
		},
		{
			v:        time.Date(2022, 11, 18, 12, 6, 39, 0, time.UTC),
			wantTime: time.Date(2022, 11, 18, 12, 6, 39, 0, time.UTC),
			wantOK:   true,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.v), tt.v), func(t *testing.T) {
			t.Parallel()
			got, ok := tryToTime(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantTime, got, "should return correct value")
			}
		})
	}
}

func Test_tryBothToTime(t *testing.T) {
	tests := []struct {
		name      string
		a         interface{}
		b         interface{}
		wantTimeA time.Time
		wantTimeB time.Time
		wantOK    bool
	}{
		{
			name:      "both times",
			a:         time.Date(2022, 11, 8, 12, 7, 4, 0, time.UTC),
			b:         time.Date(2022, 11, 8, 12, 7, 13, 0, time.UTC),
			wantTimeA: time.Date(2022, 11, 8, 12, 7, 4, 0, time.UTC),
			wantTimeB: time.Date(2022, 11, 8, 12, 7, 13, 0, time.UTC),
			wantOK:    true,
		},
		{
			name:   "both not convertable",
			a:      12,
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "one not convertable",
			a:      time.Now(),
			b:      "efficiency",
			wantOK: false,
		},
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotA, gotB, ok := tryBothToTime(tt.a, tt.b)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantTimeA, gotA, "should return correct value for a")
				assert.Equal(t, tt.wantTimeB, gotB, "should return correct value for b")
			}
		})
	}
}

func Test_tryToDuration(t *testing.T) {
	tests := []struct {
		v       interface{}
		wantDur time.Duration
		wantOK  bool
	}{
		{
			v:      uint(368),
			wantOK: false,
		},
		{
			v:      int(-934),
			wantOK: false,
		},
		{
			v:       12*time.Hour + 1*time.Minute,
			wantDur: 12*time.Hour + 1*time.Minute,
			wantOK:  true,
		},
		{
			v:      "Hello World!",
			wantOK: false,
		},
		{
			v:       time.Duration(0),
			wantDur: 0,
			wantOK:  true,
		},
		{
			v:      time.Date(2022, 11, 18, 12, 6, 39, 0, time.UTC),
			wantOK: false,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.v), tt.v), func(t *testing.T) {
			t.Parallel()
			got, ok := tryToDuration(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantDur, got, "should return correct value")
			}
		})
	}
}

func Test_tryBothToDuration(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		wantDurA time.Duration
		wantDurB time.Duration
		wantOK   bool
	}{
		{
			name:     "both durations",
			a:        12 * time.Hour,
			b:        4*time.Minute + 1*time.Second,
			wantDurA: 12 * time.Hour,
			wantDurB: 4*time.Minute + 1*time.Second,
			wantOK:   true,
		},
		{
			name:   "both not convertable",
			a:      12,
			b:      time.Now(),
			wantOK: false,
		},
		{
			name:   "one not convertable",
			a:      12 * time.Minute,
			b:      "efficiency",
			wantOK: false,
		},
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotA, gotB, ok := tryBothToDuration(tt.a, tt.b)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantDurA, gotA, "should return correct value for a")
				assert.Equal(t, tt.wantDurB, gotB, "should return correct value for b")
			}
		})
	}
}

func Test_typeOf(t *testing.T) {
	tests := []struct {
		v            interface{}
		wantTypeName typeName
		wantOK       bool
	}{
		{
			v:            uint(368),
			wantTypeName: typeUInt,
			wantOK:       true,
		},
		{
			v:            uint8(25),
			wantTypeName: typeUInt8,
			wantOK:       true,
		},
		{
			v:            uint16(487),
			wantTypeName: typeUInt16,
			wantOK:       true,
		},
		{
			v:            uint32(489),
			wantTypeName: typeUInt32,
			wantOK:       true,
		}, {
			v:            uint64(445),
			wantTypeName: typeUInt64,
			wantOK:       true,
		}, {
			v:            -934,
			wantTypeName: typeInt,
			wantOK:       true,
		},
		{
			v:            int8(-3),
			wantTypeName: typeInt8,
			wantOK:       true,
		},
		{
			v:            int16(24),
			wantTypeName: typeInt16,
			wantOK:       true,
		},
		{
			v:            int32(-849),
			wantTypeName: typeInt32,
			wantOK:       true,
		},
		{
			v:            int64(479),
			wantTypeName: typeInt64,
			wantOK:       true,
		},
		{
			v:            12 * time.Hour,
			wantTypeName: typeDuration,
			wantOK:       true,
		},
		{
			v:            "Hello World!",
			wantTypeName: typeString,
			wantOK:       true,
		},
		{
			v:            float32(1.5),
			wantTypeName: typeFloat32,
			wantOK:       true,
		},
		{
			v:            -444.2,
			wantTypeName: typeFloat64,
			wantOK:       true,
		},
		{
			v:            time.Now(),
			wantTypeName: typeTime,
			wantOK:       true,
		},
		{
			v:      struct{}{},
			wantOK: false,
		},
		{
			v:      nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.v)), func(t *testing.T) {
			t.Parallel()
			got, ok := typeOf(tt.v)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
			} else {
				assert.True(t, ok, "should not fail")
				assert.Equal(t, tt.wantTypeName, got, "should return correct value")
			}
		})
	}
}

func Test_reorderByType(t *testing.T) {
	order := []typeName{
		typeString,
		typeTime,
		typeDuration,
		typeUInt,
		typeUInt8,
		typeUInt16,
		typeUInt32,
		typeUInt64,
		typeInt,
		typeInt8,
		typeInt16,
		typeInt32,
		typeInt64,
		typeFloat32,
		typeFloat64,
	}
	tests := []struct {
		a           interface{}
		b           interface{}
		wantSwapped bool
		wantOK      bool
	}{
		{
			a:           uint8(13),
			b:           int64(2),
			wantSwapped: false,
			wantOK:      true,
		},
		{
			a:           int64(13),
			b:           uint8(2),
			wantSwapped: true,
			wantOK:      true,
		},
		{
			a:           int64(13),
			b:           uint8(233),
			wantSwapped: true,
			wantOK:      true,
		},
		{
			a:           4 * time.Hour,
			b:           time.Now(),
			wantSwapped: true,
			wantOK:      true,
		},
		{
			a:           12,
			b:           23 * time.Hour,
			wantSwapped: true,
			wantOK:      true,
		},
		{
			a:           time.Now(),
			b:           1 * time.Minute,
			wantSwapped: false,
			wantOK:      true,
		},
		{
			a:           "sir",
			b:           23,
			wantSwapped: false,
			wantOK:      true,
		},
		{
			a:           uint16(2),
			b:           "widow",
			wantSwapped: true,
			wantOK:      true,
		},
		{
			a:      nil,
			b:      23,
			wantOK: false,
		},
		{
			a:      struct{}{},
			b:      23,
			wantOK: false,
		},
		{
			a:      "move",
			b:      struct{}{},
			wantOK: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v %v", reflect.TypeOf(tt.a), reflect.TypeOf(tt.b)), func(t *testing.T) {
			t.Parallel()
			gotA, gotB, swapped, ok := reorderByType(tt.a, tt.b, order)
			if !tt.wantOK {
				assert.False(t, ok, "should fail")
				return
			}
			assert.True(t, ok, "should not fail")
			assert.Equal(t, tt.wantSwapped, swapped, "should swap correctly")
			if tt.wantSwapped {
				assert.Equal(t, tt.a, gotB, "should return swapped values")
				assert.Equal(t, tt.b, gotA, "should return swapped values")
			}
		})
	}
}
