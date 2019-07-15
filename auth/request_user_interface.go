package auth

type RequestIUser interface {
	Scan(c Context) (isAbort bool)
	User() IUser
	UserId(c Context) (userId uint, isAbort bool)
}
