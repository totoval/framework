package hub

import (
	"time"

	"github.com/golang/protobuf/proto"
)

type Listener interface {
	Name() ListenerName
	Subscribe() (eventPtrList []Eventer)

	Construct(paramPtr proto.Message) error
	Handle() error

	Retries() uint32
	Delay() time.Duration
}

type ListenerName = string
