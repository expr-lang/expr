package environment

import "os"

const (
	ExprEnvGotag = "expr__gotag"
	DefaultGotag = "expr"
)

func GetGoTag() string {
	if v := os.Getenv(ExprEnvGotag); v != "" {
		return v
	} else {
		return DefaultGotag
	}
}
