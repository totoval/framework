package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	"github.com/totoval/framework/resources/lang"
)

type Controller interface {
	Validate(c *gin.Context, _validator interface{}) bool
}

type BaseController struct{}

func (bc *BaseController) Validate(c *gin.Context, _validator interface{}) bool {
	if err := c.ShouldBindBodyWith(_validator, binding.JSON); err != nil {

		_err, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return false
		}

		v := binding.Validator.Engine().(*validator.Validate) // important, should be a new one for each request error
		translator, err := lang.Translator(v, lang.Locale(c))
		if err != nil{
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return false
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": _err.Translate(translator)})
		return false
	}

	return true
}
