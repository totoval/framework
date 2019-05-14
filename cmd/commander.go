package cmd

import "github.com/urfave/cli"

type Commander interface {
	Command() string
	Aliases() []string
	Description() string
	Handler(c *cli.Context) error
}
