package functions

import (
	"bytes"
	"fmt"
	"text/template"
)

var (
	docTplexec = FunctionDocumentation{
		Title:       `{{ tplexec (file "my.tpl") }}`,
		Description: "Execute the given template with the same function set and variables as the parent template.",
		Template:    `{{ tplexec (file "my.tpl") }}`,
		Environ: map[string]string{
			"FOO": "bar",
		},
		PreCommands: []string{
			`echo '{{ env "FOO" }}' >my.tpl`,
		},
	}
	subTemplateVariables map[string]any
)

// SetSubTemplateVariables sets variables to apply in `tplexec` function
func SetSubTemplateVariables(m map[string]any) { subTemplateVariables = m }

func init() {
	registerFunction("tplexec", tplTplexec, false, docTplexec)
}

func tplTplexec(rawTpl string) (string, error) {
	tpl, err := template.New("in").Funcs(GetFunctionMap()).Parse(rawTpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	out := new(bytes.Buffer)
	if err := tpl.Execute(out, subTemplateVariables); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return out.String(), nil
}
