package auth

type UserNotLoginError struct{}

func (e UserNotLoginError) Error() string {
	return "user not login"
}

type UserNotExistError struct{}

func (e UserNotExistError) Error() string {
	return "user not exist"
}

type RequestUser struct {
	c    Context
	user IUser
}

func (au *RequestUser) SetContext(c Context) {
	au.c = c
}
func (au *RequestUser) ScanUser() error {

	// get cached user
	if au.user != nil {
		return nil
	}

	user := au.c.IUserModel()
	userId, err := au.UserId()
	if err != nil {
		return err
	}
	if err := user.Scan(userId); err != nil {
		return UserNotExistError{}
	}

	// set cache
	au.user = user

	return nil
}

func (au *RequestUser) User() IUser {
	return au.user
}

func (au *RequestUser) UserId() (userId uint, err error) {
	exist := false
	userId, exist = au.c.AuthClaimID()
	if !exist {
		return 0, UserNotLoginError{}
	}
	return userId, nil
}
