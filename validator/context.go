package validator

import (
	"github.com/totoval/framework/context"
	"github.com/totoval/framework/resources/lang"
)

type Context interface {
	context.RequestBindingContextor
	context.ResponseContextor
	lang.Context
}
