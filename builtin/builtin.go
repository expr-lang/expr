package builtin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/antonmedv/expr/vm/runtime"
)

type Function struct {
	Name     string
	Func     func(args ...interface{}) (interface{}, error)
	Types    []reflect.Type
	Builtin1 func(arg interface{}) interface{}
	Validate func(args []reflect.Type) (reflect.Type, error)
}

var (
	Index map[string]int
	Names []string
)

func init() {
	Index = make(map[string]int)
	Names = make([]string, len(Builtins))
	for i, fn := range Builtins {
		Index[fn.Name] = i
		Names[i] = fn.Name
	}
}

var Builtins = []*Function{
	{
		Name:     "len",
		Builtin1: Len,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String, reflect.Interface:
				return integerType, nil
			}
			return anyType, fmt.Errorf("invalid argument for len (type %s)", args[0])
		},
	},
	{
		Name:     "type",
		Builtin1: Type,
		Types:    types(new(func(interface{}) string)),
	},
	{
		Name:     "abs",
		Builtin1: Abs,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Interface:
				return args[0], nil
			}
			return anyType, fmt.Errorf("invalid argument for abs (type %s)", args[0])
		},
	},
	{
		Name:     "int",
		Builtin1: Int,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return integerType, nil
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return integerType, nil
			case reflect.String:
				return integerType, nil
			}
			return anyType, fmt.Errorf("invalid argument for int (type %s)", args[0])
		},
	},
	{
		Name:     "float",
		Builtin1: Float,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return floatType, nil
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return floatType, nil
			case reflect.String:
				return floatType, nil
			}
			return anyType, fmt.Errorf("invalid argument for float (type %s)", args[0])
		},
	},
	{
		Name:     "string",
		Builtin1: String,
		Types:    types(new(func(any interface{}) string)),
	},
	{
		Name: "trim",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) == 1 {
				return strings.TrimSpace(args[0].(string)), nil
			} else if len(args) == 2 {
				return strings.Trim(args[0].(string), args[1].(string)), nil
			} else {
				return nil, fmt.Errorf("invalid number of arguments for trim (expected 1 or 2, got %d)", len(args))
			}
		},
		Types: types(
			strings.TrimSpace,
			strings.Trim,
		),
	},
	{
		Name: "trimPrefix",
		Func: func(args ...interface{}) (interface{}, error) {
			s := " "
			if len(args) == 2 {
				s = args[1].(string)
			}
			return strings.TrimPrefix(args[0].(string), s), nil
		},
		Types: types(
			strings.TrimPrefix,
			new(func(string) string),
		),
	},
	{
		Name: "trimSuffix",
		Func: func(args ...interface{}) (interface{}, error) {
			s := " "
			if len(args) == 2 {
				s = args[1].(string)
			}
			return strings.TrimSuffix(args[0].(string), s), nil
		},
		Types: types(
			strings.TrimSuffix,
			new(func(string) string),
		),
	},
	{
		Name: "upper",
		Builtin1: func(arg interface{}) interface{} {
			return strings.ToUpper(arg.(string))
		},
		Types: types(strings.ToUpper),
	},
	{
		Name: "lower",
		Builtin1: func(arg interface{}) interface{} {
			return strings.ToLower(arg.(string))
		},
		Types: types(strings.ToLower),
	},
	{
		Name: "split",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) == 2 {
				return strings.Split(args[0].(string), args[1].(string)), nil
			} else if len(args) == 3 {
				return strings.SplitN(args[0].(string), args[1].(string), runtime.ToInt(args[2])), nil
			} else {
				return nil, fmt.Errorf("invalid number of arguments for split (expected 2 or 3, got %d)", len(args))
			}
		},
		Types: types(
			strings.Split,
			strings.SplitN,
		),
	},
	{
		Name: "splitAfter",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) == 2 {
				return strings.SplitAfter(args[0].(string), args[1].(string)), nil
			} else if len(args) == 3 {
				return strings.SplitAfterN(args[0].(string), args[1].(string), runtime.ToInt(args[2])), nil
			} else {
				return nil, fmt.Errorf("invalid number of arguments for splitAfter (expected 2 or 3, got %d)", len(args))
			}
		},
		Types: types(
			strings.SplitAfter,
			strings.SplitAfterN,
		),
	},
	{
		Name: "replace",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) == 4 {
				return strings.Replace(args[0].(string), args[1].(string), args[2].(string), runtime.ToInt(args[3])), nil
			} else if len(args) == 3 {
				return strings.ReplaceAll(args[0].(string), args[1].(string), args[2].(string)), nil
			} else {
				return nil, fmt.Errorf("invalid number of arguments for replace (expected 3 or 4, got %d)", len(args))
			}
		},
		Types: types(
			strings.Replace,
			strings.ReplaceAll,
		),
	},
	{
		Name: "repeat",
		Func: func(args ...interface{}) (interface{}, error) {
			n := runtime.ToInt(args[1])
			if n > 1e6 {
				panic("memory budget exceeded")
			}
			return strings.Repeat(args[0].(string), n), nil
		},
		Types: types(strings.Repeat),
	},
	{
		Name: "join",
		Func: func(args ...interface{}) (interface{}, error) {
			glue := ""
			if len(args) == 2 {
				glue = args[1].(string)
			}
			switch args[0].(type) {
			case []string:
				return strings.Join(args[0].([]string), glue), nil
			case []interface{}:
				var s []string
				for _, arg := range args[0].([]interface{}) {
					s = append(s, arg.(string))
				}
				return strings.Join(s, glue), nil
			}
			return nil, fmt.Errorf("invalid argument for join (type %s)", reflect.TypeOf(args[0]))
		},
		Types: types(
			strings.Join,
			new(func([]interface{}, string) string),
			new(func([]interface{}) string),
			new(func([]string, string) string),
			new(func([]string) string),
		),
	},
	{
		Name: "indexOf",
		Func: func(args ...interface{}) (interface{}, error) {
			return strings.Index(args[0].(string), args[1].(string)), nil
		},
		Types: types(strings.Index),
	},
	{
		Name: "lastIndexOf",
		Func: func(args ...interface{}) (interface{}, error) {
			return strings.LastIndex(args[0].(string), args[1].(string)), nil
		},
		Types: types(strings.LastIndex),
	},
	{
		Name: "hasPrefix",
		Func: func(args ...interface{}) (interface{}, error) {
			return strings.HasPrefix(args[0].(string), args[1].(string)), nil
		},
		Types: types(strings.HasPrefix),
	},
	{
		Name: "hasSuffix",
		Func: func(args ...interface{}) (interface{}, error) {
			return strings.HasSuffix(args[0].(string), args[1].(string)), nil
		},
		Types: types(strings.HasSuffix),
	},
	{
		Name: "max",
		Func: Max,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) == 0 {
				return anyType, fmt.Errorf("not enough arguments to call max")
			}
			for _, arg := range args {
				switch kind(arg) {
				case reflect.Interface:
					return anyType, nil
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				default:
					return anyType, fmt.Errorf("invalid argument for max (type %s)", arg)
				}
			}
			return args[0], nil
		},
	},
	{
		Name: "min",
		Func: Min,
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) == 0 {
				return anyType, fmt.Errorf("not enough arguments to call min")
			}
			for _, arg := range args {
				switch kind(arg) {
				case reflect.Interface:
					return anyType, nil
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				default:
					return anyType, fmt.Errorf("invalid argument for min (type %s)", arg)
				}
			}
			return args[0], nil
		},
	},
	{
		Name: "toJSON",
		Func: func(args ...interface{}) (interface{}, error) {
			b, err := json.MarshalIndent(args[0], "", "  ")
			if err != nil {
				return nil, err
			}
			return string(b), nil
		},
		Types: types(new(func(interface{}) string)),
	},
	{
		Name: "fromJSON",
		Func: func(args ...interface{}) (interface{}, error) {
			var v interface{}
			err := json.Unmarshal([]byte(args[0].(string)), &v)
			if err != nil {
				return nil, err
			}
			return v, nil
		},
		Types: types(new(func(string) interface{})),
	},
	{
		Name: "toBase64",
		Func: func(args ...interface{}) (interface{}, error) {
			return base64.StdEncoding.EncodeToString([]byte(args[0].(string))), nil
		},
		Types: types(new(func(string) string)),
	},
	{
		Name: "fromBase64",
		Func: func(args ...interface{}) (interface{}, error) {
			b, err := base64.StdEncoding.DecodeString(args[0].(string))
			if err != nil {
				return nil, err
			}
			return string(b), nil
		},
		Types: types(new(func(string) string)),
	},
	{
		Name: "now",
		Func: func(args ...interface{}) (interface{}, error) {
			return time.Now(), nil
		},
		Types: types(new(func() time.Time)),
	},
	{
		Name: "duration",
		Func: func(args ...interface{}) (interface{}, error) {
			return time.ParseDuration(args[0].(string))
		},
		Types: types(time.ParseDuration),
	},
	{
		Name: "date",
		Func: func(args ...interface{}) (interface{}, error) {
			date := args[0].(string)
			if len(args) == 2 {
				layout := args[1].(string)
				return time.Parse(layout, date)
			}
			if len(args) == 3 {
				layout := args[1].(string)
				timeZone := args[2].(string)
				tz, err := time.LoadLocation(timeZone)
				if err != nil {
					return nil, err
				}
				t, err := time.ParseInLocation(layout, date, tz)
				if err != nil {
					return nil, err
				}
				return t, nil
			}

			layouts := []string{
				"2006-01-02",
				"15:04:05",
				"2006-01-02 15:04:05",
				time.RFC3339,
				time.RFC822,
				time.RFC850,
				time.RFC1123,
			}
			for _, layout := range layouts {
				t, err := time.Parse(layout, date)
				if err == nil {
					return t, nil
				}
			}
			return nil, fmt.Errorf("invalid date %s", date)
		},
		Types: types(
			new(func(string) time.Time),
			new(func(string, string) time.Time),
			new(func(string, string, string) time.Time),
		),
	},
	{
		Name: "first",
		Func: func(args ...interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			return runtime.Fetch(args[0], 0), nil
		},
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return anyType, nil
			case reflect.Slice, reflect.Array:
				return args[0].Elem(), nil
			}
			return anyType, fmt.Errorf("cannot get first element from %s", args[0])
		},
	},
	{
		Name: "last",
		Func: func(args ...interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			return runtime.Fetch(args[0], -1), nil
		},
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return anyType, nil
			case reflect.Slice, reflect.Array:
				return args[0].Elem(), nil
			}
			return anyType, fmt.Errorf("cannot get last element from %s", args[0])
		},
	},
	{
		Name: "get",
		Func: func(args ...interface{}) (out interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			return runtime.Fetch(args[0], args[1]), nil
		},
	},
	{
		Name: "keys",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			v := reflect.ValueOf(args[0])
			if v.Kind() != reflect.Map {
				return nil, fmt.Errorf("cannot get keys from %s", v.Kind())
			}
			keys := v.MapKeys()
			out := make([]interface{}, len(keys))
			for i, key := range keys {
				out[i] = key.Interface()
			}
			return out, nil
		},
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return arrayType, nil
			case reflect.Map:
				return arrayType, nil
			}
			return anyType, fmt.Errorf("cannot get keys from %s", args[0])
		},
	},
	{
		Name: "values",
		Func: func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			v := reflect.ValueOf(args[0])
			if v.Kind() != reflect.Map {
				return nil, fmt.Errorf("cannot get values from %s", v.Kind())
			}
			keys := v.MapKeys()
			out := make([]interface{}, len(keys))
			for i, key := range keys {
				out[i] = v.MapIndex(key).Interface()
			}
			return out, nil
		},
		Validate: func(args []reflect.Type) (reflect.Type, error) {
			if len(args) != 1 {
				return anyType, fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
			}
			switch kind(args[0]) {
			case reflect.Interface:
				return arrayType, nil
			case reflect.Map:
				return arrayType, nil
			}
			return anyType, fmt.Errorf("cannot get values from %s", args[0])
		},
	},
}
