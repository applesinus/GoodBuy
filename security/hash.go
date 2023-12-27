package security

import (
	"crypto/md5"
	"encoding/hex"
)

func CheckHash(original, hashed string) bool {
	hasher := md5.New()
	hasher.Write([]byte(original))
	return hex.EncodeToString(hasher.Sum(nil)) == hashed
}

func Hash(original string) string {
	hasher := md5.New()
	hasher.Write([]byte(original))
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}
