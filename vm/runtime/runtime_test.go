package runtime

import (
	"reflect"
	"testing"

	"github.com/expr-lang/expr/internal/testify/require"
)

type Namer interface {
	Name() string
}

type IntHolder interface {
	Int() int
}

type ConcreteWithName struct {
	Title string
}

func (ConcreteWithName) Name() string { return "" }

type ConcreteWithSkippedField struct {
	Title string `expr:"-"`
}

func (ConcreteWithSkippedField) Name() string { return "" }

type ConcreteEmptyStruct struct{}

func (ConcreteEmptyStruct) Name() string { return "" }

type ConcreteInt int

func (ConcreteInt) Int() int { return 0 }

type ConcreteWithEmbeddedInterface struct {
	Namer
}

func (ConcreteWithEmbeddedInterface) Int() int { return 0 }

type EmbeddedInterfaceOnly struct {
	Namer
}

type EmbeddedNilPointerOnly struct {
	*ConcreteWithName
}

type EmbeddedStructWithInterface struct {
	EmbeddedInterfaceOnly
}

type EmbeddedIntHolder struct {
	IntHolder
}

type PlainStruct struct {
	Title string
}

func TestFetchFromEmbeddedInterfaces(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		fieldName string
		want      any
		ok        bool
	}{
		{
			name:      "no anonymous fields",
			input:     PlainStruct{Title: "ignored"},
			fieldName: "Title",
			ok:        false,
		},
		{
			name: "embedded interface with field on concrete struct",
			input: EmbeddedInterfaceOnly{
				Namer: ConcreteWithName{Title: "hello"},
			},
			fieldName: "Title",
			want:      "hello",
			ok:        true,
		},
		{
			name: "embedded interface, concrete missing field",
			input: EmbeddedInterfaceOnly{
				Namer: ConcreteWithName{Title: "hello"},
			},
			fieldName: "Missing",
			ok:        false,
		},
		{
			name: "embedded interface holding pointer to struct",
			input: EmbeddedInterfaceOnly{
				Namer: &ConcreteWithName{Title: "pointer"},
			},
			fieldName: "Title",
			want:      "pointer",
			ok:        true,
		},
		{
			name:      "embedded interface with nil concrete value",
			input:     EmbeddedInterfaceOnly{Namer: nil},
			fieldName: "Title",
			ok:        false,
		},
		{
			name:      "embedded nil pointer to struct",
			input:     EmbeddedNilPointerOnly{ConcreteWithName: nil},
			fieldName: "Title",
			ok:        false,
		},
		{
			name: "embedded struct containing embedded interface with field",
			input: EmbeddedStructWithInterface{
				EmbeddedInterfaceOnly: EmbeddedInterfaceOnly{
					Namer: ConcreteWithName{Title: "nested"},
				},
			},
			fieldName: "Title",
			want:      "nested",
			ok:        true,
		},
		{
			name: "embedded interface whose concrete embeds another interface",
			input: EmbeddedIntHolder{
				IntHolder: ConcreteWithEmbeddedInterface{
					Namer: ConcreteWithName{Title: "deep"},
				},
			},
			fieldName: "Title",
			want:      "deep",
			ok:        true,
		},
		{
			name: "embedded interface with non-struct concrete value",
			input: EmbeddedIntHolder{
				IntHolder: ConcreteInt(5),
			},
			fieldName: "Title",
			ok:        false,
		},
		{
			name: "field is skipped via expr:\"-\" tag",
			input: EmbeddedInterfaceOnly{
				Namer: ConcreteWithSkippedField{Title: "hidden"},
			},
			fieldName: "Title",
			ok:        false,
		},
		{
			name: "embedded interface with empty concrete struct, recurses to nothing",
			input: EmbeddedInterfaceOnly{
				Namer: ConcreteEmptyStruct{},
			},
			fieldName: "Title",
			ok:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := fetchFromEmbeddedInterfaces(reflect.ValueOf(tt.input), tt.fieldName)
			require.Equal(t, tt.ok, ok)
			if tt.ok {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
