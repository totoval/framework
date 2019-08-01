package auth

type RequestIUser interface {
	ScanUser() error
	User() IUser
	UserId() (userId uint, err error)
	ScanUserWithJSON() (isAbort bool)
}
