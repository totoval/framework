package middleware

import (
	"crypto"
	"fmt"
	"net/http"

	"github.com/totoval/framework/cache"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/request"

	"github.com/totoval/framework/helpers/bytes"
)

var limiter *cache.RateLimit

func Throttle(maxAttempts uint, decayMinutes uint) request.HandlerFunc {
	return func(c request.Context) {
		if limiter == nil {
			limiter = cache.NewRateLimit(cache.Cache())
		}

		key := requestSignature(c)

		if limiter.TooManyAttempts(key, int64(maxAttempts)) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, toto.V{"error": "Too Many Attempts"})
			return
		}

		limiter.Hit(key, int(decayMinutes))

		setHeader(c, maxAttempts, uint(calculateRemainingAttempts(key, maxAttempts, 0)), 0)

		c.Next()
	}
}

func calculateRemainingAttempts(key string, maxAttempts uint, retryAfter zone.Duration) int64 {
	if retryAfter == 0 {
		return limiter.RetriesLeft(key, int64(maxAttempts))
	}
	return 0
}

func setHeader(c request.Context, maxAttempts uint, remainingAttempts uint, retryAfter zone.Duration) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxAttempts))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingAttempts))

	if retryAfter != 0 {
		c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", zone.Now().Add(retryAfter).Unix()))
	}
}

func requestSignature(c request.Context) string {
	userId, exist := c.AuthClaimID()

	sha1 := crypto.SHA1.New()
	if exist {
		// has user
		sha1.Write([]byte(bytes.FromUint64(uint64(userId))))
		return string(sha1.Sum(nil))
	}

	// do not has user
	sha1.Write([]byte(c.Request().Host + "|" + c.ClientIP()))
	return string(sha1.Sum(nil))
}
