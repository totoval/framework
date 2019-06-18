package cache

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/totoval/framework/helpers/zone"
)

type RateLimit struct {
	cache cacher
}

func NewRateLimit(cache cacher) *RateLimit {
	return &RateLimit{
		cache: cache,
	}
}

const RATE_LIMITE_CACHE_KEY = "TOTOVAL_RATE_LIMIT_%s"
const RATE_LIMITE_TIMER_CACHE_KEY = "TOTOVAL_RATE_LIMIT_%s_TIMER"

func rateLimitCacheKey(key string) string {
	return fmt.Sprintf(RATE_LIMITE_CACHE_KEY, key)
}
func rateLimitTimerCacheKey(key string) string {
	return fmt.Sprintf(RATE_LIMITE_TIMER_CACHE_KEY, key)
}

func (rl *RateLimit) TooManyAttempts(key string, maxAttempts int64) bool {
	if rl.Attempts(key) >= maxAttempts {
		if rl.cache.Has(rateLimitTimerCacheKey(key)) {
			return true
		}

		rl.ResetAttempts(key)
	}
	return false
}
func (rl *RateLimit) Hit(key string, decayMinutes int) int64 {
	expiredAt := zone.Now().Add(zone.Duration(decayMinutes) * zone.Minute)
	rl.cache.Add(rateLimitTimerCacheKey(key), expiredAt.Unix(), expiredAt)

	added := rl.cache.Add(rateLimitCacheKey(key), int64(0), expiredAt)
	hits, success := rl.cache.Increment(rateLimitCacheKey(key), 1)

	if !added && (success && hits == 1) {
		rl.cache.Put(rateLimitCacheKey(key), int64(1), expiredAt)
	}

	return hits
}
func (rl *RateLimit) Attempts(key string) int64 {
	attempts, err := parseInt64(rl.cache.Get(rateLimitCacheKey(key), int64(0)))
	if err != nil {
		return -1
	}

	return attempts
}
func (rl *RateLimit) ResetAttempts(key string) bool {
	return rl.cache.Forget(rateLimitCacheKey(key))
}
func (rl *RateLimit) RetriesLeft(key string, maxAttempts int64) int64 {
	attempts := rl.Attempts(key)
	return maxAttempts - attempts
}
func (rl *RateLimit) Clear(key string) {
	rl.ResetAttempts(key)
	rl.cache.Forget(rateLimitTimerCacheKey(key))
}
func (rl *RateLimit) AvailableIn(key string) zone.Duration {
	expiredAtUnix, err := parseInt64(rl.cache.Get(rateLimitTimerCacheKey(key)))
	if err != nil {
		expiredAtUnix = 0
	}
	expiredAt := zone.Unix(expiredAtUnix, 0)
	return expiredAt.Sub(zone.Now())
}
func parseInt64(val interface{}) (int64, error) {
	switch val.(type) {
	case int64:
		return val.(int64), nil
	case int32:
		return int64(val.(int32)), nil
	case int:
		return int64(val.(int)), nil
	case string:
		_int64, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return -1, err
		}
		return _int64, nil
	default:
		return -1, errors.New("unrecognized value")
	}
}
