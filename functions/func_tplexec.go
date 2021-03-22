package functions

import (
	"bytes"
	"fmt"
	"text/template"
)

var subTemplateVariables map[string]interface{}

func SetSubTemplateVariables(m map[string]interface{}) { subTemplateVariables = m }

func init() {
	registerFunction("tplexec", func(rawTpl string) (string, error) {
		tpl, err := template.New("in").Funcs(GetFunctionMap()).Parse(rawTpl)
		if err != nil {
			return "", fmt.Errorf("parse template: %w", err)
		}

		out := new(bytes.Buffer)
		if err := tpl.Execute(out, subTemplateVariables); err != nil {
			return "", fmt.Errorf("execute template: %w", err)
		}

		return out.String(), nil
	})
}
