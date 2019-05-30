package driver

import (
	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/helpers/zone"
)

type ProtoCacheGetter interface {
	Pget(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error
}
type ProtoCacher interface {
	Ppull(key string, valuePtr proto.Message, defaultValuePtr ...proto.Message) error
	Pput(key string, valuePtr proto.Message, future zone.Time) bool
	Padd(key string, valuePtr proto.Message, future zone.Time) bool
	Pforever(key string, valuePtr proto.Message) bool

	ProtoCacheGetter
}
