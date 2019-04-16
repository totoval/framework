package migration

import (
    "fmt"

    "github.com/jinzhu/gorm"

    "github.com/totoval/framework/database"

    "github.com/totoval/framework/config"
)

type Migration struct {
    ID        uint   `gorm:"column:id;primary_key;auto_increment;"`
    Migration string `gorm:"column:migration;type:varchar(255)"`
    Batch     uint   `gorm:"column:batch;"`
}

// 建立migration表
func (m *Migration) up(db *gorm.DB) {
    tx := db.Begin()
    {
        tx.CreateTable(&m)
    }
    tx.Commit()
}

func (m *Migration) TableName() string {
    tableName := config.GetString("database.migrations")
    if tableName == "" {
        panic("migrations table name parse error")
    }

    tableNameWithPrefix := fmt.Sprintf("%s%s", database.Prefix(), tableName)
    return tableNameWithPrefix
}

func (m *Migration) Name() string {
    return m.Migration
}
