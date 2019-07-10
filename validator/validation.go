package validator

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"gopkg.in/go-playground/validator.v9"

	"github.com/totoval/framework/helpers/locale"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/trans"
)

type Validation struct {
}

func (v *Validation) Validate(c Context, _validator interface{}, onlyFirstError bool) (isAbort bool) {
	if err := c.ShouldBindBodyWith(_validator, binding.JSON); err != nil {

		_err, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, toto.V{"error": err.Error()})
			return false
		}

		v := binding.Validator.Engine().(*validator.Validate) // important, should be a new one for each request error

		errorResult := trans.ValidationTranslate(v, locale.Locale(c), _err)
		if onlyFirstError {
			c.JSON(http.StatusUnprocessableEntity, toto.V{"error": errorResult.First()})
		} else {
			c.JSON(http.StatusUnprocessableEntity, toto.V{"error": errorResult})
		}

		return false
	}

	return true
}
