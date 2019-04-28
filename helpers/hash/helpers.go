package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
