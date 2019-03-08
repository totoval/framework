package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/totoval/framework/resources/lang"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
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

		v := binding.Validator.Engine().(*validator.Validate)
		translator, found := lang.Translator(v, lang.Locale(c))
		if !found {
			log.Println(errors.New("translation locale xx.New() not found, use default translation"))
			//c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			//return false
		}


		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": _err.Translate(translator)})
		return false
	}

	return true
}
