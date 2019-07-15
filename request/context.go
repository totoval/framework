package request

import (
	"reflect"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
)

const CONTEXT_CLAIM_KEY = "TOTOVAL_CONTEXT_CLAIM"
const CONTEXT_IUSER_KEY = "TOTOVAL_CONTEXT_IUSER"

type Context struct {
	*gin.Context
}
type HandlerFunc func(*Context)

func (c *Context) SetAuthClaim(claims *jwt.UserClaims) {
	c.Set(CONTEXT_CLAIM_KEY, claims)
}
func (c *Context) AuthClaimID() (ID uint, exist bool) {
	claims, exist := c.Get(CONTEXT_CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}
func (c *Context) SetIUserModel(iuser auth.IUser) {
	c.Set(CONTEXT_IUSER_KEY, iuser)
}
func (c *Context) IUserModel() auth.IUser {
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

func ConvertHandlers(handlers []HandlerFunc) (ginHandlers []gin.HandlerFunc) {
	for _, h := range handlers {
		handler := h // must new a variable for `range's val`, or the `val` in anonymous funcs will be overwrite every loop

		ginHandlers = append(ginHandlers, func(c *gin.Context) {
			totovalContext := Context{Context: c}
			handler(&totovalContext)
		})
	}
	return
}
