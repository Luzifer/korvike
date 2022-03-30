package functions

import "strings"

func init() {
	registerFunction("join", strings.Join)
	registerFunction("split", strings.Split)
}
