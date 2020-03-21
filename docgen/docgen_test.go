package docgen_test

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"

	. "github.com/antonmedv/expr/docgen"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

type Tweet struct {
	Size    int
	Message string
}

type Env struct {
	Tweets []Tweet
	Config struct {
		MaxSize int32
	}
	Env map[string]interface{}
}

type Duration int

func (Duration) String() string {
	return ""
}

func (*Env) Duration(s string) Duration {
	return Duration(0)
}

func TestCreateDoc(t *testing.T) {
	Operators = nil
	Builtins = nil
	doc := CreateDoc(&Env{})
	expected := &Context{
		Variables: map[Identifier]*Type{
			"Tweets": {
				Kind: "array",
				Type: &Type{
					Kind: "struct",
					Name: "Tweet",
				},
			},
			"Config": {
				Kind: "struct",
				Fields: map[Identifier]*Type{
					"MaxSize": {Kind: "int"},
				},
			},
			"Env": {
				Kind: "map",
				Key:  &Type{Kind: "string"},
				Type: &Type{Kind: "any"},
			},
			"Duration": {
				Kind: "func",
				Arguments: []*Type{
					{Kind: "string"},
				},
				Return: &Type{Kind: "struct", Name: "Duration"},
			},
		},
		Types: map[TypeName]*Type{
			"Tweet": {
				Kind: "struct",
				Fields: map[Identifier]*Type{
					"Size":    {Kind: "int"},
					"Message": {Kind: "string"},
				},
			},
			"Duration": {
				Kind: "struct",
				Fields: map[Identifier]*Type{
					"String": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "string",
						},
					},
				},
			},
		},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(doc))
}

func TestCreateDoc_FromMap(t *testing.T) {
	env := map[string]interface{}{
		"Tweets": []*Tweet{},
		"Config": struct {
			MaxSize int
		}{},
		"Max": math.Max,
	}
	Operators = nil
	Builtins = nil
	doc := CreateDoc(env)
	expected := &Context{
		Variables: map[Identifier]*Type{
			"Tweets": {
				Kind: "array",
				Type: &Type{
					Kind: "struct",
					Name: "Tweet",
				},
			},
			"Config": {
				Kind: "struct",
				Fields: map[Identifier]*Type{
					"MaxSize": {Kind: "int"},
				},
			},
			"Max": {
				Kind: "func",
				Arguments: []*Type{
					{Kind: "float"},
					{Kind: "float"},
				},
				Return: &Type{Kind: "float"},
			},
		},
		Types: map[TypeName]*Type{
			"Tweet": {
				Kind: "struct",
				Fields: map[Identifier]*Type{
					"Size":    {Kind: "int"},
					"Message": {Kind: "string"},
				},
			},
		},
	}

	require.Equal(t, litter.Sdump(expected), litter.Sdump(doc))
}

func TestContext_Markdown(t *testing.T) {
	doc := CreateDoc(&Env{})
	md := doc.Markdown()
	require.True(t, len(md) > 0)
}
