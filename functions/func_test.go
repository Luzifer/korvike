package functions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
)

func renderHelper(tpl string, ctx map[string]interface{}) string {
	t, err := template.New("mytemplate").Funcs(GetFunctionMap()).Parse(tpl)
	if err != nil {
		panic(err)
	}

	SetSubTemplateVariables(ctx)

	buf := bytes.NewBufferString("")
	if err := t.Execute(buf, ctx); err != nil {
		panic(err)
	}

	return buf.String()
}

func randomString() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func Test_GetFunctionMap(t *testing.T) {
	f := GetFunctionMap()
	if f == nil || len(f) < 1 {
		t.Fatal("No functions were registered.")
	}
}

func Test_env(t *testing.T) {
	result := randomString()
	os.Setenv("KORVIKE_TESTING", result)

	if r := renderHelper(`{{env "KORVIKE_TESTING"}}`, nil); r != result {
		t.Errorf("[env] did not receive expected string: %q (expected %q)", r, result)
	}
}

func Test_file(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	p := f.Name()
	result := randomString()
	fmt.Fprint(f, result)
	f.Close()
	defer os.Remove(p)

	if r := renderHelper(fmt.Sprintf("{{file %q}}", p), nil); r != result {
		t.Errorf("[file] did not receive expected string: %q (expected %q)", r, result)
	}
}

func Test_now(t *testing.T) {
	if _, err := time.Parse(time.RFC3339Nano, renderHelper(fmt.Sprintf("{{now %q}}", time.RFC3339Nano), nil)); err != nil {
		t.Errorf("[now] did not produce expected time format")
	}
}

func Test_tplexec(t *testing.T) {
	result := randomString()
	os.Setenv("KORVIKE_TESTING", result)

	if res := renderHelper(`{{ tplexec "{{ .var }}:{{ env \"KORVIKE_TESTING\" }}" }}`, map[string]interface{}{"var": "test"}); res != strings.Join([]string{
		"test",
		result,
	}, ":") {
		t.Errorf("[template] did not produce expected result %q != test", res)
	}
}
