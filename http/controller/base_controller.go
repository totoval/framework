package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	"github.com/totoval/framework/helpers/locale"
	"github.com/totoval/framework/helpers/trans"
)

type Controller interface {
	Validate(c *gin.Context, _validator interface{}) bool
}

type BaseController struct{}

func (bc *BaseController) Validate(c *gin.Context, _validator interface{}, onlyFirstError bool) bool {
	if err := c.ShouldBindBodyWith(_validator, binding.JSON); err != nil {

		_err, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return false
		}

		v := binding.Validator.Engine().(*validator.Validate) // important, should be a new one for each request error

		errorResult := trans.ValidationTranslate(v, locale.Locale(c), _err)
		if onlyFirstError {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errorResult.First()})
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errorResult})
		}

		return false
	}

	return true
}
