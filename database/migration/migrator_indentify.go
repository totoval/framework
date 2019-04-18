package migration

import (
	"reflect"
	"strconv"
	"strings"
)

const TIMESTAMP_LENGTH = 10

type MigratorIdentifier interface {
	Name(migrator *Migrator) string
	Timestamp(migrator *Migrator) uint
}

type MigratorIdentify struct{}

func (m *MigratorIdentify) Name(migrator *Migrator) string {
	tmp := strings.Split(reflect.TypeOf(*migrator).String(), ".")
	return tmp[len(tmp)-1]
}

func (m *MigratorIdentify) Timestamp(migrator *Migrator) uint {
	name := m.Name(migrator)
	t, err := strconv.ParseUint(name[len(name)-TIMESTAMP_LENGTH:], 10, 32)
	if err != nil {
		panic(err)
	}
	//@todo check whether `t` is a timestamp
	return uint(t)
}
