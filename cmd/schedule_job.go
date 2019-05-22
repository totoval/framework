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
			var __err error
			if _err, ok := err.(error); ok {
				__err = _err
			} else {
				__err = errors.New(command.When() + "schedule panic") //@todo err.(string) may be down when `panic(123)`
			}
			_ = log.Error(__err)
		}
	}(j.scmd)

	if err := j.scmd.handler()(j.scmd.args()); err != nil {
		_ = log.Error(err)
	}
}
