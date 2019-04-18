package migration

import (
	"sort"

	"github.com/jinzhu/gorm"
)

type Migrator interface {
	Up(db *gorm.DB) *gorm.DB
	Down(db *gorm.DB) *gorm.DB
	MigratorIdentifier
}

// contains all the migrators
var migratorList []Migrator

func AddMigrator(migrator Migrator) {
	// add migrator
	migratorList = append(migratorList, migrator)

	// do sort by timestamp
	sort.Slice(migratorList, func(i, j int) bool {
		return migratorList[i].Timestamp(&migratorList[i]) < migratorList[j].Timestamp(&migratorList[j])
	})
}

func newMigrator(name string) Migrator {
	for _, migrator := range migratorList {
		if name == migrator.Name(&migrator) {
			return migrator
		}
	}
	return nil
}
