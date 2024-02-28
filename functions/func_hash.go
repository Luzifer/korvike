package functions

import (
	"crypto/sha512"
	"encoding/hex"
)

func init() {
	registerFunction("sha512sum", func(data string) string {
		hash := sha512.Sum512([]byte(data))
		return hex.EncodeToString(hash[:])
	})
}
