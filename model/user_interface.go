package model

type UserIF interface {
	Scan(userId uint) error
}
