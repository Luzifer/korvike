package functions

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderHelper(t *testing.T, rawTpl string, ctx map[string]interface{}) string {
	tpl, err := template.New("mytemplate").Funcs(GetFunctionMap()).Parse(rawTpl)
	require.NoError(t, err)

	SetSubTemplateVariables(ctx)

	buf := bytes.NewBufferString("")
	require.NoError(t, tpl.Execute(buf, ctx))

	return buf.String()
}

func randomString() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	//#nosec:G404 // Okay to generate non-deterministic strings
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func Test_GetFunctionMap(t *testing.T) {
	f := GetFunctionMap()
	require.NotNil(t, f)
}

func Test_base64(t *testing.T) {
	var (
		b64   = "aGVsbG8="
		plain = "hello"
	)

	assert.Equal(t, plain, renderHelper(t, fmt.Sprintf(`{{ b64dec "%s" }}`, b64), nil))
	assert.Equal(t, b64, renderHelper(t, fmt.Sprintf(`{{ b64enc "%s" }}`, plain), nil))
}

func Test_env(t *testing.T) {
	result := randomString()
	require.NoError(t, os.Setenv("KORVIKE_TESTING", result))

	assert.Equal(t, result, renderHelper(t, `{{ env "KORVIKE_TESTING" }}`, nil))
}

func Test_file(t *testing.T) {
	f, err := os.CreateTemp("", "")
	require.NoError(t, err)

	p := f.Name()
	result := randomString()
	fmt.Fprint(f, result)
	require.NoError(t, f.Close())

	t.Cleanup(func() {
		require.NoError(t, os.Remove(p))
	})

	assert.Equal(t, result, renderHelper(t, fmt.Sprintf("{{ file %q }}", p), nil))
}

func Test_hash(t *testing.T) {
	input := "I'm a string to hash"
	for algo, exp := range map[string]string{
		"sha1sum":   "bd41599338445f401b8d3751fbe718e8a0b52004",
		"sha256sum": "ba32f090baf28862816a10da05509b31393704184ae49c68f9eb2933afa9e4d1",
		"sha512sum": "3288edcff4f28526fe8ecfb6d5182f2a446ab0572550c9591d1f5cacd377397af25c0274c9c2428c35422d215ccdc304d0353c093c76e750f9d7c4d54e64eed8",
	} {
		assert.Equal(t, exp, renderHelper(t, fmt.Sprintf("{{ %s %q }}", algo, input), nil), algo)
	}
}

func Test_now(t *testing.T) {
	_, err := time.Parse(
		time.RFC3339Nano,
		renderHelper(t, fmt.Sprintf("{{ now | date %q }}", time.RFC3339Nano), nil),
	)
	assert.NoError(t, err)
}

func Test_tplexec(t *testing.T) {
	result := randomString()
	require.NoError(t, os.Setenv("KORVIKE_TESTING", result))

	assert.Equal(
		t,
		strings.Join([]string{
			"test",
			result,
		}, ":"),
		renderHelper(t, `{{ tplexec "{{ .var }}:{{ env \"KORVIKE_TESTING\" }}" }}`, map[string]any{
			"var": "test",
		}),
	)
}
