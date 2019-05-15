package cmd

type Schedule struct {
	commandMap map[commandName]*scheduleCommand
}

var scheduleHub *Schedule

func ScheduleCommandMap() *map[commandName]*scheduleCommand {
	return &scheduleHub.commandMap
}

func NewSchedule() *Schedule {
	scheduleHub = &Schedule{commandMap: make(map[commandName]*scheduleCommand)}
	return scheduleHub
}

func (s *Schedule) Command(commandWithArgData string) *scheduleCommand {
	name, argDataList := parseArgData(commandWithArgData)

	c, err := getParsedCommand(name)
	if err != nil {
		panic(err)
	}

	if len(argDataList) <= 0 {
		s.commandMap[name] = newScheduleCommand(c.Handler, newArg(nil))
		return s.commandMap[name]
	}

	argMap := c.mapArg(argDataList)
	s.commandMap[name] = newScheduleCommand(c.Handler, newArg(&argMap))
	return s.commandMap[name]
}
