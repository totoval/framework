package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/model"
)

type UserNotPermitError struct{}

func (e UserNotPermitError) Error() string {
	return "user has no permission"
}

func Middleware(policy Policier, action Action) gin.HandlerFunc {

	return func(c *gin.Context) {
		user := auth.NewUser().(model.IUser)
		if middleware.AuthUser(c, user) {
			c.Abort()
			return
		}

		sub := user   // the user that wants to access a resource.
		obj := policy // the resource that is going to be accessed.
		act := action // the operation that the user performs on the resource.
		if !enfc.Enforce(sub, obj, act) {
			c.JSON(http.StatusForbidden, gin.H{"error": UserNotPermitError{}.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
