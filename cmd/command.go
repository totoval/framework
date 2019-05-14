package cmd

import (
	"strings"

	"github.com/urfave/cli"
)

var commandMap map[string]Commander

func init() {
	commandMap = make(map[string]Commander)
}

func Add(c Commander) {
	commandMap[c.Command()] = c
}

type Command struct {
}

func commandCategory(command string) string {
	if !strings.Contains(command, ":") {
		return ""
	}
	tmp := strings.Split(command, ":")
	if len(tmp) > 1 {
		return tmp[0]
	}
	return ""
}

func List() (cmdList []cli.Command) {
	for _, c := range commandMap {
		cmdList = append(cmdList, cli.Command{
			Category: commandCategory(c.Command()),
			Name:     c.Command(),
			Aliases:  c.Aliases(),
			Usage:    c.Description(),
			Action:   c.Handler,
		})
	}
	return
}
