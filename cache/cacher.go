package cache

import (
	"github.com/totoval/framework/helpers/zone"
)

type cacher interface {
	Prefix() string

	Has(key string) bool
	Get(key string, defaultValue ...interface{}) interface{}
	Pull(key string, defaultValue ...interface{}) interface{}
	Put(key string, value interface{}, future zone.Time)
	Add(key string, value interface{}, future zone.Time) bool
	Increment(key string, value int64) (incremented int64, success bool)
	Decrement(key string, value int64) (decremented int64, success bool)
	Forever(key string, value interface{})
	Forget(key string) bool
}
