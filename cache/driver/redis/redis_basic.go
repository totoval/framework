package redis

import (
	r "github.com/go-redis/redis"

	. "github.com/totoval/framework/cache/utils"
	"github.com/totoval/framework/helpers/zone"
)

type redisBasic struct {
	cache  *r.Client
	prefix string
}

func (re *redisBasic) Cache() *r.Client {
	return re.cache
}

func (re *redisBasic) Prefix() string {
	return re.prefix
}
func (re *redisBasic) Has(key string) bool {
	k := NewKey(key, re.Prefix())

	exists, err := re.cache.Exists(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if exists <= 0 {
		return false
	}

	return true
}
func (re *redisBasic) Get(key string, defaultValue ...interface{}) interface{} {
	k := NewKey(key, re.Prefix())

	if !re.Has(k.Raw()) {
		//@todo Event CacheMissed
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	valStr, err := re.cache.Get(k.Prefixed()).Result()
	if err != nil {
		return err
	}

	//@todo Event CacheHit
	return valStr
}
func (re *redisBasic) Pull(key string, defaultValue ...interface{}) interface{} {
	k := NewKey(key, re.Prefix())

	val := re.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}

	re.Forget(k.Raw())

	return val
}
func (re *redisBasic) Put(key string, value interface{}, future zone.Time) bool {
	k := NewKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, DurationFromNow(future)).Result()
	if err != nil {
		return false
	}

	return true

	//@todo Event KeyWritten
}
func (re *redisBasic) Add(key string, value interface{}, future zone.Time) bool {
	k := NewKey(key, re.Prefix())

	// if expired return false
	ttl, err := re.cache.TTL(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if ttl > 0 {
		return false
	}

	// if exists return false
	if re.Has(k.Raw()) {
		return false
	}

	result := re.Put(k.Raw(), value, future)

	//@todo Event KeyWritten
	return result
}
func (re *redisBasic) Increment(key string, value int64) (incremented int64, success bool) {
	k := NewKey(key, re.Prefix())

	incremented, err := re.cache.IncrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return incremented, true
}
func (re *redisBasic) Decrement(key string, value int64) (decremented int64, success bool) {
	k := NewKey(key, re.Prefix())

	decremented, err := re.cache.DecrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return decremented, true
}
func (re *redisBasic) Forever(key string, value interface{}) bool {
	k := NewKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, 0).Result()
	if err != nil {
		return false
	}

	//@todo Event KeyWritten
	return true
}
func (re *redisBasic) Forget(key string) bool {
	k := NewKey(key, re.Prefix())

	result, err := re.cache.Del(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if result <= 0 {
		return false
	}

	//@todo Event KeyForget
	return true
}
func (re *redisBasic) Close() error {
	return re.cache.Close()
}
