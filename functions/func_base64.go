package functions

import "encoding/base64"

func init() {
	registerFunction("b64decode", func(name string, v ...string) (string, error) {
		b, err := base64.StdEncoding.DecodeString(name)
		return string(b), err
	})
	registerFunction("b64encode", func(name string, v ...string) string {
		return base64.StdEncoding.EncodeToString([]byte(name))
	})
}
