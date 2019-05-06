package hub

import (
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

type EventerName = string

type Eventer interface {
	//Name() string

	SetParam(paramPtr proto.Message)
	paramData() proto.Message
	ParamProto() proto.Message
}

func EventName(e Eventer) string {
	tmp := strings.Split(reflect.TypeOf(e).String(), ".")
	return tmp[len(tmp)-1]
}
