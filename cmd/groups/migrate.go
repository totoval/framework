package groups

import (
	"github.com/urfave/cli"
	"framework/database/migration"
)

type MigrateCommand struct {
	MigratorInitializer func()
	ChLog chan interface{}
}

func (mc *MigrateCommand) MigrationInit() (command cli.Command) {
	return cli.Command{
		Name:    "migration:init",
		Aliases: []string{"c"},
		Usage:   "complete a task on the list",
		Action: func(c *cli.Context) error {
			m := &migration.MigrationUtils{}
			m.Init(mc.ChLog)
			go m.SetUp()
			return nil
		}}
}
func (mc *MigrateCommand) Migrate() (command cli.Command) {
	return cli.Command{
		Name:    "migrate",
		Aliases: []string{"c"},
		Usage:   "complete a task on the list",
		Action: func(c *cli.Context) error {
			mc.MigratorInitializer()
			m := &migration.MigrationUtils{}
			m.Init(mc.ChLog)
			go m.Migrate()

			return nil
		}}
}
func (mc *MigrateCommand) MigrateRollBack() (command cli.Command) {
	return cli.Command{
		Name:    "migration:rollback",
		Aliases: []string{"c"},
		Usage:   "complete a task on the list",
		Action: func(c *cli.Context) error {
			mc.MigratorInitializer()
			m := &migration.MigrationUtils{}
			m.Init(mc.ChLog)
			go m.Rollback()
			return nil
		}}
}