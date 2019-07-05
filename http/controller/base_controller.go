package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/model"
	"github.com/totoval/framework/policy"
	"github.com/totoval/framework/validator"
)

type Controller interface {
	Validate(c *gin.Context, _validator interface{}, onlyFirstError bool) (isAbort bool)

	Authorize(c *gin.Context, policies policy.Policier, action policy.Action) (permit bool, user model.IUser)

	Scan(c *gin.Context) (isAbort bool)
	User() model.IUser
	UserId(c *gin.Context) (userId uint, isAbort bool)
}

type BaseController struct {
	policy.Authorization
	auth.RequestUser
	validator.Validation
}
