package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/request"
)

func Recovery() request.HandlerFunc {
	return func(c *request.Context) {
		gin.Recovery()(c.Context)
	}
}
