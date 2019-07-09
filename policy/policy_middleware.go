package policy

import (
	"net/http"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/request"
)

type UserNotPermitError struct{}

func (e UserNotPermitError) Error() string {
	return "user has no permission"
}

func Middleware(policy Policier, action Action) request.HandlerFunc {

	return func(c *request.Context) {
		// get route url param
		routeParamMap := make(map[string]string)
		for _, param := range c.Params {
			routeParamMap[param.Key] = param.Value
		}

		// get user
		requestUser := &auth.RequestUser{}
		if requestUser.Scan(c) {
			c.Abort()
			return
		}

		// validate policy
		if !policyValidate(requestUser.User(), policy, action, routeParamMap) {
			forbid(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func forbid(c *request.Context) {
	c.JSON(http.StatusForbidden, toto.V{"error": UserNotPermitError{}.Error()})
}
