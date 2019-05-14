package cmd

import (
	"regexp"
	"strings"

	"github.com/urfave/cli"

	"github.com/totoval/framework/helpers/debug"
)

var commandMap map[string]parsedCommand

func init() {
	commandMap = make(map[string]parsedCommand)
}

func Add(c Commander) {
	cmdName, argList := parseCommand(c.Command())
	commandMap[cmdName] = parsedCommand{
		c,
		cmdName,
		argList,
	}
}

func parseCommand(commandWithArgs string) (name string, argList []string) {

	tmp := strings.SplitN(commandWithArgs, " ", 2)
	name = tmp[0]

	const reg = `(?m).*?(\s\{(.*?)\})`
	re := regexp.MustCompile(reg)
	for _, match := range re.FindAllStringSubmatch(commandWithArgs, -1) {
		if len(match) > 2 {
			argList = append(argList, match[2])
		}
	}

	return
}

func List() (cmdList []cli.Command) {
	for _, v := range commandMap {
		debug.Dump(v.Name(), v.commandCategory())
		cmdList = append(cmdList, cli.Command{
			Category: v.commandCategory(),
			Name:     v.Name(),
			//Aliases:  c.Aliases(),
			Usage: v.Description(),
			Action: func(_c *cli.Context) error {
				c := commandMap[_c.Command.Name]

				if !_c.Args().Present() {
					return c.Handler(newArg(nil))
				}

				argMap := c.mapArg(_c.Args())
				return c.Handler(newArg(&argMap))

			},
			ArgsUsage: v.argUsage(),
		})
	}
	return
}
