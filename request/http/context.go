package http

import (
	"net/http"
	"reflect"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/request/http/auth"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
)

const CONTEXT_CLAIM_KEY = "TOTOVAL_CONTEXT_CLAIM"
const CONTEXT_IUSER_MODEL_KEY = "TOTOVAL_CONTEXT_IUSER_MODEL"

type httpContext struct {
	*gin.Context
	*auth.RequestUser
}

func (c *httpContext) GinContext() *gin.Context {
	return c.Context
}
func (c *httpContext) Request() *http.Request {
	return c.Context.Request
}
func (c *httpContext) Writer() gin.ResponseWriter {
	return c.Context.Writer
}
func (c *httpContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}
func (c *httpContext) SetWriter(w gin.ResponseWriter) {
	c.Context.Writer = w
}

func (c *httpContext) Params() gin.Params {
	return c.Context.Params
}
func (c *httpContext) Accepted() []string {
	return c.Context.Accepted
}
func (c *httpContext) Keys() map[string]interface{} {
	return c.Context.Keys
}
func (c *httpContext) Errors() []*gin.Error {
	return c.Context.Errors
}
func (c *httpContext) SetAuthClaim(claims *jwt.UserClaims) {
	c.Set(CONTEXT_CLAIM_KEY, claims)
}
func (c *httpContext) AuthClaimID() (ID uint, exist bool) {
	claims, exist := c.Get(CONTEXT_CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}
func (c *httpContext) SetIUserModel(iuser auth.IUser) {
	c.Set(CONTEXT_IUSER_MODEL_KEY, iuser)
}
func (c *httpContext) IUserModel() auth.IUser {
	iuser, exist := c.Get(CONTEXT_IUSER_MODEL_KEY)

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

func (c *httpContext) ScanUserWithJSON() (isAbort bool) {
	if err := c.ScanUser(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, toto.V{"error": err})
		return true
	}
	return false
}
func ConvertContext(c *gin.Context) *httpContext {
	_c := &httpContext{Context: c, RequestUser: &auth.RequestUser{}}

	_c.RequestUser.SetContext(_c)

	return _c
}
