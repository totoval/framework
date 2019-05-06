package job

import (
	"time"

	"github.com/golang/protobuf/proto"
)

type Job struct {
	param proto.Message
	delay time.Duration
}

func (j *Job) Name() string {
	panic("need implements")
}

func (j *Job) Handle(paramPtr proto.Message) error {
	panic("need implements")
}

func (j *Job) SetParam(paramPtr proto.Message) {
	j.param = paramPtr
}

func (j *Job) paramData() proto.Message {
	return j.param
}

func (j *Job) ParamProto() proto.Message {
	panic("need implements")
}

// default retry 3 times
func (j *Job) Retries() uint32 {
	return 3
}

// default no delay
func (j *Job) SetDelay(delay time.Duration) {
	j.delay = delay
}

// default no delay
func (j *Job) Delay() time.Duration {
	return j.delay
}
