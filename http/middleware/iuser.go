package middleware

import (
	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/request"
)

func IUser(userModelPtr auth.IUser) request.HandlerFunc {
	return func(c request.Context) {
		c.SetIUserModel(userModelPtr)

		c.Next()
	}
}
