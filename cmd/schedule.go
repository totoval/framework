package cmd

type Schedule struct {
	commandList []*scheduleCommand
}

var scheduleHub *Schedule

func ScheduleCommandList() *[]*scheduleCommand {
	return &scheduleHub.commandList
}

func NewSchedule() *Schedule {
	var cmdList []*scheduleCommand
	scheduleHub = &Schedule{commandList: cmdList}
	return scheduleHub
}

func (s *Schedule) Command(commandWithArgData string) *scheduleCommand {
	name, argDataList := parseArgData(commandWithArgData)

	c, err := getParsedCommand(name)
	if err != nil {
		panic(err)
	}

	if len(argDataList) <= 0 {
		_cmd := newScheduleCommand(c.Handler, newArg(nil))
		s.commandList = append(s.commandList, _cmd)
		return _cmd
	}

	argMap := c.mapArg(argDataList)
	_cmd := newScheduleCommand(c.Handler, newArg(&argMap))
	s.commandList = append(s.commandList, _cmd)
	return _cmd
}
