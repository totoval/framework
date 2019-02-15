package migration

import (
	"github.com/jinzhu/gorm"
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
	tableName, ok := config.Get("database.migrations").(string);
	if !ok {
		panic("migrations table name parse error")
	}
	return tableName
}

func (m *Migration) Name() string {
	return m.Migration
}

