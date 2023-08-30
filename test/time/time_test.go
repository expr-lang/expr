package time_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/antonmedv/expr"
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
		a       any
		b       any
		op      string
		want    any
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
			env := map[string]any{
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

func TestTime_duration(t *testing.T) {
	env := map[string]any{
		"foo": time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
	}
	program, err := expr.Compile(`now() - duration("1h") < now() && foo + duration("24h") < now()`, expr.Env(env))
	require.NoError(t, err)

	output, err := expr.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, true, output)
}

func TestTime_date(t *testing.T) {
	var tests = []struct {
		input string
		want  time.Time
	}{
		{
			`date('2017-10-23')`,
			time.Date(2017, 10, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			`date('24.11.1987 20:30', "02.01.2006 15:04", "Europe/Zurich")`,
			time.Date(1987, 11, 24, 20, 30, 0, 0, time.FixedZone("Europe/Zurich", 3600)),
		},
		{
			`date('24.11.1987 20:30 MSK', "02.01.2006 15:04 MST", "Europe/Zurich")`,
			time.Date(1987, 11, 24, 20, 30, 0, 0, time.FixedZone("MSK", 0)),
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program, err := expr.Compile(test.input)
			require.NoError(t, err)

			output, err := expr.Run(program, nil)
			require.NoError(t, err)
			require.Truef(t, test.want.Equal(output.(time.Time)), "want %v, got %v", test.want, output)
		})
	}
}
