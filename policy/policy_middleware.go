package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
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
		authUser := &auth.AuthUser{}
		if authUser.Scan(c) {
			c.Abort()
			return
		}

		// validate policy
		if !policyValidate(authUser.User(), policy, action, routeParamMap) {
			forbid(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func forbid(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": UserNotPermitError{}.Error()})
}
