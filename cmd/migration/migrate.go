package migration

import (
	"github.com/urfave/cli"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/database/migration"
)

func init() {
	cmd.Add(&Migrate{})
}

type Migrate struct {
}

func (mi *Migrate) Command() string {
	return "migrate"
}

func (mi *Migrate) Aliases() []string {
	return []string{"c"}
}

func (mi *Migrate) Description() string {
	return "complete a task on the list"
}

func (mi *Migrate) Handler(c *cli.Context) error {
	m := &migration.MigrationUtils{}
	m.Migrate()

	return nil
}
