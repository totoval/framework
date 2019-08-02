package middleware

import (
	l "github.com/totoval/framework/helpers/locale"
	"github.com/totoval/framework/request"

	"github.com/totoval/framework/config"
)

func Locale() request.HandlerFunc {
	return func(c request.Context) {
		locale := c.Request().Header.Get("locale")
		if locale == "" {
			locale = c.DefaultQuery("locale", config.GetString("app.locale"))
		}

		l.SetLocale(c, locale)

		c.Next()
	}
}
