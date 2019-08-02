package request

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/context"
	"github.com/totoval/framework/request/http/auth"
	"github.com/totoval/framework/utils/jwt"
)

type Context interface {
	// http context
	context.HttpContextor

	GinContext() *gin.Context

	SetAuthClaim(claims *jwt.UserClaims) //@todo move into a new interface
	SetIUserModel(iUser auth.IUser)      //@todo move into a new interface

	auth.Context
	auth.RequestIUser
}
