package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type H = map[string]interface{}

func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
