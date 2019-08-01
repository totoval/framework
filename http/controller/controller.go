package controller

import (
	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/validator"
)

type Controller interface {
	Validate(c validator.Context, _validator interface{}, onlyFirstError bool) (isAbort bool)

	Authorize(c auth.Context, policies policy.Policier, action policy.Action) (permit bool, user auth.IUser)

	Scan(c auth.Context) (isAbort bool)
	User(c auth.Context) auth.IUser
	UserId(c auth.Context) (userId uint, isAbort bool)
}
