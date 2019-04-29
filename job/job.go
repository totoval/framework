package job

import (
	"errors"

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
	if err := queue.NewProducer("job", j.Name(), j.ParamData(), j.Retries()).Push(); err != nil {
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
	Handle(paramPtr proto.Message) error
	ParamProto() proto.Message
	Retries() uint32
}

type Job struct {
	param proto.Message
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

func (j *Job) Retries() uint32 {
	return 0
}
