package auth

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/model"
)

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

type AuthUser struct {
	user model.IUser
}

func (au *AuthUser) Scan(c *gin.Context) (isAbort bool) {
	if au.user != nil {
		return false
	}

	user := newUser().(model.IUser)
	userId, exist := middleware.AuthClaimID(c)
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotLoginError{}.Error()})
		return true
	}
	if err := user.Scan(userId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotExistError{}.Error()})
		return true
	}

	au.user = user

	return false
}

func (au *AuthUser) User() model.IUser {
	return au.user
}

func (au *AuthUser) UserId(c *gin.Context) (userId uint, isAbort bool) {
	exist := false
	userId, exist = middleware.AuthClaimID(c)
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotLoginError{}.Error()})
		return 0, true
	}
	return userId, false
}
