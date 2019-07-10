package request

import (
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/utils/jwt"
)

const CONTEXT_CLAIM_KEY = "TOTOVAL_CONTEXT_CLAIM"

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
