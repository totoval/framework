package request

import (
	"net/http"
	"reflect"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
)

type websocketContext struct {
	*CommonContext
	*gin.Context
	ws *websocket.Conn
}

func (c *websocketContext) GinContext() *gin.Context {
	return c.Context
}
func (c *websocketContext) Request() *http.Request {
	return c.Context.Request
}
func (c *websocketContext) Writer() gin.ResponseWriter {
	return c.Context.Writer
}
func (c *websocketContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}
func (c *websocketContext) SetWriter(w gin.ResponseWriter) {
	c.Context.Writer = w
}

func (c *websocketContext) Params() gin.Params {
	return c.Context.Params
}
func (c *websocketContext) Accepted() []string {
	return c.Context.Accepted
}
func (c *websocketContext) Keys() map[string]interface{} {
	return c.Context.Keys
}
func (c *websocketContext) Errors() []*gin.Error {
	return c.Context.Errors
}
func (c *websocketContext) SetAuthClaim(claims *jwt.UserClaims) {
	c.Set(CONTEXT_CLAIM_KEY, claims)
}
func (c *websocketContext) AuthClaimID() (ID uint, exist bool) {
	claims, exist := c.Get(CONTEXT_CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}
func (c *websocketContext) SetIUserModel(iuser auth.IUser) {
	c.Set(CONTEXT_IUSER_KEY, iuser)
}
func (c *websocketContext) IUserModel() auth.IUser {
	iuser, exist := c.Get(CONTEXT_IUSER_KEY)

	var typeof reflect.Type
	if !exist {
		// default use config IUser model
		typeof = reflect.TypeOf(config.GetInterface("auth.model_ptr"))
	} else {
		// or use user set by middleware.IUser
		typeof = reflect.TypeOf(iuser.(auth.IUser))
	}

	ptr := reflect.New(typeof).Elem()
	val := reflect.New(typeof.Elem())
	ptr.Set(val)
	return ptr.Interface().(auth.IUser)
}
func (c *websocketContext) WS() *websocket.Conn {
	return c.ws
}
