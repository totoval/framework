package job

import (
	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/helpers/zone"
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

	SetDelay(delay zone.Duration)
	Delay() zone.Duration
}
