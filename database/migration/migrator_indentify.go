package migration

import (
	"reflect"
	"strings"
)

type MigratorIdentifier interface {
	Name(*Migrator) string
}

type MigratorIdentify struct{}

func (m *MigratorIdentify) Name(migrator *Migrator) string {
	tmp := strings.Split(reflect.TypeOf(*migrator).String(), ".")
	return tmp[len(tmp)-1]
}
