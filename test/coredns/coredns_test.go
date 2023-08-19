package coredns_test

import (
	"context"
	"testing"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/test/coredns"
	"github.com/stretchr/testify/assert"
)

func TestCoreDNS(t *testing.T) {
	env := coredns.DefaultEnv(context.Background(), &coredns.Request{})

	tests := []struct {
		input string
	}{
		{`metadata('geoip/city/name') == 'Exampleshire'`},
		{`(type() == 'A' && name() == 'example.com') || client_ip() == '1.2.3.4'`},
		{`name() matches '^abc\\..*\\.example\\.com\\.$'`},
		{`type() in ['A', 'AAAA']`},
		{`incidr(client_ip(), '192.168.0.0/16')`},
		{`incidr(client_ip(), '127.0.0.0/24')`},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			_, err := expr.Compile(test.input, expr.Env(env), expr.DisableBuiltin("type"))
			assert.NoError(t, err)
		})
	}
}
