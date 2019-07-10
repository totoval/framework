package policy

import (
	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/context"
)

type Context interface {
	context.LifeCycleContextor
	auth.Context
}
