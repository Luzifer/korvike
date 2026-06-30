package functions

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

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
	raw, err := yaml.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("marshalling data: %w", err)
	}

	return string(raw), nil
}
