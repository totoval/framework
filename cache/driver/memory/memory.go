package memory

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	c "github.com/patrickmn/go-cache"

	"github.com/totoval/framework/cache/driver/common"
	. "github.com/totoval/framework/cache/utils"
	"github.com/totoval/framework/helpers/zone"
)

func NewMemory(prefix string, defaultExpirationMinute uint, cleanUpIntervalMinute uint) *memory {
	return &memory{
		memoryBasic: memoryBasic{
			cache:  c.New(zone.Duration(defaultExpirationMinute)*zone.Minute, zone.Duration(cleanUpIntervalMinute)*zone.Minute),
			prefix: prefix,
		},
	}
}

type memory struct {
	memoryBasic
	common.Common
}

func (m *memory) Pget(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := NewKey(key, m.Prefix())

	valueInterface, found := m.cache.Get(k.Prefixed())
	if !found {
		//@todo Event CacheMissed
		if len(defaultValuePtr) > 0 {
			return copier.Copy(valuePtr, defaultValuePtr[0])
		}
		return errors.New("key not exist")
	}

	//@todo Event CacheHit
	valueBytes, ok := valueInterface.([]byte)
	if !ok {
		return errors.New("key's value is not a valid proto buffer")
	}
	if err := proto.Unmarshal(valueBytes, valuePtr); err != nil {
		return err
	}
	return nil
}
