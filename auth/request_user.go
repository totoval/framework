package auth

import (
	"net/http"

	"github.com/totoval/framework/helpers/toto"
)

const CONTEXT_REQUEST_USER_KEY = "TOTOVAL_CONTEXT_REQUEST_USER"

type UserNotLoginError struct{}

func (e UserNotLoginError) Error() string {
	return "user not login"
}

type UserNotExistError struct{}

func (e UserNotExistError) Error() string {
	return "user not exist"
}

type RequestUser struct {
}

func (au *RequestUser) Scan(c Context) (isAbort bool) {

	// get cached user
	if au.User(c) != nil {
		return false
	}

	user := c.IUserModel()
	userId, exist := c.AuthClaimID()
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": UserNotLoginError{}.Error()})
		return true
	}
	if err := user.Scan(userId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": UserNotExistError{}.Error()})
		return true
	}

	// set cache
	c.Set(CONTEXT_REQUEST_USER_KEY, user)

	return false
}

func (au *RequestUser) User(c Context) IUser {
	if _requestUser, exists := c.Get(CONTEXT_REQUEST_USER_KEY); exists {
		if requestUser, ok := _requestUser.(IUser); ok {
			return requestUser
		}
	}
	return nil
}

func (au *RequestUser) UserId(c Context) (userId uint, isAbort bool) {
	exist := false
	userId, exist = c.AuthClaimID()
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": UserNotLoginError{}.Error()})
		return 0, true
	}
	return userId, false
}
