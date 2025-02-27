package saas

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(encodeString string) string {
	h := md5.New()
	h.Write([]byte(encodeString))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
