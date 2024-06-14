package types_test

import (
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"
	. "github.com/expr-lang/expr/types"
)

func TestType_Equal(t *testing.T) {
	tests := []struct {
		index string
		a, b  Type
		want  bool
	}{
		{"1", Int, Int, true},
		{"2", Int, Int8, false},
		{"3", Int, Uint, false},
		{"4", Int, Float, false},
		{"5", Int, String, false},
		{"6", Int, Bool, false},
		{"7", Int, Nil, false},
		{"8", Int, Array(Int), false},
		{"9", Int, Map{"foo": Int}, false},
		{"10", Int, StrictMap{"foo": Int}, false},
		{"11", Int, Array(Int), false},
		{"12", Array(Int), Array(Int), true},
		{"13", Array(Int), Array(Float), false},
		{"14", Map{"foo": Int}, Map{"foo": Int}, true},
		{"15", Map{"foo": Int}, Map{"foo": Float}, false},
		{"16", Map{"foo": Int}, StrictMap{"foo": Int}, false},
		{"17", StrictMap{"foo": Int}, StrictMap{"foo": Int}, true},
		{"18", StrictMap{"foo": Int}, StrictMap{"foo": Float}, false},
		{"19", Map{"foo": Map{"bar": Int}}, Map{"foo": Map{"bar": Int}}, true},
		{"20", Map{"foo": Map{"bar": Int}}, Map{"foo": Map{"bar": Float}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.index, func(t *testing.T) {
			if tt.want {
				require.True(t, tt.a.Equal(tt.b), tt.a.String()+" == "+tt.b.String())
			} else {
				require.False(t, tt.a.Equal(tt.b), tt.a.String()+" == "+tt.b.String())
			}
		})
	}
}
