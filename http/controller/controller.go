package controller

import (
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/request/http/auth"
	"github.com/totoval/framework/validator"
)

type Controller interface {
	Validate(c validator.Context, _validator interface{}, onlyFirstError bool) (isAbort bool)

	Authorize(c policy.Context, policies policy.Policier, action policy.Action) (permit bool, user auth.IUser)
}
