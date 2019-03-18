package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/resources/lang"

	"github.com/totoval/framework/config"
)

func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Request.Header.Get("locale")
		if locale == "" {
			locale = c.DefaultQuery("locale", config.GetString("app.locale"))
		}

		lang.SetLocale(c, locale)

		c.Next()
	}
}
