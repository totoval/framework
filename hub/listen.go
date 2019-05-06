package hub

import (
	"time"

	"github.com/golang/protobuf/proto"
)

type Listen struct {
}

func (l *Listen) Name() ListenerName {
	panic("need implements")
}

func (l *Listen) Subscribe() (eventPtrList [] Eventer) {
	panic("need implements")
}

func (l *Listen) Construct(paramPtr proto.Message) error {
	panic("need implements")
}

func (l *Listen) Handle() error {
	panic("need implements")
}

func (l *Listen) Retries() uint32 {
	return 0
}

func (l *Listen) Delay() time.Duration {
	return 0
}
