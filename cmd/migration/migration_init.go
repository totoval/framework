package migration

import (
	"github.com/urfave/cli"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/database/migration"
)

func init() {
	cmd.Add(&MigrationInit{})
}

type MigrationInit struct {
}

func (mi *MigrationInit) Command() string {
	return "migration:init"
}

func (mi *MigrationInit) Aliases() []string {
	return []string{"c"}
}

func (mi *MigrationInit) Description() string {
	return "complete a task on the list"
}

func (mi *MigrationInit) Handler(c *cli.Context) error {
	m := &migration.MigrationUtils{}
	m.SetUp()
	return nil
}
