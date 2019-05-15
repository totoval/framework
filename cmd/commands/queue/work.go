package queue

import (
	"errors"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/job"
)

func init() {
	cmd.Add(&Work{})
}

type Work struct {
}

func (w *Work) Command() string {
	return "queue:work {job_name}"
}

func (w *Work) Description() string {
	return "Process job"
}

func (w *Work) Handler(arg *cmd.Arg) error {
	jobNamePtr, err := arg.Get("job_name")
	if err != nil {
		return err
	}

	if jobNamePtr == nil {
		return errors.New("job_name is invalid")
	}

	job.Process(*jobNamePtr)

	select {}

	return nil
}
