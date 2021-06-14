package vm_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/compiler"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"github.com/stretchr/testify/require"
)

const (
	assertDelta = 0.001
)

func TestRun_nanmin(t *testing.T) {
	var input = `nanmin([-1, 1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, float64(-1), out)
}

func TestRun_nanmax(t *testing.T) {
	var input = `nanmax([-1, 1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, float64(2), out)
}

func TestRun_nansum(t *testing.T) {
	var input = `nansum([-1, 0, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, float64(1), out)
}

func TestRun_nanprod(t *testing.T) {
	var input = `nanprod([-1, 1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, float64(-2), out)
}

func TestRun_nanmean(t *testing.T) {
	var input = `nanmean([-1, 0, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDelta(t, float64(0.3333), out, assertDelta)
}

func TestRun_nanstd(t *testing.T) {
	var input = `nanstd([-1, 0, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDelta(t, float64(1.247219128924647), out, assertDelta)
}

func TestRun_abs(t *testing.T) {
	var input = `abs([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, []float64{1, 1}, out)
}

func TestRun_acos(t *testing.T) {
	var input = `acos([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{3.141592653589793, 0}, out, assertDelta)
}

func TestRun_acosh(t *testing.T) {
	var input = `acosh([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0, 1.3169578969248166}, out, assertDelta)
}

func TestRun_asin(t *testing.T) {
	var input = `asin([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1.5707963267948966, 1.5707963267948966}, out, assertDelta)
}

func TestRun_asinh(t *testing.T) {
	var input = `asinh([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.881373587019543, 0.881373587019543}, out, assertDelta)
}

func TestRun_atan(t *testing.T) {
	var input = `atan([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.7853981633974483, 0.7853981633974483}, out, assertDelta)
}

func TestRun_atanh(t *testing.T) {
	var input = `atanh([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(-1), math.Inf(1)}, out, assertDelta)
}

func TestRun_cos(t *testing.T) {
	var input = `cos([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0.5403023058681398, 0.5403023058681398}, out, assertDelta)
}

func TestRun_cosh(t *testing.T) {
	var input = `cosh([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{1.5430806348152437, 1.5430806348152437}, out, assertDelta)
}

func TestRun_sin(t *testing.T) {
	var input = `sin([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.8414709848078965, 0.8414709848078965}, out, assertDelta)
}

func TestRun_sinh(t *testing.T) {
	var input = `sinh([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1.1752011936438014, 1.1752011936438014}, out, assertDelta)
}

func TestRun_tan(t *testing.T) {
	var input = `tan([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1.557407724654902, 1.557407724654902}, out, assertDelta)
}

func TestRun_tanh(t *testing.T) {
	var input = `tanh([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.7615941559557649, 0.7615941559557649}, out, assertDelta)
}

func TestRun_ceil(t *testing.T) {
	var input = `ceil([-1.1, 1.1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, []float64{-1, 2}, out)
}

func TestRun_floor(t *testing.T) {
	var input = `floor([-1.1, 1.1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.Equal(t, []float64{-2, 1}, out)
}

func TestRun_cbrt(t *testing.T) {
	var input = `cbrt([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1, 1}, out, assertDelta)
}

func TestRun_erf(t *testing.T) {
	var input = `erf([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.8427007929497149, 0.8427007929497149}, out, assertDelta)
}

func TestRun_erfc(t *testing.T) {
	var input = `erfc([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{1.8427007929497148, 0.15729920705028513}, out, assertDelta)
}

func TestRun_erfcinv(t *testing.T) {
	var input = `erfcinv([0, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(1), 0}, out, assertDelta)
}

func TestRun_erfinv(t *testing.T) {
	var input = `erfinv([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(-1), math.Inf(1)}, out, assertDelta)
}

func TestRun_exp(t *testing.T) {
	var input = `exp([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0.36787944117144233, 2.718281828459045}, out, assertDelta)
}

func TestRun_exp2(t *testing.T) {
	var input = `exp2([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{2, 4}, out, assertDelta)
}

func TestRun_expm1(t *testing.T) {
	var input = `expm1([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.6321205588285577, 1.718281828459045}, out, assertDelta)
}

func TestRun_gamma(t *testing.T) {
	var input = `gamma([0, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(1), 1}, out, assertDelta)
}

func TestRun_j0(t *testing.T) {
	var input = `j0([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0.7651976865579666, 0.7651976865579666}, out, assertDelta)
}

func TestRun_j1(t *testing.T) {
	var input = `j1([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-0.4400505857449335, 0.4400505857449335}, out, assertDelta)
}

func TestRun_log(t *testing.T) {
	var input = `log([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0, 0.6931471805599453}, out, assertDelta)
}

func TestRun_log10(t *testing.T) {
	var input = `log10([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0, 0.3010299956639812}, out, assertDelta)
}

func TestRun_log1p(t *testing.T) {
	var input = `log1p([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(-1), 0.6931471805599453}, out, assertDelta)
}

func TestRun_log2(t *testing.T) {
	var input = `log2([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0, 1}, out, assertDelta)
}

func TestRun_logb(t *testing.T) {
	var input = `logb([-1, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{0, 0}, out, assertDelta)
}

func TestRun_round(t *testing.T) {
	var input = `round([-1.1, 1.1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1, 1}, out, assertDelta)
}

func TestRun_roundtoeven(t *testing.T) {
	var input = `roundtoeven([-1.1, 1.1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1, 1}, out, assertDelta)
}

func TestRun_sqrt(t *testing.T) {
	var input = `sqrt([1, 2])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{1, 1.4142135623730951}, out, assertDelta)
}

func TestRun_trunc(t *testing.T) {
	var input = `trunc([-1.1, 1.1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{-1, 1}, out, assertDelta)
}

func TestRun_y0(t *testing.T) {
	var input = `y0([0, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(-1), 0.08825696421567697}, out, assertDelta)
}

func TestRun_y1(t *testing.T) {
	var input = `y1([0, 1])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)
	require.InDeltaSlice(t, []float64{math.Inf(-1), -0.7812128213002887}, out, assertDelta)
}

func TestRun_minimum(t *testing.T) {
	env := map[string]interface{}{
		"a": []float32{1, 2},
	}
	var input = `minimum(a, [0, 4])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, []float64{0, 2}, out)
}

func TestRun_maximum(t *testing.T) {
	env := map[string]interface{}{
		"a": []float32{1, 2},
	}
	var input = `maximum(a, [0, 4])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, []float64{1, 4}, out)
}

func TestRun_mod(t *testing.T) {
	env := map[string]interface{}{
		"a": []float32{1, 2},
	}
	var input = `mod(a, [0, 4])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, []float64{0, 0}, out)
}

func TestRun_pow(t *testing.T) {
	env := map[string]interface{}{
		"a": []float32{1, 2},
	}
	var input = `pow(a, [0, 4])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, []float64{0, 16}, out)
}

func TestRun_remainder(t *testing.T) {
	env := map[string]interface{}{
		"a": []float32{1, 2},
	}
	var input = `remainder(a, [0, 4])`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	out, err := vm.Run(program, env)
	require.NoError(t, err)
	require.Equal(t, []float64{0, 0}, out)
}

func TestRun_debug(t *testing.T) {
	var input = `[1, 2, 3]`

	node, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(node, nil)
	require.NoError(t, err)

	_, err = vm.Run(program, nil)
	require.NoError(t, err)
}

func TestRun_cast(t *testing.T) {
	input := `1`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, &conf.Config{Expect: reflect.Float64})
	require.NoError(t, err)

	out, err := vm.Run(program, nil)
	require.NoError(t, err)

	require.Equal(t, float64(1), out)
}

func TestRun_helpers(t *testing.T) {
	values := []interface{}{
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		float32(1),
		float64(1),
	}
	ops := []string{"+", "-", "*", "/", "%", "==", ">=", "<=", "<", ">"}

	for _, a := range values {
		for _, b := range values {
			for _, op := range ops {

				if op == "%" {
					switch a.(type) {
					case float32, float64:
						continue
					}
					switch b.(type) {
					case float32, float64:
						continue
					}
				}

				input := fmt.Sprintf("a %v b", op)
				env := map[string]interface{}{
					"a": a,
					"b": b,
				}

				tree, err := parser.Parse(input)
				require.NoError(t, err)

				_, err = checker.Check(tree, nil)
				require.NoError(t, err)

				program, err := compiler.Compile(tree, nil)
				require.NoError(t, err)

				_, err = vm.Run(program, env)
				require.NoError(t, err)
			}
		}
	}
}

func checkVecAt(t *testing.T, expected float64, actual interface{}) {
	require.EqualValues(t, expected, actual)
}

func checkVec(t *testing.T, expected []float64, actual interface{}) {
	switch vec := actual.(type) {
	case []uint:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []uint8:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []uint16:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []uint32:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []uint64:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []int:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []int8:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []int16:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []int32:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []int64:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []float32:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	case []float64:
		for j, val := range vec {
			checkVecAt(t, expected[j], val)
		}
	}
}

func testScalar(t *testing.T, op string, value int, expected []float64) {
	singletons := []interface{}{
		uint(value),
		uint8(value),
		uint16(value),
		uint32(value),
		uint64(value),
		int(value),
		int8(value),
		int16(value),
		int32(value),
		int64(value),
		float32(value),
		float64(value),
	}

	values := []interface{}{
		[]uint{1, 2},
		[]uint8{1, 2},
		[]uint16{1, 2},
		[]uint32{1, 2},
		[]uint64{1, 2},
		[]int{1, 2},
		[]int8{1, 2},
		[]int16{1, 2},
		[]int32{1, 2},
		[]int64{1, 2},
		[]float32{1, 2},
		[]float64{1, 2},
	}
	for _, s := range singletons {
		for _, vec := range values {
			input := fmt.Sprintf("vec %v s", op)
			env := map[string]interface{}{
				"s":   s,
				"vec": vec,
			}

			tree, err := parser.Parse(input)
			require.NoError(t, err)

			_, err = checker.Check(tree, nil)
			require.NoError(t, err)

			program, err := compiler.Compile(tree, nil)
			require.NoError(t, err)

			out, err := vm.Run(program, env)
			require.NoError(t, err)

			checkVec(t, expected, out)
		}
	}
}

func testVec(t *testing.T, op string, expected []float64) {
	a := []interface{}{
		[]uint{1, 2},
		[]uint8{1, 2},
		[]uint16{1, 2},
		[]uint32{1, 2},
		[]uint64{1, 2},
		[]int{1, 2},
		[]int8{1, 2},
		[]int16{1, 2},
		[]int32{1, 2},
		[]int64{1, 2},
		[]float32{1, 2},
		[]float64{1, 2},
	}
	b := []interface{}{
		[]uint{1, 2},
		[]uint8{1, 2},
		[]uint16{1, 2},
		[]uint32{1, 2},
		[]uint64{1, 2},
		[]int{1, 2},
		[]int8{1, 2},
		[]int16{1, 2},
		[]int32{1, 2},
		[]int64{1, 2},
		[]float32{1, 2},
		[]float64{1, 2},
	}

	for _, x := range a {
		for _, y := range b {
			input := fmt.Sprintf("x %v y", op)
			env := map[string]interface{}{
				"x": x,
				"y": y,
			}

			tree, err := parser.Parse(input)
			require.NoError(t, err)

			_, err = checker.Check(tree, nil)
			require.NoError(t, err)

			program, err := compiler.Compile(tree, nil)
			require.NoError(t, err)

			out, err := vm.Run(program, env)
			require.NoError(t, err)
			checkVec(t, expected, out)
		}
	}
}

func TestRun_AddScalar(t *testing.T) {
	expected := []float64{2, 3}
	testScalar(t, "+", 1, expected)
}

func TestRun_SubScalar(t *testing.T) {
	expected := []float64{0, 1}
	testScalar(t, "-", 1, expected)
}

func TestRun_MulScalar(t *testing.T) {
	expected := []float64{2, 4}
	testScalar(t, "*", 2, expected)
}

func TestRun_DivScalar(t *testing.T) {
	expected := []float64{1, 2}
	testScalar(t, "/", 1, expected)
}

func TestRun_AddVec(t *testing.T) {
	expected := []float64{2, 4}
	testVec(t, "+", expected)
}

func TestRun_SubVec(t *testing.T) {
	expected := []float64{0, 0}
	testVec(t, "-", expected)
}

func TestRun_MulVec(t *testing.T) {
	expected := []float64{1, 4}
	testVec(t, "*", expected)
}

func TestRun_DivVec(t *testing.T) {
	expected := []float64{1, 1}
	testVec(t, "/", expected)
}

func TestRun_memory_budget(t *testing.T) {
	input := `map(1..100, {map(1..100, {map(1..100, {0})})})`

	tree, err := parser.Parse(input)
	require.NoError(t, err)

	program, err := compiler.Compile(tree, nil)
	require.NoError(t, err)

	_, err = vm.Run(program, nil)
	require.Error(t, err)
}
