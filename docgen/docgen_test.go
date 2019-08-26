package docgen_test

import (
	. "github.com/antonmedv/expr/docgen"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func (*Env) Duration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestCreateDoc(t *testing.T) {
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
					"Hours": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "float",
						},
					},
					"Minutes": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "float",
						},
					},
					"Nanoseconds": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "int",
						},
					},
					"Round": {
						Kind: "func",
						Arguments: []*Type{
							{
								Name: "Duration",
								Kind: "struct",
							},
						},
						Return: &Type{
							Name: "Duration",
							Kind: "struct",
						},
					},
					"Seconds": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "float",
						},
					},
					"String": {
						Kind:      "func",
						Arguments: []*Type{},
						Return: &Type{
							Kind: "string",
						},
					},
					"Truncate": {
						Kind: "func",
						Arguments: []*Type{
							{
								Name: "Duration",
								Kind: "struct",
							},
						},
						Return: &Type{
							Name: "Duration",
							Kind: "struct",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, litter.Sdump(expected), litter.Sdump(doc))
}
