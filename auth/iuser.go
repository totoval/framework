package auth

type IUser interface {
	Scan(userId uint) error
	Value() interface{}
}
