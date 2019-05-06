package hub

import "github.com/gogo/protobuf/proto"

type Event struct {
	param proto.Message
}

func (e *Event) SetParam(paramPtr proto.Message) {
	e.param = paramPtr
}

func (e *Event) paramData() proto.Message {
	return e.param
}

func (e *Event) ParamProto() proto.Message {
	panic("need implements")
}
