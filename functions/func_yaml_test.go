package functions

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type brokenYAML struct{}

func TestFromYaml(t *testing.T) {
	rawYAML := `
array:
  - first
  - second
dict:
  nested: value
int: 42
string: hello
`

	assert.Equal(
		t,
		"hello|42|second|value",
		renderHelper(
			t,
			`{{ $yaml := fromYaml .raw }}{{ index $yaml "string" }}|{{ index $yaml "int" }}|{{ index (index $yaml "array") 1 }}|{{ index (index $yaml "dict") "nested" }}`,
			map[string]any{"raw": rawYAML},
		),
	)
}

func TestFromYamlInvalidInput(t *testing.T) {
	_, err := tplFromYAML("invalid:\n  - yaml\n broken\n")
	require.Error(t, err)
	assert.ErrorContains(t, err, "unmarshalling YAML")
}

func TestFromYamlLiteralArray(t *testing.T) {
	assert.Equal(
		t,
		"first|second",
		renderHelper(
			t,
			`{{ $yaml := fromYaml .raw }}{{ index $yaml 0 }}|{{ index $yaml 1 }}`,
			map[string]any{"raw": "- first\n- second\n"},
		),
	)
}

func TestFromYamlLiteralString(t *testing.T) {
	assert.Equal(
		t,
		"hello",
		renderHelper(
			t,
			`{{ fromYaml .raw }}`,
			map[string]any{"raw": "hello\n"},
		),
	)
}

func TestToYaml(t *testing.T) {
	rawYAML, err := tplToYAML(map[string]any{
		"array":  []string{"first", "second"},
		"dict":   map[string]string{"nested": "value"},
		"int":    42,
		"string": "hello",
	})
	require.NoError(t, err)

	assert.YAMLEq(t, `
array:
  - first
  - second
dict:
  nested: value
int: 42
string: hello
`, rawYAML)
}

func TestToYamlInvalidInput(t *testing.T) {
	_, err := tplToYAML(brokenYAML{})
	require.Error(t, err)
	assert.ErrorContains(t, err, "marshalling data")
}

func (brokenYAML) MarshalYAML() (any, error) {
	return nil, errors.New("broken yaml")
}
