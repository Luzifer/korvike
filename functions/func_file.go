package functions

import (
	"fmt"
	"os"
)

func init() {
	registerFunction("file", func(name string) string {
		fc, err := tplReadFile(name)
		if err != nil {
			return ""
		}
		return fc
	})

	registerFunction("mustFile", tplReadFile)
}

func tplReadFile(name string) (string, error) {
	rawValue, err := os.ReadFile(name) //#nosec:G304 // Intended to load custom file
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(rawValue), nil
}
