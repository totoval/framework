package middleware

import (
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/request/http/auth"
)

func IUser(userModelPtr auth.IUser) request.HandlerFunc {
	return func(c request.Context) {
		c.SetIUserModel(userModelPtr)

		c.Next()
	}
}
