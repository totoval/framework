package driver

import (
	"fmt"

	r "github.com/go-redis/redis"

	"github.com/totoval/framework/helpers/zone"
)

func NewRedis(host string, port string, password string, dbIndex int, prefix string) *redis {
	client := r.NewClient(&r.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbIndex,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return &redis{
		cache:  client,
		prefix: prefix,
	}
}

type redis struct {
	cache  *r.Client
	prefix string
}

func (re *redis) Prefix() string {
	return re.prefix
}
func (re *redis) Has(key string) bool {
	k := newKey(key, re.Prefix())

	exists, err := re.cache.Exists(k.Prefixed()).Result()
	if err != nil {
		return false
	}
	if exists <= 0 {
		return false
	}

	return true
}
func (re *redis) Get(key string, defaultValue ...interface{}) interface{} {
	k := newKey(key, re.Prefix())

	var val interface{}
	err := re.cache.Get(k.Prefixed()).Scan(&val)
	if err != nil {
		return err
	}
	if val == nil {
		//@todo Event CacheMissed
		return nil
	}

	//@todo Event CacheHit
	return val
}
func (re *redis) Pull(key string, defaultValue ...interface{}) interface{} {
	k := newKey(key, re.Prefix())

	val := re.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}

	re.Forget(k.Raw())

	return val
}
func (re *redis) Put(key string, value interface{}, future zone.Time) bool {
	k := newKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, durationFromNow(future)).Result()
	if err != nil {
		return false
	}

	return true

	//@todo Event KeyWritten
}
func (re *redis) Add(key string, value interface{}, future zone.Time) bool {
	k := newKey(key, re.Prefix())

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
func (re *redis) Increment(key string, value int64) (incremented int64, success bool) {
	k := newKey(key, re.Prefix())

	incremented, err := re.cache.IncrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return incremented, true
}
func (re *redis) Decrement(key string, value int64) (decremented int64, success bool) {
	k := newKey(key, re.Prefix())

	decremented, err := re.cache.DecrBy(k.Prefixed(), value).Result()
	if err != nil {
		return 0, false
	}

	return decremented, true
}
func (re *redis) Forever(key string, value interface{}) bool {
	k := newKey(key, re.Prefix())

	_, err := re.cache.Set(k.Prefixed(), value, 0).Result()
	if err != nil {
		return false
	}

	//@todo Event KeyWritten
	return true
}
func (re *redis) Forget(key string) bool {
	k := newKey(key, re.Prefix())

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
