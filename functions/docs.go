package functions

import "slices"

type (
	// FunctionDocumentation describes a function's documentation entry and example.
	FunctionDocumentation struct {
		// Title of the doc-entry, i.e. `{{ myFunc .var }}`
		Title string
		// What does the function do?
		Description string
		// Template to show (and execute) to demonstrate function
		Template string

		// Optional: Mark function as deprecated (to be removed)
		Deprecated bool
		// Optional: ENV-vars to set before running documentation
		Environ map[string]string
		// Optional: Override output in case function needs special access
		// to remote services which are not present in every environment
		Output string
		// Optional: Commands to prepare the environment
		PreCommands []string
		// Whether the function is included in `--safe` mode
		Safe bool
		// Optional: Vars to provide during doc execution
		Vars map[string]string
	}
)

var registeredDocumentation = []FunctionDocumentation{
	{
		Title:       `{{ .<variable name> }}`,
		Description: "Take key-value pairs from the CLI and replace them inside the template",
		Template:    `{{ .foo }}`,
		Safe:        true,
		Vars: map[string]string{
			"foo": "bar",
		},
	},
}

// GetFunctionDocumentation returns the registered function documentation.
func GetFunctionDocumentation() []FunctionDocumentation {
	return slices.Clone(registeredDocumentation)
}
