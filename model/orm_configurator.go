package model

type OrmConfigurator interface {
	ConnectionArgs() string
	Driver() string
	Prefix() string
}
