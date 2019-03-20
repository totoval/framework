package middleware

import (
	"github.com/gin-gonic/gin"

	l "github.com/totoval/framework/helpers/locale"

	"github.com/totoval/framework/config"
)

func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Request.Header.Get("locale")
		if locale == "" {
			locale = c.DefaultQuery("locale", config.GetString("app.locale"))
		}

		l.SetLocale(c, locale)

		c.Next()
	}
}
