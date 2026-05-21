package runtime_test

import (
	"strconv"
	"testing"

	"github.com/expr-lang/expr/vm/runtime"
)

// BenchmarkIn benchmarks the `in` operator over the common slice shapes at
// representative list sizes. The interesting comparison is between the typed
// slice variants (which previously paid one heap alloc per element through
// reflect.Value.Index(i).Interface()) and the []any variant (which has always
// been zero-alloc per element because the slice's element type is interface).
//
// Run with:
//
//	go test -bench=BenchmarkIn -benchmem ./vm/runtime/
func BenchmarkIn(b *testing.B) {
	sizes := []int{8, 64, 256}

	for _, n := range sizes {
		// Plant a hit roughly halfway through so the loop's short-circuit
		// fires at the same position in every variant.
		strs := make([]string, n)
		anys := make([]any, n)
		for i := 0; i < n; i++ {
			s := strconv.Itoa(i)
			strs[i] = s
			anys[i] = s
		}
		strs[n/2] = "needle"
		anys[n/2] = "needle"

		b.Run("StringSlice/N="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !runtime.In("needle", strs) {
					b.Fatal("expected hit")
				}
			}
		})
		b.Run("AnySliceOfString/N="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !runtime.In("needle", anys) {
					b.Fatal("expected hit")
				}
			}
		})

		floats := make([]float64, n)
		floatAnys := make([]any, n)
		for i := 0; i < n; i++ {
			floats[i] = float64(i)
			floatAnys[i] = float64(i)
		}
		floats[n/2] = 99999.0
		floatAnys[n/2] = 99999.0

		b.Run("Float64Slice/N="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !runtime.In(99999.0, floats) {
					b.Fatal("expected hit")
				}
			}
		})

		ints := make([]int64, n)
		intAnys := make([]any, n)
		for i := 0; i < n; i++ {
			ints[i] = int64(i)
			intAnys[i] = int64(i)
		}
		ints[n/2] = 99999
		intAnys[n/2] = int64(99999)

		b.Run("Int64Slice/N="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if !runtime.In(int64(99999), ints) {
					b.Fatal("expected hit")
				}
			}
		})
	}
}
