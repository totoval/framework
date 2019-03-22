package driver

import (
    "time"

    c "github.com/patrickmn/go-cache"
)

func NewMemory(prefix string, defaultExpirationMinute uint, cleanUpIntervalMinute uint) *memory {
    return &memory{
        cache: c.New(time.Duration(defaultExpirationMinute)*time.Minute, time.Duration(cleanUpIntervalMinute)*time.Minute),
        prefix: prefix,
    }
}

type memory struct {
    cache *c.Cache
    prefix string
}

func durationFromNow(future time.Time) time.Duration {
    return future.Sub(time.Now())
}
func (m *memory)prefixedKey(k string) string {
    return m.Prefix() + k
}

func (m *memory)Prefix() string {
    return m.prefix
}
func (m *memory)Has(key string) bool {
    _, found := m.cache.Get(key)
    return found
}
func (m *memory)Get(key string, defaultValue ...interface{}) interface{}{
    val, found := m.cache.Get(m.prefixedKey(key))
    if !found {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return nil
    }
    return val
}
func (m *memory)Pull(key string, defaultValue ...interface{}) interface{}{
    result := m.Get(key, defaultValue)
    m.Forget(key)
    return result
}
func (m *memory)Put(key string, value interface{}, future time.Time){
    m.cache.Set(m.prefixedKey(key), value, durationFromNow(future))
}
func (m *memory)Add(key string, value interface{}, future time.Time) bool{
    if err := m.cache.Add(m.prefixedKey(key), value, durationFromNow(future)); err != nil {
        return false
    }
    return true
}
func (m *memory)Increment(key string, value int) bool {
    if err := m.cache.Increment(m.prefixedKey(key), int64(value)); err != nil {
        return false
    }
    return true
}
func (m *memory)Decrement(key string, value int) bool{
    if err := m.cache.Decrement(m.prefixedKey(key), int64(value)); err != nil {
        return false
    }
    return true
}
func (m *memory)Forever(key string, value interface{}){
    m.cache.Set(m.prefixedKey(key), value, -1)
}
func (m *memory)Forget(key string) bool{
    m.cache.Delete(m.prefixedKey(key))
    return true
}