package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{
			a:    1,
			b:    2,
			want: false,
		},
		{
			a:    23,
			b:    23,
			want: true,
		},
		{
			a:    uint(12),
			b:    83,
			want: false,
		}, {
			a:    int8(-23),
			b:    int32(-23),
			want: true,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: true,
		},
		{
			a:    11.0,
			b:    11,
			want: true,
		},
		{
			a:    2 * time.Hour,
			b:    120 * time.Minute,
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:    time.Now(),
			b:    120 * time.Hour,
			want: false,
		},
		{
			a:    23,
			b:    "least",
			want: false,
		},
		{
			a:    nil,
			b:    12,
			want: false,
		},
		{
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			a:    "organ",
			b:    "organ",
			want: true,
		},
		{
			a:    "table",
			b:    "sit",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v == %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Equal(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{
			a:    1,
			b:    2,
			want: true,
		},
		{
			a:    23,
			b:    23,
			want: false,
		},
		{
			a:    uint(12),
			b:    83,
			want: true,
		},
		{
			a:    int8(-23),
			b:    int32(-23),
			want: false,
		},
		{
			a:    int16(-23),
			b:    int8(-23),
			want: false,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: false,
		},
		{
			a:    11.1,
			b:    11,
			want: false,
		},
		{
			a:    1 * time.Hour,
			b:    120 * time.Minute,
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: false,
		},
		{
			a:    time.Date(2022, 18, 12, 1, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:       time.Now(),
			b:       120 * time.Hour,
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    "dry",
			b:    "dry",
			want: false,
		},
		{
			a:    "wire",
			b:    "pay",
			want: false,
		},
		{
			a:    "go",
			b:    "pay",
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v < %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Less(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestMore(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: true,
		},
		{
			a:    23,
			b:    23,
			want: false,
		},
		{
			a:    uint(12),
			b:    83,
			want: false,
		},
		{
			a:    int8(-23),
			b:    int32(-24),
			want: true,
		},
		{
			a:    int64(-23),
			b:    int32(-24),
			want: true,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: false,
		},
		{
			a:    11.1,
			b:    11,
			want: true,
		},
		{
			a:    3 * time.Hour,
			b:    120 * time.Minute,
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: false,
		},
		{
			a:    time.Date(2022, 18, 23, 1, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:       time.Now(),
			b:       120 * time.Hour,
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    "dry",
			b:    "dry",
			want: false,
		},
		{
			a:    "wire",
			b:    "pay",
			want: true,
		},
		{
			a:    "go",
			b:    "pay",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v > %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := More(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestLessOrEqual(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{
			a:    1,
			b:    2,
			want: true,
		},
		{
			a:    23,
			b:    23,
			want: true,
		},
		{
			a:    uint(12),
			b:    83,
			want: true,
		},
		{
			a:    int8(-23),
			b:    int32(-23),
			want: true,
		},
		{
			a:    int32(-23),
			b:    int8(-23),
			want: true,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: true,
		},
		{
			a:    11.1,
			b:    11,
			want: false,
		},
		{
			a:    1 * time.Hour,
			b:    120 * time.Minute,
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 1, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:       time.Now(),
			b:       120 * time.Hour,
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    "dry",
			b:    "dry",
			want: true,
		},
		{
			a:    "wire",
			b:    "pay",
			want: false,
		},
		{
			a:    "go",
			b:    "pay",
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v <= %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := LessOrEqual(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestMoreOrEqual(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: true,
		},
		{
			a:    23,
			b:    23,
			want: true,
		},
		{
			a:    uint(12),
			b:    83,
			want: false,
		},
		{
			a:    int8(-23),
			b:    int32(-24),
			want: true,
		},
		{
			a:    int32(2),
			b:    int8(8),
			want: false,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: true,
		},
		{
			a:    11.1,
			b:    11,
			want: true,
		},
		{
			a:    3 * time.Hour,
			b:    120 * time.Minute,
			want: true,
		},
		{
			a:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:    time.Date(2022, 18, 23, 1, 35, 25, 0, time.UTC),
			b:    time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: true,
		},
		{
			a:       time.Now(),
			b:       120 * time.Hour,
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    "dry",
			b:    "dry",
			want: true,
		},
		{
			a:    "wire",
			b:    "pay",
			want: true,
		},
		{
			a:    "go",
			b:    "pay",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v >= %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := MoreOrEqual(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: 3,
		},
		{
			a:    23,
			b:    23,
			want: 46,
		},
		{
			a:    int32(4),
			b:    int8(2),
			want: 6,
		},
		{
			a:    uint(12),
			b:    83,
			want: 95,
		}, {
			a:    int8(-23),
			b:    int32(-24),
			want: -47,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: 533 * 2,
		},
		{
			a:    11.1,
			b:    11,
			want: 22.1,
		},
		{
			a:    3 * time.Hour,
			b:    120 * time.Minute,
			want: 5 * time.Hour,
		},
		{
			a:    2*time.Minute + 4*time.Second,
			b:    int(time.Second),
			want: 2*time.Minute + 5*time.Second,
		},
		{
			a:    4 * time.Second,
			b:    float64(time.Second),
			want: 5 * time.Second,
		},
		{
			a:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			wantErr: true,
		},
		{
			a:    time.Date(2000, 2, 28, 22, 0, 0, 0, time.UTC),
			b:    2 * time.Hour,
			want: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			a:    2 * time.Hour,
			b:    time.Date(2000, 2, 28, 22, 0, 0, 0, time.UTC),
			want: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			a:    "a",
			b:    "b",
			want: "ab",
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v + %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Add(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: 1,
		},
		{
			a:    23,
			b:    23,
			want: 0,
		},
		{
			a:    uint(12),
			b:    83,
			want: -71,
		},
		{
			a:    int8(-23),
			b:    int32(-24),
			want: 1,
		},
		{
			a:    int32(-23),
			b:    int8(-24),
			want: 1,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: 0,
		},
		{
			a:    11.5,
			b:    12,
			want: -0.5,
		},
		{
			a:    3 * time.Hour,
			b:    120 * time.Minute,
			want: 1 * time.Hour,
		},
		{
			a: time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b: time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			want: time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).
				Sub(time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local)),
		},
		{
			a:    time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
			b:    2 * time.Hour,
			want: time.Date(2000, 2, 29, 22, 0, 0, 0, time.UTC),
		},
		{
			a:       "a",
			b:       "b",
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    4.5,
			b:    0.5,
			want: 4.0,
		},
		{
			a:       4 * time.Minute,
			b:       time.Now(),
			wantErr: true,
		},
		{
			a:    3 * time.Second,
			b:    int(1 * time.Second),
			want: 2 * time.Second,
		},
		{
			a:       4,
			b:       2 * time.Hour,
			wantErr: true,
		},
		{
			a:       4.0,
			b:       2 * time.Minute,
			wantErr: true,
		},
		{
			a:    7 * time.Hour,
			b:    float64(30 * time.Minute),
			want: 6*time.Hour + 30*time.Minute,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v - %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Subtract(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: 2 * 1,
		},
		{
			a:    23,
			b:    23,
			want: 23 * 23,
		},
		{
			a:    uint(12),
			b:    83,
			want: 12 * 83,
		}, {
			a:    int8(-23),
			b:    int32(-24),
			want: -23 * -24,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: 533 * 533,
		},
		{
			a:    11.5,
			b:    11,
			want: 11.5 * 11,
		},
		{
			a:       3 * time.Hour,
			b:       120 * time.Minute,
			wantErr: true,
		},
		{
			a:    3 * time.Hour,
			b:    2,
			want: 6 * time.Hour,
		},
		{
			a:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			wantErr: true,
		},
		{
			a:       time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
			b:       2 * time.Hour,
			wantErr: true,
		},
		{
			a:       "a",
			b:       "b",
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
		{
			a:    12 * time.Hour,
			b:    2.0,
			want: 24 * time.Hour,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v * %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Multiply(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			a:    2,
			b:    1,
			want: 2.0,
		},
		{
			a:    23,
			b:    23,
			want: 1.0,
		},
		{
			a:    uint(12),
			b:    8,
			want: 1.5,
		},
		{
			a:    uint32(8),
			b:    uint8(16),
			want: 0.5,
		},
		{
			a:    int8(-23),
			b:    int32(-24),
			want: float64(-23) / float64(-24),
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: 1.0,
		},
		{
			a:    11.5,
			b:    11,
			want: 11.5 / float64(11),
		},
		{
			a:    3 * time.Hour,
			b:    60 * time.Minute,
			want: 3.0,
		},
		{
			a:    4 * time.Hour,
			b:    2,
			want: 2 * time.Hour,
		},
		{
			a:       4,
			b:       2 * time.Hour,
			wantErr: true,
		},
		{
			a:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			wantErr: true,
		},
		{
			a:       time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
			b:       2 * time.Hour,
			wantErr: true,
		},
		{
			a:       "a",
			b:       "b",
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v / %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Divide(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}

func TestModulo(t *testing.T) {
	tests := []struct {
		a       interface{}
		b       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			a:    9,
			b:    3,
			want: 0,
		},
		{
			a:    23,
			b:    3,
			want: 2,
		},
		{
			a:    uint(12),
			b:    8,
			want: 4,
		},
		{
			a:    uint32(8),
			b:    uint8(16),
			want: 8,
		},
		{
			a:    -23,
			b:    int32(-2),
			want: -1,
		},
		{
			a:    uint16(533),
			b:    int32(533),
			want: 0,
		},
		{
			a:    11,
			b:    4,
			want: 3,
		},
		{
			a:       3 * time.Hour,
			b:       60 * time.Minute,
			wantErr: true,
		},
		{
			a:       4 * time.Hour,
			b:       2,
			wantErr: true,
		},
		{
			a:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC),
			b:       time.Date(2022, 18, 12, 12, 35, 25, 0, time.UTC).In(time.Local),
			wantErr: true,
		},
		{
			a:       time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
			b:       2 * time.Hour,
			wantErr: true,
		},
		{
			a:       "a",
			b:       "b",
			wantErr: true,
		},
		{
			a:       23,
			b:       "least",
			wantErr: true,
		},
		{
			a:       nil,
			b:       12,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%v %v / %v %v", reflect.TypeOf(tt.a), tt.a, reflect.TypeOf(tt.b), tt.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := assert.NotPanics
			if tt.wantErr {
				runner = assert.Panics
			}
			runner(t, func() {
				got := Modulo(tt.a, tt.b)
				if tt.wantErr {
					return
				}
				assert.Equal(t, tt.want, got, "should return correct value")
			})
		})
	}
}
