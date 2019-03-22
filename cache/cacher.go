package cache

import "time"

type cacher interface {
    Prefix() string

    Has(key string) bool
    Get(key string, defaultValue ...interface{}) interface{}
    Pull(key string, defaultValue ...interface{}) interface{}
    Put(key string, value interface{}, future time.Time)
    Add(key string, value interface{}, future time.Time) bool
    Increment(key string, value int) bool
    Decrement(key string, value int) bool
    Forever(key string, value interface{})
    Forget(key string) bool
}
