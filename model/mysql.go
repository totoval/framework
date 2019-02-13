package model

import "totoval-framework/config"

type mysql struct {
	conn string
}

func NewMysql(connection string) *mysql {
	_db := new(mysql)
	_db.setConnection(connection)
	return _db
}

func (_mys *mysql) setConnection(connection string) {
	_mys.conn = connection
}

func (_mys *mysql) connection() string {
	return _mys.conn
}
func (_mys *mysql) username() string {
	return _mys.config("username")
}
func (_mys *mysql) password() string {
	return _mys.config("password")
}
func (_mys *mysql) host() string {
	return _mys.config("host")
}
func (_mys *mysql) port() string {
	return _mys.config("port")
}
func (_mys *mysql) database() string {
	return _mys.config("database")
}
func (_mys *mysql) charset() string {
	return _mys.config("charset")
}
func (_mys *mysql) Prefix() string {
	return _mys.config("prefix")
}
func (_mys *mysql) Driver() string {
	return _mys.config("driver")
}
func (_mys *mysql) collation() string {
	return _mys.config("collation")
}
func (_mys *mysql) config(key string) string {
	if value, ok := config.Get("database.connections."+ _mys.connection() + "." + key).(string); ok {
		return value
	}
	panic("database "+key+" parse error")
}
func (_mys *mysql) ConnectionArgs() string {
	return _mys.username()+":"+_mys.password()+"@"+"tcp("+_mys.host()+":"+_mys.port()+")/"+_mys.database()+"?charset="+_mys.charset()+"&parseTime=True&loc=Local"
}
