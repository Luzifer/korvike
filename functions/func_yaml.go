package functions

import (
	"bytes"
	"fmt"

	"go.yaml.in/yaml/v3"
)

const yamlIndent = 2 // Default is 4, like crazy peeps

func init() {
	registerFunction("fromYaml", tplFromYAML, true)
	registerFunction("toYaml", tplToYAML, true)
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
