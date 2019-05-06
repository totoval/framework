package job

import (
	"time"

	"github.com/golang/protobuf/proto"
)

var jobMap map[string]jobber

func init() {
	jobMap = make(map[string]jobber)
}

func Add(j jobber) {
	j.SetParam(j.ParamProto()) // for init
	jobMap[j.Name()] = j
}

type jobber interface {
	Name() string

	SetParam(paramPtr proto.Message)
	paramData() proto.Message
	ParamProto() proto.Message

	Handle(paramPtr proto.Message) error

	Retries() uint32

	SetDelay(delay time.Duration)
	Delay() time.Duration
}
