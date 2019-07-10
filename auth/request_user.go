package auth

import (
	"net/http"
	"reflect"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/toto"
)

const CONTEXT_REQUEST_USER_KEY = "TOTOVAL_CONTEXT_REQUEST_USER"

func newUser() interface{} {
	typeof := reflect.TypeOf(config.GetInterface("auth.model_ptr"))
	ptr := reflect.New(typeof).Elem()
	val := reflect.New(typeof.Elem())
	ptr.Set(val)
	return ptr.Interface()
}

type UserNotLoginError struct{}

func (e UserNotLoginError) Error() string {
	return "user not login"
}

type UserNotExistError struct{}

func (e UserNotExistError) Error() string {
	return "user not exist"
}

type RequestUser struct {
	user IUser
}

func (au *RequestUser) Scan(c Context) (isAbort bool) {
	// if already scanned
	if au.user != nil {
		return false
	}

	// get cached user
	if _requestUser, exists := c.Get(CONTEXT_REQUEST_USER_KEY); exists {
		if requestUser, ok := _requestUser.(IUser); ok {
			au.user = requestUser
			return false
		}
	}

	user := newUser().(IUser)
	userId, exist := c.AuthClaimID()
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": UserNotLoginError{}.Error()})
		return true
	}
	if err := user.Scan(userId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": UserNotExistError{}.Error()})
		return true
	}

	au.user = user

	// set cache
	c.Set(CONTEXT_REQUEST_USER_KEY, user)

	return false
}

func (au *RequestUser) User() IUser {
	return au.user
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
