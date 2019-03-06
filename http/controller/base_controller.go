package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type Controller interface {
	Validate(c *gin.Context, _validator interface{}) bool
}

type BaseController struct{}

func (bc *BaseController) Validate(c *gin.Context, _validator interface{}) bool {
	if err := c.ShouldBindBodyWith(_validator, binding.JSON); err != nil {

		_ = err.(validator.ValidationErrors)
		//@todo translate
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return false
	}

	return true
}
