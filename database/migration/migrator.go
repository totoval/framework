package migration

import (
	"github.com/jinzhu/gorm"
)

type Migrator interface {
	Up(db *gorm.DB) (*gorm.DB)
	Down(db *gorm.DB) (*gorm.DB)
	MigratorIdentifier
}

// contains all the migrators
var migratorList []Migrator

func AddMigrator(migrator Migrator) {
	migratorList = append(migratorList, migrator)
}

func newMigrator(name string) (Migrator) {
	for _, migrator := range migratorList {
		if name == migrator.Name(&migrator) {
			return migrator
		}
	}
	return nil
}
