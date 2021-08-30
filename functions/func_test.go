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

func Test_base64(t *testing.T) {
	var (
		b64   = "aGVsbG8="
		plain = "hello"
	)

	if r := renderHelper(fmt.Sprintf(`{{b64decode "%s"}}`, b64), nil); r != plain {
		t.Errorf("[b64decode] did not yield expected string: %q (expected %q)", r, plain)
	}

	if r := renderHelper(fmt.Sprintf(`{{b64encode "%s"}}`, plain), nil); r != b64 {
		t.Errorf("[b64encode] did not yield expected string: %q (expected %q)", r, b64)
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

func Test_hash(t *testing.T) {
	input := "I'm a string to hash"
	for algo, exp := range map[string]string{
		"md5":    "d5adddaa0fd9f924b85e7874dc85f814",
		"sha1":   "bd41599338445f401b8d3751fbe718e8a0b52004",
		"sha256": "ba32f090baf28862816a10da05509b31393704184ae49c68f9eb2933afa9e4d1",
		"sha512": "3288edcff4f28526fe8ecfb6d5182f2a446ab0572550c9591d1f5cacd377397af25c0274c9c2428c35422d215ccdc304d0353c093c76e750f9d7c4d54e64eed8",
	} {
		if res := renderHelper(fmt.Sprintf("{{ hash %q %q }}", algo, input), nil); res != exp {
			t.Errorf("Hash algo %q yield unexpected result: exp=%q res=%q", algo, exp, res)
		}
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
