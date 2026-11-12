package functions

import "net/url"

var docURLEscape = FunctionDocumentation{
	Title:       `{{ urlescape <input string> }}`,
	Description: "Escape input for use inside a URL query parameter.",
	Template:    `{{ urlescape "Hellö Wörld@Golang" }}`,
}

func init() {
	registerFunction("urlescape", url.QueryEscape, true, docURLEscape)
}
