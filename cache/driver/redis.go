package driver

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/totoval/framework/helpers/zone"
	"log"
	"reflect"
	"strconv"
)

func NewRedis(hostname string, port uint, db uint, auth string, prefix string, poolMaxIdle uint, poolMaxActive uint) (*redisCache, error) {
	if port > 65535 {
		return nil, errors.New("ports invalid")
	}
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", hostname, port))
			if err != nil {
				return nil, err
			}
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		MaxIdle:   int(poolMaxIdle),
		MaxActive: int(poolMaxActive),
	}
	// test connect
	testC := pool.Get()
	if err := testC.Flush(); err != nil {
		return nil, err
	}

	return &redisCache{
		cachePool: pool,
		prefix:    prefix,
	}, nil
}

type redisCache struct {
	cachePool *redis.Pool
	prefix    string
}

func (r *redisCache) parseVariable(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.FormatUint(uint64(val.(int)), 10)
	case uint:
		return strconv.FormatUint(uint64(val.(uint)), 10)
	case int64:
		return strconv.FormatUint(uint64(val.(int64)), 10)
	case uint64:
		return strconv.FormatUint(uint64(val.(uint64)), 10)
	default:
		panic(fmt.Sprintf("type:%v variable type invalid", reflect.TypeOf(val)))
	}
}

func (r *redisCache) prefixedKey(k string) string {
	return r.Prefix() + k
}

func (r *redisCache) resultToString(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case ([]uint8):
		return string(val.([]uint8))
	case int64:
		return strconv.FormatUint(uint64(val.(int64)), 10)
	case uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case int:
		return strconv.FormatUint(uint64(val.(int)), 10)
	default:
		panic("type is not invalid")
	}
}

func (r *redisCache) Prefix() string {
	return r.prefix
}
func (r *redisCache) Has(key string) bool {
	_, err := r.cachePool.Get().Do("GET", r.prefixedKey(key))
	if err != nil {
		return false
	}
	return true
}
func (r *redisCache) Get(key string, defaultValue ...interface{}) interface{} {
	val, err := r.cachePool.Get().Do("GET", r.prefixedKey(key))
	if err != nil || val == nil {
		// @todo Event CacheMissed
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	// @todo Event CacheHit
	return r.resultToString(val)
}
func (r *redisCache) Pull(key string, defaultValue ...interface{}) interface{} {
	result := r.Get(key, defaultValue...)
	r.Forget(key)
	return r.resultToString(result)
}
func (r *redisCache) Put(key string, value interface{}, future zone.Time) {
	conn := r.cachePool.Get()
	if _, err := conn.Do("MULTI"); err != nil {
		// @todo return the error
		return
	}
	conn.Do("SET", r.prefixedKey(key), value)
	conn.Do("EXPIREAT", r.prefixedKey(key), future.Unix())
	if _, err := conn.Do("EXEC"); err != nil {
		log.Fatal(err)
		conn.Do("DISCARD")
	}
	// @todo Event KeyWritten
}
func (r *redisCache) Add(key string, value interface{}, future zone.Time) bool {
	conn := r.cachePool.Get()
	if _, err := conn.Do("MULTI"); err != nil {
		// @todo return the error
		return false
	}
	conn.Do("SET", r.prefixedKey(key), value)
	conn.Do("EXPIREAT", r.prefixedKey(key), future.Unix())
	if _, err := conn.Do("EXEC"); err != nil {
		conn.Do("DISCARD")
		return false
	}

	// @todo Event KeyWritten
	return true
}
func (r *redisCache) Increment(key string, value int64) (incremented int64, success bool) {
	val, err := r.cachePool.Get().Do("INCRBY", r.prefixedKey(key), r.parseVariable(value))
	if err != nil {
		return 0, false
	}
	incremented, _ = strconv.ParseInt(r.resultToString(val), 10, 64)
	return incremented, true
}
func (r *redisCache) Decrement(key string, value int64) (decremented int64, success bool) {
	val, err := r.cachePool.Get().Do("DECRBY", r.prefixedKey(key), value)
	if err != nil {
		return 0, false
	}

	decremented, _ = strconv.ParseInt(r.resultToString(val), 10, 64)
	return decremented, true
}
func (r *redisCache) Forever(key string, value interface{}) {
	r.cachePool.Get().Do("EXPIRE", r.prefixedKey(key), value, -1)

	// @todo Event KeyWritten
}
func (r *redisCache) Forget(key string) bool {
	r.cachePool.Get().Do("DEL", r.prefixedKey(key))

	// @todo Event KeyForget
	return true
}

func (r *redisCache) Close() error {
	return r.cachePool.Close()
}
