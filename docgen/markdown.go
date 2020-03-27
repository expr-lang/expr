package docgen

import (
	"fmt"
	"sort"
	"strings"
)

func (c *Context) Markdown() string {
	var variables []string
	for name := range c.Variables {
		variables = append(variables, string(name))
	}

	var types []string
	for name := range c.Types {
		types = append(types, string(name))
	}

	sort.Strings(variables)
	sort.Strings(types)

	out := `### Variables
| Name | Type |
|------|------|
`
	for _, name := range variables {
		v := c.Variables[Identifier(name)]
		if v.Kind == "func" {
			continue
		}
		if v.Kind == "operator" {
			continue
		}
		out += fmt.Sprintf("| %v | %v |\n", name, link(v))
	}

	out += `
### Functions
| Name | Return type |
|------|-------------|
`
	for _, name := range variables {
		v := c.Variables[Identifier(name)]
		if v.Kind == "func" {
			args := make([]string, len(v.Arguments))
			for i, arg := range v.Arguments {
				args[i] = link(arg)
			}
			out += fmt.Sprintf("| %v(%v) | %v |\n", name, strings.Join(args, ", "), link(v.Return))
		}
	}

	out += "\n### Types\n"
	for _, name := range types {
		t := c.Types[TypeName(name)]
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
	var fields []string
	for field := range t.Fields {
		fields = append(fields, string(field))
	}
	sort.Strings(fields)

	out := ""
	foundFields := false
	for _, name := range fields {
		v := t.Fields[Identifier(name)]
		if v.Kind != "func" {
			if !foundFields {
				out += "| Field | Type |\n|---|---|\n"
			}
			foundFields = true

			out += fmt.Sprintf("| %v | %v |\n", name, link(v))
		}
	}
	foundMethod := false
	for _, name := range fields {
		v := t.Fields[Identifier(name)]
		if v.Kind == "func" {
			if !foundMethod {
				out += "\n| Method | Returns |\n|---|---|\n"
			}
			foundMethod = true

			args := make([]string, len(v.Arguments))
			for i, arg := range v.Arguments {
				args[i] = link(arg)
			}
			out += fmt.Sprintf("| %v(%v) | %v |\n", name, strings.Join(args, ", "), link(v.Return))
		}
	}
	return out
}
