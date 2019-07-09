package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/model"
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/validator"
)

type Controller interface {
	Validate(c *request.Context, _validator interface{}, onlyFirstError bool) (isAbort bool)

	Authorize(c *request.Context, policies policy.Policier, action policy.Action) (permit bool, user model.IUser)

	Scan(c *request.Context) (isAbort bool)
	User() model.IUser
	UserId(c *request.Context) (userId uint, isAbort bool)
}

type BaseController struct {
	ginContext *gin.Context
	policy.Authorization
	auth.RequestUser
	validator.Validation
}
