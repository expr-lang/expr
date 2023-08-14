package time_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	testTime := time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	testDuration := time.Duration(1)

	tests := []struct {
		a       interface{}
		b       interface{}
		op      string
		want    interface{}
		wantErr bool
	}{
		{a: testTime, b: testTime, op: "<", wantErr: false, want: false},
		{a: testTime, b: testTime, op: ">", wantErr: false, want: false},
		{a: testTime, b: testTime, op: "<=", wantErr: false, want: true},
		{a: testTime, b: testTime, op: ">=", wantErr: false, want: true},
		{a: testTime, b: testTime, op: "==", wantErr: false, want: true},
		{a: testTime, b: testTime, op: "!=", wantErr: false, want: false},
		{a: testTime, b: testTime, op: "-", wantErr: false},
		{a: testTime, b: testDuration, op: "+", wantErr: false},
		{a: testTime, b: testDuration, op: "-", wantErr: false},

		// error cases
		{a: testTime, b: int64(1), op: "<", wantErr: true},
		{a: testTime, b: float64(1), op: "<", wantErr: true},
		{a: testTime, b: testDuration, op: "<", wantErr: true},

		{a: testTime, b: int64(1), op: ">", wantErr: true},
		{a: testTime, b: float64(1), op: ">", wantErr: true},
		{a: testTime, b: testDuration, op: ">", wantErr: true},

		{a: testTime, b: int64(1), op: "<=", wantErr: true},
		{a: testTime, b: float64(1), op: "<=", wantErr: true},
		{a: testTime, b: testDuration, op: "<=", wantErr: true},

		{a: testTime, b: int64(1), op: ">=", wantErr: true},
		{a: testTime, b: float64(1), op: ">=", wantErr: true},
		{a: testTime, b: testDuration, op: ">=", wantErr: true},

		{a: testTime, b: int64(1), op: "==", wantErr: false, want: false},
		{a: testTime, b: float64(1), op: "==", wantErr: false, want: false},
		{a: testTime, b: testDuration, op: "==", wantErr: false, want: false},

		{a: testTime, b: int64(1), op: "!=", wantErr: false, want: true},
		{a: testTime, b: float64(1), op: "!=", wantErr: false, want: true},
		{a: testTime, b: testDuration, op: "!=", wantErr: false, want: true},

		{a: testTime, b: int64(1), op: "-", wantErr: true},
		{a: testTime, b: float64(1), op: "-", wantErr: true},

		{a: testTime, b: testTime, op: "+", wantErr: true},
		{a: testTime, b: int64(1), op: "+", wantErr: true},
		{a: testTime, b: float64(1), op: "+", wantErr: true},
		{a: testDuration, b: testTime, op: "+", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("time helper test `%T %s %T`", tt.a, tt.op, tt.b), func(t *testing.T) {
			input := fmt.Sprintf("a %v b", tt.op)
			env := map[string]interface{}{
				"a": tt.a,
				"b": tt.b,
			}

			config := conf.CreateNew()

			tree, err := parser.Parse(input)
			require.NoError(t, err)

			_, err = checker.Check(tree, config)
			require.NoError(t, err)

			program, err := compiler.Compile(tree, config)
			require.NoError(t, err)

			got, err := vm.Run(program, env)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.want != nil {
					require.Equal(t, tt.want, got)
				}
			}
		})
	}
}
