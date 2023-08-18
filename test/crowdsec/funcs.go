package crowdsec

import (
	"time"
)

var CustomFunctions = []struct {
	Name string
	Func []interface{}
}{
	{
		Name: "Distance",
		Func: []interface{}{
			new(func(string, string, string, string) (float64, error)),
		},
	},
	{
		Name: "GetFromStash",
		Func: []interface{}{
			new(func(string, string) (string, error)),
		},
	},
	{
		Name: "Atof",
		Func: []interface{}{
			new(func(string) float64),
		},
	},
	{
		Name: "JsonExtract",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "JsonExtractUnescape",
		Func: []interface{}{
			new(func(string, ...string) string),
		},
	},
	{
		Name: "JsonExtractLib",
		Func: []interface{}{
			new(func(string, ...string) string),
		},
	},
	{
		Name: "JsonExtractSlice",
		Func: []interface{}{
			new(func(string, string) []interface{}),
		},
	},
	{
		Name: "JsonExtractObject",
		Func: []interface{}{
			new(func(string, string) map[string]interface{}),
		},
	},
	{
		Name: "ToJsonString",
		Func: []interface{}{
			new(func(interface{}) string),
		},
	},
	{
		Name: "File",
		Func: []interface{}{
			new(func(string) []string),
		},
	},
	{
		Name: "RegexpInFile",
		Func: []interface{}{
			new(func(string, string) bool),
		},
	},
	{
		Name: "Upper",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "Lower",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "IpInRange",
		Func: []interface{}{
			new(func(string, string) bool),
		},
	},
	{
		Name: "TimeNow",
		Func: []interface{}{
			new(func() string),
		},
	},
	{
		Name: "ParseUri",
		Func: []interface{}{
			new(func(string) map[string][]string),
		},
	},
	{
		Name: "PathUnescape",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "QueryUnescape",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "PathEscape",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "QueryEscape",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "XMLGetAttributeValue",
		Func: []interface{}{
			new(func(string, string, string) string),
		},
	},
	{
		Name: "XMLGetNodeValue",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "IpToRange",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "IsIPV6",
		Func: []interface{}{
			new(func(string) bool),
		},
	},
	{
		Name: "IsIPV4",
		Func: []interface{}{
			new(func(string) bool),
		},
	},
	{
		Name: "IsIP",
		Func: []interface{}{
			new(func(string) bool),
		},
	},
	{
		Name: "LookupHost",
		Func: []interface{}{
			new(func(string) []string),
		},
	},
	{
		Name: "GetDecisionsCount",
		Func: []interface{}{
			new(func(string) int),
		},
	},
	{
		Name: "GetDecisionsSinceCount",
		Func: []interface{}{
			new(func(string, string) int),
		},
	},
	{
		Name: "Sprintf",
		Func: []interface{}{
			new(func(string, ...interface{}) string),
		},
	},
	{
		Name: "ParseUnix",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "SetInStash",
		Func: []interface{}{
			new(func(string, string, string, *time.Duration) error),
		},
	},
	{
		Name: "Fields",
		Func: []interface{}{
			new(func(string) []string),
		},
	},
	{
		Name: "Index",
		Func: []interface{}{
			new(func(string, string) int),
		},
	},
	{
		Name: "IndexAny",
		Func: []interface{}{
			new(func(string, string) int),
		},
	},
	{
		Name: "Join",
		Func: []interface{}{
			new(func([]string, string) string),
		},
	},
	{
		Name: "Split",
		Func: []interface{}{
			new(func(string, string) []string),
		},
	},
	{
		Name: "SplitAfter",
		Func: []interface{}{
			new(func(string, string) []string),
		},
	},
	{
		Name: "SplitAfterN",
		Func: []interface{}{
			new(func(string, string, int) []string),
		},
	},
	{
		Name: "SplitN",
		Func: []interface{}{
			new(func(string, string, int) []string),
		},
	},
	{
		Name: "Replace",
		Func: []interface{}{
			new(func(string, string, string, int) string),
		},
	},
	{
		Name: "ReplaceAll",
		Func: []interface{}{
			new(func(string, string, string) string),
		},
	},
	{
		Name: "Trim",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimLeft",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimRight",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimSpace",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "TrimPrefix",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "TrimSuffix",
		Func: []interface{}{
			new(func(string, string) string),
		},
	},
	{
		Name: "Get",
		Func: []interface{}{
			new(func([]string, int) string),
		},
	},
	{
		Name: "ToString",
		Func: []interface{}{
			new(func(interface{}) string),
		},
	},
	{
		Name: "Match",
		Func: []interface{}{
			new(func(string, string) bool),
		},
	},
	{
		Name: "KeyExists",
		Func: []interface{}{
			new(func(string, map[string]interface{}) bool),
		},
	},
	{
		Name: "LogInfo",
		Func: []interface{}{
			new(func(string, ...interface{}) bool),
		},
	},
	{
		Name: "B64Decode",
		Func: []interface{}{
			new(func(string) string),
		},
	},
	{
		Name: "UnmarshalJSON",
		Func: []interface{}{
			new(func(string, map[string]interface{}, string) error),
		},
	},
	{
		Name: "ParseKV",
		Func: []interface{}{
			new(func(string, map[string]interface{}, string) error),
		},
	},
	{
		Name: "Hostname",
		Func: []interface{}{
			new(func() (string, error)),
		},
	},
}
