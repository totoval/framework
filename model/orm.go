package model

type Orm interface {
	ConnectionArgs() string
	Driver() string
	Prefix() string
}
