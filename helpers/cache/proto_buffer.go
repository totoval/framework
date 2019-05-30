package cache

import (
	"github.com/golang/protobuf/proto"

	c "github.com/totoval/framework/cache"
	"github.com/totoval/framework/helpers/zone"
)

func Pget(key string, valuePtr proto.Message, defaultValue ...proto.Message) error {
	return c.Cache().Pget(key, valuePtr, defaultValue...)
}
func Ppull(key string, valuePtr proto.Message, defaultValue ...proto.Message) error {
	return c.Cache().Ppull(key, valuePtr, defaultValue...)
}
func Pput(key string, value proto.Message, future zone.Time) bool {
	return c.Cache().Pput(key, value, future)
}
func Padd(key string, value proto.Message, future zone.Time) bool {
	return c.Cache().Padd(key, value, future)
}
func Pforever(key string, value proto.Message) bool {
	return c.Cache().Pforever(key, value)
}
