package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/helpers/toto"
)

type UserNotPermitError struct{}

func (e UserNotPermitError) Error() string {
	return "user has no permission"
}

func Middleware(policy Policier, action Action, c Context, params []gin.Param) {

	// get route url param
	routeParamMap := make(map[string]string)
	for _, param := range params {
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

func forbid(c Context) {
	c.JSON(http.StatusForbidden, toto.V{"error": UserNotPermitError{}.Error()})
}
