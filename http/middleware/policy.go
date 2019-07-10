package middleware

import (
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/request"
)

func Policy(_policy policy.Policier, action policy.Action) request.HandlerFunc {
	return func(c *request.Context) {
		policy.Middleware(_policy, action, c, c.Params)
	}
}
