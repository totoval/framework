package cmd

import (
	"fmt"

	"github.com/totoval/framework/cmd/schedule"
)

type scheduleCommand struct {
	commandHandler func(arg *Arg) error
	arg            *Arg
	when           schedule.When
}

func newScheduleCommand(handler func(arg *Arg) error, arg *Arg) *scheduleCommand {
	return &scheduleCommand{commandHandler: handler, arg: arg}
}

func (sc *scheduleCommand) Yearly() {
	sc.when = schedule.EveryYear
}
func (sc *scheduleCommand) Monthly() {
	sc.when = schedule.EveryMonth
}
func (sc *scheduleCommand) Daily() {
	sc.when = schedule.EveryDay
}
func (sc *scheduleCommand) Hourly() {
	sc.when = schedule.EveryHour
}
func (sc *scheduleCommand) EveryMinute() {
	sc.when = schedule.EveryMinute
}
func (sc *scheduleCommand) EverySecond() {
	sc.when = schedule.EverySecond
}
func (sc *scheduleCommand) EveryDays(d uint) {
	sc.when = fmt.Sprintf(schedule.EveryHoursFormat, d*24)
}
func (sc *scheduleCommand) EveryHours(h uint) {
	sc.when = fmt.Sprintf(schedule.EveryHoursFormat, h)
}
func (sc *scheduleCommand) EveryMinutes(m uint) {
	sc.when = fmt.Sprintf(schedule.EveryMinutesFormat, m)
}
func (sc *scheduleCommand) EverySeconds(s uint) {
	sc.when = fmt.Sprintf(schedule.EverySecondsFormat, s)
}

func (sc *scheduleCommand) When() schedule.When {
	return sc.when
}
func (sc *scheduleCommand) args() *Arg {
	return sc.arg
}
func (sc *scheduleCommand) handler() func(arg *Arg) error {
	return sc.commandHandler
}
