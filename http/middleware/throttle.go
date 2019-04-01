package middleware

import (
	"crypto"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/cache"

	"github.com/totoval/framework/helpers/bytes"
)

var limiter *cache.RateLimit

func Throttle(maxAttempts uint, decayMinutes uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			limiter = cache.NewRateLimit(cache.Cache())
		}

		key := requestSignature(c)

		if limiter.TooManyAttempts(key, int64(maxAttempts)) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Attempts"})
			return
		}

		limiter.Hit(key, int(decayMinutes))

		setHeader(c, maxAttempts, uint(calculateRemainingAttempts(key, maxAttempts, 0)), 0)

		c.Next()
	}
}

func calculateRemainingAttempts(key string, maxAttempts uint, retryAfter time.Duration) int64 {
	if retryAfter == 0 {
		return limiter.RetriesLeft(key, int64(maxAttempts))
	}
	return 0
}

func setHeader(c *gin.Context, maxAttempts uint, remainingAttempts uint, retryAfter time.Duration) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxAttempts))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingAttempts))

	if retryAfter != 0 {
		c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(retryAfter).Unix()))
	}
}

func requestSignature(c *gin.Context) string {
	userId, exist := authClaimID(c)

	sha1 := crypto.SHA1.New()
	if exist {
		// has user
		sha1.Write([]byte(bytes.FromUint64(uint64(userId))))
		return string(sha1.Sum(nil))
	}

	// do not has user
	sha1.Write([]byte(c.Request.Host + "|" + c.ClientIP()))
	return string(sha1.Sum(nil))
}
