package hub

import (
	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/helpers/zone"
)

type Listener interface {
	Name() ListenerName
	Subscribe() (eventPtrList []Eventer)

	Construct(paramPtr proto.Message) error
	Handle() error

	Retries() uint32
	Delay() zone.Duration
}

type ListenerName = string
