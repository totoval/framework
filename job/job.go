package job

import (
	"errors"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/queue"
)

var jobMap map[string]jobber

func init() {
	queue.Queue()
	jobMap = make(map[string]jobber)
}

func Add(j jobber) {
	j.SetParam(j.ParamProto())
	jobMap[j.Name()] = j
}

func Dispatch(j jobber) error {
	if err := queue.NewProducer("job", j.Name(), j.ParamData(), j.Retries(), j.Delay()).Push(); err != nil {
		return err
	}
	return nil
}

func Process(jobName string) {
	j := jobMap[jobName]
	if j == nil {
		panic(errors.New("job " + jobName + " doesn't exist"))
	}
	err := queue.NewConsumer("job", j.Name(), j.ParamProto(), j.Handle).Pop()
	if err != nil {
		panic(err)
	}
}

type jobber interface {
	Name() string

	SetParam(paramPtr proto.Message)
	ParamData() proto.Message
	ParamProto() proto.Message

	Handle(paramPtr proto.Message) error

	Retries() uint32

	SetDelay(delay time.Duration)
	Delay() time.Duration
}

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

func (j *Job) ParamData() proto.Message {
	return j.param
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
