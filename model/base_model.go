package model

import (
	"github.com/jinzhu/gorm"
	"github.com/totoval/framework/database"
)

type BaseModeller interface {
	DB() *gorm.DB
	SetTX(db *gorm.DB)
}

type BaseModel struct {
	db *gorm.DB
}

func (bm *BaseModel) DB() *gorm.DB {
	if bm.db == nil {
		return database.DB()
	}
	return bm.db
}
func (bm *BaseModel) SetTX(db *gorm.DB) {
	bm.SetDB(db)
}
func (bm *BaseModel) SetDB(db *gorm.DB) {
	bm.db = db
}
