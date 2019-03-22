package database

import (
	"github.com/jinzhu/gorm"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/database/driver"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"

var db *gorm.DB
var dber databaser

func Initialize() {
	dber, db = setConnection("default")
	configOrm(dber)
}

func setConnection(conn string) (dber databaser, _db *gorm.DB) {
	// get database connection name
	_conn := conn
	if conn == "default" {
		//var _ok bool
		//_conn, _ok = config.Get("database." + conn).(string)
		//if !_ok {
		//	panic("database connection parse error")
		//}
		_conn = config.GetString("database." + conn)
		if _conn == "" {
			panic("database connection parse error")
		}
	}

	// get orm instance
	switch _conn {
	case "mysql":
		dber = driver.NewMysql(_conn)
		break
	default:
		panic("incorrect database connection provided")
	}

	// connect database
	_db, err := gorm.Open(_conn, dber.ConnectionArgs())
	if err != nil {
		panic("failed to connect database")
	}

	err = _db.DB().Ping()
	if err != nil {
		panic("failed to connect database by ping")
	}

	// debug mode
	if config.GetBool("app.debug") {
		_db = _db.Debug().LogMode(true)
	}

	_db.DB().SetMaxIdleConns(config.GetInt("database.max_idle_connections"))
	_db.DB().SetMaxOpenConns(config.GetInt("database.max_open_connections"))

	//defer _db.Close()
	return dber, _db
}

func Connection(conn string) (_db *gorm.DB) {
	_, _db = setConnection(conn)
	return _db
}

func DB() *gorm.DB {
	return db
}

func Prefix() string {
	return dber.Prefix()
}

func configOrm(orm databaser) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return Prefix() + defaultTableName
	}
}
