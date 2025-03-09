package profile_test

import (
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/profile"
	"github.com/expr-lang/expr/vm"
)

func TestProfileExample(t *testing.T) {
	prg, err := expr.Compile(`(a + b) + (a + c) + func1() + func2()`,
		expr.Profile(),
		expr.Env(map[string]any{
			"a": int64(1),
			"b": int64(1),
			"c": int64(1),
		}),
		expr.Function("func1", func(params ...any) (any, error) {
			time.Sleep(time.Second)
			return 3, nil
		}),
		expr.Function("func2", func(params ...any) (any, error) {
			time.Sleep(time.Second)
			return 4, nil
		}),
	)
	if err != nil {
		t.Error(err)
	}
	out, err := expr.Run(prg, map[string]any{"a": int64(3), "b": int64(2), "c": int64(3)})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	t.Logf("%s", profile.GeneratePprofProfile(vm.GetSpan(prg), "./profile.pprof"))
}
