package cmd

import (
	"errors"

	"github.com/totoval/framework/helpers/log"
)

type job struct {
	scmd *scheduleCommand
}

func NewJob(scmd *scheduleCommand) *job {
	return &job{scmd: scmd}
}

func (j *job) Run() {
	defer func(command *scheduleCommand) {
		if err := recover(); err != nil {
			_err := errors.New(command.When() + "schedule panic")
			if __err, ok := err.(error); ok {
				_err = __err
			}
			_ = log.Error(_err)
		}
	}(j.scmd)

	if err := j.scmd.handler()(j.scmd.args()); err != nil {
		_ = log.Error(err)
	}
}
