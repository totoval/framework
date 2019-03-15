package orm

type OrmConfigurator interface {
	ConnectionArgs() string
	Driver() string
	Prefix() string
}
