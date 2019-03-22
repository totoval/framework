package database

type databaser interface {
	ConnectionArgs() string
	Driver() string
	Prefix() string
}
