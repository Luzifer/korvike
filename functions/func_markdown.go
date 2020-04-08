package functions

import "github.com/russross/blackfriday/v2"

func init() {
	registerFunction("markdown", func(name string, v ...string) string {
		return string(blackfriday.Run([]byte(name)))
	})
}
