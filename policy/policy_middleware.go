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
		// get route url param
		routeParamMap := make(map[string]string)
		for _, param := range c.Params {
			routeParamMap[param.Key] = param.Value
		}

		// get user
		user := auth.NewUser().(model.IUser)
		if middleware.AuthUser(c, user) {
			c.Abort()
			return
		}

		// validate policy
		if !policyValidate(user, policy, action, routeParamMap) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": UserNotPermitError{}.Error()})
			return
		}

		c.Next()
	}
}
