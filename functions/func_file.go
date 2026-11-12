package functions

import (
	"fmt"
	"os"
)

var (
	docFile = FunctionDocumentation{
		Title:       `{{ file <file name> }}`,
		Description: "Read a file and place it inside the template. Returns an empty string on error.",
		Template:    `{{ file "hello" }}`,
		PreCommands: []string{
			`echo "Hello World" > hello`,
		},
	}

	docMustFile = FunctionDocumentation{
		Title:       `{{ mustFile <file name> }}`,
		Description: "Read a file and place it inside the template. Returns an error when the file cannot be read.",
		Template:    `{{ mustFile "hello" }}`,
		PreCommands: []string{
			`echo "Hello World" > hello`,
		},
	}
)

func init() {
	registerFunction("file", tplFile, false, docFile)
	registerFunction("mustFile", tplMustFile, false, docMustFile)
}

func tplFile(name string) string {
	fc, err := tplMustFile(name)
	if err != nil {
		return ""
	}
	return fc
}

func tplMustFile(name string) (string, error) {
	rawValue, err := os.ReadFile(name) //#nosec:G304 // Intended to load custom file
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(rawValue), nil
}
