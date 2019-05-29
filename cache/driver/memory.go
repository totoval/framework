package driver

import (
	c "github.com/patrickmn/go-cache"

	"github.com/totoval/framework/helpers/zone"
)

func NewMemory(prefix string, defaultExpirationMinute uint, cleanUpIntervalMinute uint) *memory {
	return &memory{
		cache:  c.New(zone.Duration(defaultExpirationMinute)*zone.Minute, zone.Duration(cleanUpIntervalMinute)*zone.Minute),
		prefix: prefix,
	}
}

type memory struct {
	cache  *c.Cache
	prefix string
}

func (m *memory) Prefix() string {
	return m.prefix
}
func (m *memory) Has(key string) bool {
	k := newKey(key, m.Prefix())

	_, found := m.cache.Get(k.Prefixed())
	return found
}
func (m *memory) Get(key string, defaultValue ...interface{}) interface{} {
	k := newKey(key, m.Prefix())

	val, found := m.cache.Get(k.Prefixed())
	if !found {
		//@todo Event CacheMissed
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	//@todo Event CacheHit
	return val
}
func (m *memory) Pull(key string, defaultValue ...interface{}) interface{} {
	k := newKey(key, m.Prefix())

	val := m.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}

	m.Forget(k.Raw())
	return val
}
func (m *memory) Put(key string, value interface{}, future zone.Time) bool {
	k := newKey(key, m.Prefix())

	m.cache.Set(k.Prefixed(), value, durationFromNow(future))

	//@todo Event KeyWritten
	return true
}
func (m *memory) Add(key string, value interface{}, future zone.Time) bool {
	k := newKey(key, m.Prefix())

	// if exist or expired return false
	if err := m.cache.Add(k.Prefixed(), value, durationFromNow(future)); err != nil {
		return false
	}

	//@todo Event KeyWritten
	return true
}
func (m *memory) Increment(key string, value int64) (incremented int64, success bool) {
	k := newKey(key, m.Prefix())

	incremented, err := m.cache.IncrementInt64(k.Prefixed(), value)
	if err != nil {
		return 0, false
	}
	return incremented, true
}
func (m *memory) Decrement(key string, value int64) (decremented int64, success bool) {
	k := newKey(key, m.Prefix())

	decremented, err := m.cache.DecrementInt64(k.Prefixed(), value)
	if err != nil {
		return 0, false
	}
	return decremented, true
}
func (m *memory) Forever(key string, value interface{}) bool {
	k := newKey(key, m.Prefix())

	m.cache.Set(k.Prefixed(), value, -1)

	//@todo Event KeyWritten
	return true
}
func (m *memory) Forget(key string) bool {
	k := newKey(key, m.Prefix())

	m.cache.Delete(k.Prefixed())

	//@todo Event KeyForget
	return true
}
