package model

type IUser interface {
	Scan(userId uint) error
	Value() interface{}
}
