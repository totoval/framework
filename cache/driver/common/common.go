package common

import (
	"github.com/totoval/framework/cache/driver"
	"github.com/totoval/framework/helpers/zone"

	"github.com/golang/protobuf/proto"

	. "github.com/totoval/framework/cache/utils"
)

type Common struct {
	driver.BasicCacher
	driver.ProtoCacheGetter
}

func (c *Common) Ppull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error {
	k := NewKey(key, c.Prefix())

	err := c.Pget(k.Raw(), valuePtr, defaultValuePtr...)
	if err != nil {
		return err
	}

	c.Forget(k.Raw())

	return nil
}
func (c *Common) Pput(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return c.Put(key, valueBytes, future)
}
func (c *Common) Padd(key string, valuePtr proto.Message, future zone.Time) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return c.Add(key, valueBytes, future)
}
func (c *Common) Pforever(key string, valuePtr proto.Message) bool {
	valueBytes, err := proto.Marshal(valuePtr)
	if err != nil {
		return false
	}
	return c.Forever(key, valueBytes)
}
