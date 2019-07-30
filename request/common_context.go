package request

import (
	"net/http"
	"reflect"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
)

type CommonContext struct {
	*gin.Context
}

func (c *CommonContext) GinContext() *gin.Context {
	return c.Context
}
func (c *CommonContext) Request() *http.Request {
	return c.Context.Request
}
func (c *CommonContext) Writer() gin.ResponseWriter {
	return c.Context.Writer
}
func (c *CommonContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}
func (c *CommonContext) SetWriter(w gin.ResponseWriter) {
	c.Context.Writer = w
}

func (c *CommonContext) Params() gin.Params {
	return c.Context.Params
}
func (c *CommonContext) Accepted() []string {
	return c.Context.Accepted
}
func (c *CommonContext) Keys() map[string]interface{} {
	return c.Context.Keys
}
func (c *CommonContext) Errors() []*gin.Error {
	return c.Context.Errors
}
func (c *CommonContext) SetAuthClaim(claims *jwt.UserClaims) {
	c.Set(CONTEXT_CLAIM_KEY, claims)
}
func (c *CommonContext) AuthClaimID() (ID uint, exist bool) {
	claims, exist := c.Get(CONTEXT_CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}
func (c *CommonContext) SetIUserModel(iuser auth.IUser) {
	c.Set(CONTEXT_IUSER_KEY, iuser)
}
func (c *CommonContext) IUserModel() auth.IUser {
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
