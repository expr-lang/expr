package crowdsec

import (
	"time"
)

var CustomFunctions = []struct {
	Name string
	Func []any
}{
	{
		Name: "Distance",
		Func: []any{
			new(func(string, string, string, string) (float64, error)),
		},
	},
	{
		Name: "GetFromStash",
		Func: []any{
			new(func(string, string) (string, error)),
		},
	},
	{
		Name: "Atof",
		Func: []any{
			new(func(string) float64),
		},
	},
	{
		Name: "JsonExtract",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "JsonExtractUnescape",
		Func: []any{
			new(func(string, ...string) string),
		},
	},
	{
		Name: "JsonExtractLib",
		Func: []any{
			new(func(string, ...string) string),
		},
	},
	{
		Name: "JsonExtractSlice",
		Func: []any{
			new(func(string, string) []any),
		},
	},
	{
		Name: "JsonExtractObject",
		Func: []any{
			new(func(string, string) map[string]any),
		},
	},
	{
		Name: "ToJsonString",
		Func: []any{
			new(func(any) string),
		},
	},
	{
		Name: "File",
		Func: []any{
			new(func(string) []string),
		},
	},
	{
		Name: "RegexpInFile",
		Func: []any{
			new(func(string, string) bool),
		},
	},
	{
		Name: "Upper",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "Lower",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "IpInRange",
		Func: []any{
			new(func(string, string) bool),
		},
	},
	{
		Name: "TimeNow",
		Func: []any{
			new(func() string),
		},
	},
	{
		Name: "ParseUri",
		Func: []any{
			new(func(string) map[string][]string),
		},
	},
	{
		Name: "PathUnescape",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "QueryUnescape",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "PathEscape",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "QueryEscape",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "XMLGetAttributeValue",
		Func: []any{
			new(func(string, string, string) string),
		},
	},
	{
		Name: "XMLGetNodeValue",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "IpToRange",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "IsIPV6",
		Func: []any{
			new(func(string) bool),
		},
	},
	{
		Name: "IsIPV4",
		Func: []any{
			new(func(string) bool),
		},
	},
	{
		Name: "IsIP",
		Func: []any{
			new(func(string) bool),
		},
	},
	{
		Name: "LookupHost",
		Func: []any{
			new(func(string) []string),
		},
	},
	{
		Name: "GetDecisionsCount",
		Func: []any{
			new(func(string) int),
		},
	},
	{
		Name: "GetDecisionsSinceCount",
		Func: []any{
			new(func(string, string) int),
		},
	},
	{
		Name: "Sprintf",
		Func: []any{
			new(func(string, ...any) string),
		},
	},
	{
		Name: "ParseUnix",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "SetInStash",
		Func: []any{
			new(func(string, string, string, *time.Duration) error),
		},
	},
	{
		Name: "Fields",
		Func: []any{
			new(func(string) []string),
		},
	},
	{
		Name: "Index",
		Func: []any{
			new(func(string, string) int),
		},
	},
	{
		Name: "IndexAny",
		Func: []any{
			new(func(string, string) int),
		},
	},
	{
		Name: "Join",
		Func: []any{
			new(func([]string, string) string),
		},
	},
	{
		Name: "Split",
		Func: []any{
			new(func(string, string) []string),
		},
	},
	{
		Name: "SplitAfter",
		Func: []any{
			new(func(string, string) []string),
		},
	},
	{
		Name: "SplitAfterN",
		Func: []any{
			new(func(string, string, int) []string),
		},
	},
	{
		Name: "SplitN",
		Func: []any{
			new(func(string, string, int) []string),
		},
	},
	{
		Name: "Replace",
		Func: []any{
			new(func(string, string, string, int) string),
		},
	},
	{
		Name: "ReplaceAll",
		Func: []any{
			new(func(string, string, string) string),
		},
	},
	{
		Name: "Trim",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimLeft",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimRight",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimSpace",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "TrimPrefix",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimSuffix",
		Func: []any{
			new(func(string, string) string),
		},
	},
	{
		Name: "Get",
		Func: []any{
			new(func([]string, int) string),
		},
	},
	{
		Name: "ToString",
		Func: []any{
			new(func(any) string),
		},
	},
	{
		Name: "Match",
		Func: []any{
			new(func(string, string) bool),
		},
	},
	{
		Name: "KeyExists",
		Func: []any{
			new(func(string, map[string]any) bool),
		},
	},
	{
		Name: "LogInfo",
		Func: []any{
			new(func(string, ...any) bool),
		},
	},
	{
		Name: "B64Decode",
		Func: []any{
			new(func(string) string),
		},
	},
	{
		Name: "UnmarshalJSON",
		Func: []any{
			new(func(string, map[string]any, string) error),
		},
	},
	{
		Name: "ParseKV",
		Func: []any{
			new(func(string, map[string]any, string) error),
		},
	},
	{
		Name: "Hostname",
		Func: []any{
			new(func() (string, error)),
		},
	},
}
