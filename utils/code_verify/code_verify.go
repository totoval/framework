package code_verify

import (
	"github.com/totoval/framework/helpers/cache"

	"github.com/totoval/framework/helpers/str"
	"fmt"
	"time"
)

const VALIDATION_CACHE_KEY = "TOTOVAL_VALIDATION_CACHE_KEY_%s"

func validationCacheKey(index string) string {
	return fmt.Sprintf(VALIDATION_CACHE_KEY, index)
}

func Generate(index string, codeLen uint, expiredMinute uint) string {
	code := str.RandNumberString(codeLen)
	cache.Put(validationCacheKey(index), code, time.Now().Add(time.Duration(expiredMinute) * time.Minute))
	return code
}

func Verify(index string, code string) bool {
	cacheCode := cache.Pull(validationCacheKey(index))
	if cacheCode == nil {
		return false
	}

	if cacheCode.(string) == code {
		return true
	}

	return false
}