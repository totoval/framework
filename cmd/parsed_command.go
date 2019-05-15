package cmd

import (
	"fmt"
	"strings"
)

type parsedCommand struct {
	Commander
	name    commandName
	argList []string
}

func (pc *parsedCommand) Name() commandName {
	return pc.name
}
func (pc *parsedCommand) ArgList() []string {
	return pc.argList
}
func (pc *parsedCommand) argUsage() (argUsage string) {
	argUsage = ""
	for _, arg := range pc.ArgList() {
		argUsage += fmt.Sprintf(" <%s>", arg)
	}

	return strings.TrimSpace(argUsage)
}
func (pc *parsedCommand) mapArg(argDataList []string) (argData map[string]string) {
	argData = make(map[string]string)
	for i, argKey := range pc.ArgList() {
		if len(argDataList) < i+1 {
			break
		}
		argData[argKey] = argDataList[i]
	}

	return
}
func (pc *parsedCommand) commandCategory() string {
	if !strings.Contains(pc.name, ":") {
		return ""
	}
	tmp := strings.Split(pc.name, ":")
	if len(tmp) > 1 {
		return tmp[0]
	}
	return ""
}
