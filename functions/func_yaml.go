package functions

import (
	"bytes"
	"fmt"

	"go.yaml.in/yaml/v3"
)

const yamlIndent = 2 // Default is 4, like crazy peeps

var (
	docFromYAML = FunctionDocumentation{
		Title:       `{{ fromYaml <source> }}`,
		Description: "Parse YAML into a template value.",
		Template:    `{{ index (fromYaml "answer: 42") "answer" }}`,
	}

	docToYAML = FunctionDocumentation{
		Title:       `{{ toYaml <value> }}`,
		Description: "Serialize a value as YAML with two-space indentation.",
		Template:    `{{ toYaml (dict "answer" 42) }}`,
	}
)

func init() {
	registerFunction("fromYaml", tplFromYAML, true, docFromYAML)
	registerFunction("toYaml", tplToYAML, true, docToYAML)
}

func tplFromYAML(raw string) (out any, err error) {
	if err = yaml.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	return out, nil
}

func tplToYAML(in any) (out string, err error) {
	var (
		buf = new(bytes.Buffer)
		enc = yaml.NewEncoder(buf)
	)

	enc.SetIndent(yamlIndent)

	if err = enc.Encode(in); err != nil {
		return "", fmt.Errorf("marshalling data: %w", err)
	}

	return buf.String(), nil
}
