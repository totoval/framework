package cmd

import (
	"errors"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

type commandName = string

var commandMap map[commandName]parsedCommand

func init() {
	commandMap = make(map[commandName]parsedCommand)
}

func Add(c Commander) {
	cmdName, argList := parseCommand(c.Command())
	commandMap[cmdName] = parsedCommand{
		c,
		cmdName,
		argList,
	}
}

func getParsedCommand(name commandName) (*parsedCommand, error) {
	pcmd, ok := commandMap[name]
	if !ok {
		return nil, errors.New("command not found")
	}
	return &pcmd, nil
}

func parseCommand(commandWithArgs string) (name commandName, argList []string) {

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

func parseArgData(commandWithArgData string) (name commandName, argData []string) {
	tmp := strings.Split(commandWithArgData, " ")
	if len(tmp) > 0 {
		name = tmp[0]
	}
	argData = tmp[1:]
	return
}

func List() (cmdList []cli.Command) {
	for _, v := range commandMap {
		cmdList = append(cmdList, cli.Command{
			Category: v.commandCategory(),
			Name:     v.Name(),
			//Aliases:  c.Aliases(),
			Usage: v.Description(),
			Action: func(_c *cli.Context) error {
				c, err := getParsedCommand(_c.Command.Name)
				if err != nil {
					return err
				}

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
