package policy

import (
	"github.com/totoval/framework/context"
	"github.com/totoval/framework/request/http/auth"
)

type Context interface {
	context.LifeCycleContextor
	context.ResponseContextor
	auth.Context
	auth.RequestIUser
}
