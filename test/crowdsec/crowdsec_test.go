package crowdsec_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/test/crowdsec"
	"github.com/stretchr/testify/require"
)

func TestCrowdsec(t *testing.T) {
	b, err := os.ReadFile("../../testdata/crowdsec.json")
	require.NoError(t, err)

	var examples []string
	err = json.Unmarshal(b, &examples)
	require.NoError(t, err)

	env := map[string]any{
		"evt": &crowdsec.Event{},
	}

	var opt = []expr.Option{
		expr.Env(env),
	}
	for _, fn := range crowdsec.CustomFunctions {
		opt = append(
			opt,
			expr.Function(
				fn.Name,
				func(params ...any) (any, error) {
					return nil, nil
				},
				fn.Func...,
			),
		)
	}

	for _, line := range examples {
		t.Run(line, func(t *testing.T) {
			_, err = expr.Compile(line, opt...)
			require.NoError(t, err)
		})
	}
}
