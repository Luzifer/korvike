package functions

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
)

func init() {
	registerFunction("hash", func(name string, v ...string) (string, error) {
		if len(v) < 1 {
			return "", errors.New("no string to hash")
		}

		var hash hash.Hash
		switch name {
		case "md5":
			hash = md5.New()
		case "sha1":
			hash = sha1.New()
		case "sha256":
			hash = sha256.New()
		case "sha512":
			hash = sha512.New()

		default:
			return "", fmt.Errorf("hash algo %q not supported", name)
		}

		if _, err := hash.Write([]byte(v[0])); err != nil {
			return "", fmt.Errorf("writing to hash: %w", err)
		}

		return fmt.Sprintf("%x", hash.Sum(nil)), nil
	})
}
