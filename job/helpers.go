package job

import (
	"errors"

	"github.com/totoval/framework/queue"
)

func Dispatch(j jobber) error {
	if err := queue.NewProducer("job", j.Name(), j.paramData(), j.Retries(), j.Delay()).Push(); err != nil {
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
