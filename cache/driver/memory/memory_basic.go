package memory

import (
	"encoding"

	c "github.com/patrickmn/go-cache"

	. "github.com/totoval/framework/cache/utils"
	"github.com/totoval/framework/helpers/zone"
)

type memoryBasic struct {
	cache  *c.Cache
	prefix string
}

func (m *memoryBasic) Prefix() string {
	return m.prefix
}
func (m *memoryBasic) Has(key string) bool {
	k := NewKey(key, m.Prefix())

	_, found := m.cache.Get(k.Prefixed())
	return found
}
func (m *memoryBasic) Get(key string, defaultValue ...interface{}) interface{} {
	k := NewKey(key, m.Prefix())

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
func (m *memoryBasic) Pull(key string, defaultValue ...interface{}) interface{} {
	k := NewKey(key, m.Prefix())

	val := m.Get(k.Raw(), defaultValue...)
	if val == nil {
		return nil
	}

	m.Forget(k.Raw())
	return val
}
func (m *memoryBasic) parseValue(v interface{}) interface{} {
	switch v := v.(type) {
	//case nil:
	//case string:
	//case []byte:
	//case int:
	//case int8:
	//case int16:
	//case int32:
	//case int64:
	//case uint:
	//case uint8:
	//case uint16:
	//case uint32:
	//case uint64:
	//case float32:
	//case float64:
	//case bool:
	//	return v
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil {
			return v
		}
		return b
	default:
		return v
	}
}

func (m *memoryBasic) Put(key string, value interface{}, future zone.Time) bool {
	k := NewKey(key, m.Prefix())

	m.cache.Set(k.Prefixed(), m.parseValue(value), DurationFromNow(future))

	//@todo Event KeyWritten
	return true
}
func (m *memoryBasic) Add(key string, value interface{}, future zone.Time) bool {
	k := NewKey(key, m.Prefix())

	// if exist or expired return false
	if err := m.cache.Add(k.Prefixed(), m.parseValue(value), DurationFromNow(future)); err != nil {
		return false
	}

	//@todo Event KeyWritten
	return true
}
func (m *memoryBasic) Increment(key string, value int64) (incremented int64, success bool) {
	k := NewKey(key, m.Prefix())

	incremented, err := m.cache.IncrementInt64(k.Prefixed(), m.parseValue(value).(int64))
	if err != nil {
		return 0, false
	}
	return incremented, true
}
func (m *memoryBasic) Decrement(key string, value int64) (decremented int64, success bool) {
	k := NewKey(key, m.Prefix())

	decremented, err := m.cache.DecrementInt64(k.Prefixed(), m.parseValue(value).(int64))
	if err != nil {
		return 0, false
	}
	return decremented, true
}
func (m *memoryBasic) Forever(key string, value interface{}) bool {
	k := NewKey(key, m.Prefix())

	m.cache.Set(k.Prefixed(), m.parseValue(value), -1)

	//@todo Event KeyWritten
	return true
}
func (m *memoryBasic) Forget(key string) bool {
	k := NewKey(key, m.Prefix())

	m.cache.Delete(k.Prefixed())

	//@todo Event KeyForget
	return true
}
func (m *memoryBasic) Close() error {
	m.cache.Flush()
	return nil
}
