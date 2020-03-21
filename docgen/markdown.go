package docgen

import (
	"fmt"
	"strings"
)

func (c *Context) Markdown() string {
	out := `### Variables
| Name | Type |
|------|------|
`
	for ident, v := range c.Variables {
		if v.Kind == "func" {
			continue
		}
		if v.Kind == "operator" {
			continue
		}
		out += fmt.Sprintf("| %v | %v |\n", ident, link(v))
	}

	out += `
### Functions
| Name | Return type |
|------|-------------|
`
	for ident, v := range c.Variables {
		if v.Kind == "func" {
			args := make([]string, len(v.Arguments))
			for i, arg := range v.Arguments {
				args[i] = link(arg)
			}
			out += fmt.Sprintf("| %v(%v) | %v |\n", ident, strings.Join(args, ", "), link(v.Return))
		}
	}

	out += "\n### Types\n"
	for name, t := range c.Types {
		out += fmt.Sprintf("#### %v\n", name)
		out += fields(t)
		out += "\n"
	}

	return out
}

func link(t *Type) string {
	if t == nil {
		return "nil"
	}
	if t.Name != "" {
		return fmt.Sprintf("[%v](#%v)", t.Name, t.Name)
	}
	if t.Kind == "array" {
		return fmt.Sprintf("array(%v)", link(t.Type))
	}
	if t.Kind == "map" {
		return fmt.Sprintf("map(%v => %v)", link(t.Key), link(t.Type))
	}
	return fmt.Sprintf("`%v`", t.Kind)
}

func fields(t *Type) string {
	out := ""
	foundFields := false
	for ident, v := range t.Fields {
		if v.Kind != "func" {
			if !foundFields {
				out += "| Field | Type |\n|---|---|\n"
			}
			foundFields = true

			out += fmt.Sprintf("| %v | %v |\n", ident, link(v))
		}
	}
	foundMethod := false
	for ident, v := range t.Fields {
		if v.Kind == "func" {
			if !foundMethod {
				out += "\n| Method | Returns |\n|---|---|\n"
			}
			foundMethod = true

			args := make([]string, len(v.Arguments))
			for i, arg := range v.Arguments {
				args[i] = link(arg)
			}
			out += fmt.Sprintf("| %v(%v) | %v |\n", ident, strings.Join(args, ", "), link(v.Return))
		}
	}
	return out
}
