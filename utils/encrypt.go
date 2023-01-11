package utils

import (
	"fmt"
	"crypto/sha256"
)

func Entcrypt(s string) string {
	encrypted := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return encrypted
}