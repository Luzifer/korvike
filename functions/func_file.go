package functions

import (
	"os"
)

func init() {
	registerFunction("file", func(name string, v ...string) string {
		defaultValue := ""
		if len(v) > 0 {
			defaultValue = v[0]
		}
		if _, err := os.Stat(name); err == nil {
			//#nosec:G304 // Intended to load custom file
			if rawValue, err := os.ReadFile(name); err == nil {
				return string(rawValue)
			}
		}
		return defaultValue
	})
}
