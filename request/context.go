package request

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/context"
	"github.com/totoval/framework/utils/jwt"
)

type Context interface {
	// http context
	context.HttpContextor

	GinContext() *gin.Context
	
	SetAuthClaim(claims *jwt.UserClaims)
	AuthClaimID() (ID uint, exist bool)
	SetIUserModel(iuser auth.IUser)
	IUserModel() auth.IUser
}

type WsContext interface {
	Context
	WS() *websocket.Conn
}
