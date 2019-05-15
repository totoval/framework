package cmd

import "log"

type job struct {
	scmd *scheduleCommand
}

func NewJob(scmd *scheduleCommand) *job {
	return &job{scmd: scmd}
}

func (j *job) Run() {
	if err := j.scmd.handler()(j.scmd.args()); err != nil {
		log.Println("......schedule run failed......", err)
	}
}
