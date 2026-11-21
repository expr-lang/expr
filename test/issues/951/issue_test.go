package issue951

import (
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/internal/testify/require"
)

type Node interface {
	ID() string
}

type Base struct {
	Name string
}

func (b Base) ID() string { return b.Name }

type Container struct {
	Base
	Items []*Item
}

type Item struct {
	Kind  string
	Value string
}

type Wrapper struct {
	Node // embedded interface
}

type Proxy struct {
	*Wrapper
}

type Nodes []Node

func (ns Nodes) GetByID(id string) Node {
	for _, n := range ns {
		if n.ID() == id {
			return n
		}
	}
	return nil
}

func TestFieldAccessThroughEmbeddedInterface(t *testing.T) {
	container := &Container{
		Base: Base{Name: "test"},
		Items: []*Item{
			{Kind: "card", Value: "some_value"},
		},
	}
	proxy := &Proxy{
		Wrapper: &Wrapper{
			Node: container,
		},
	}

	tests := []struct {
		name   string
		expr   string
		env    any
		expect any
	}{
		{
			name:   "field through GetByID returning interface",
			expr:   `data.GetByID("test").Items[0].Value`,
			env:    map[string]any{"data": Nodes{proxy}},
			expect: "some_value",
		},
		{
			name:   "optional chaining with embedded interface",
			expr:   `data.GetByID("test")?.Items[0].Value`,
			env:    map[string]any{"data": Nodes{proxy}},
			expect: "some_value",
		},
		{
			name:   "optional chaining nil result",
			expr:   `data.GetByID("missing")?.Items`,
			env:    map[string]any{"data": Nodes{proxy}},
			expect: nil,
		},
		{
			name:   "promoted field through interface",
			expr:   `data.GetByID("test").Name`,
			env:    map[string]any{"data": Nodes{proxy}},
			expect: "test",
		},
		{
			name:   "method on interface still works",
			expr:   `data.GetByID("test").ID()`,
			env:    map[string]any{"data": Nodes{proxy}},
			expect: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expr.Eval(tt.expr, tt.env)
			require.NoError(t, err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestFieldAccessEmbeddedInterfaceNil(t *testing.T) {
	proxy := &Proxy{
		Wrapper: &Wrapper{
			Node: nil,
		},
	}

	_, err := expr.Eval(`Items[0].Value`, proxy)
	require.Error(t, err)
}
