package schedule

import (
	"github.com/robfig/cron"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/helpers/zone"
)

func init() {
	cmd.Add(&Schedule{})
}

type Schedule struct {
}

func (s *Schedule) Command() string {
	return "schedule:run"
}

func (s *Schedule) Description() string {
	return "Schedule a failed queue from database"
}

func (s *Schedule) Handler(arg *cmd.Arg) error {
	c := cron.NewWithLocation(zone.GetLocation())
	for _, scheduleCommand := range *cmd.ScheduleCommandList() {
		if err := c.AddJob(scheduleCommand.When(), cmd.NewJob(scheduleCommand)); err != nil {
			panic(err)
		}
	}
	c.Start()
	defer c.Stop()

	select {}

	return nil
}
