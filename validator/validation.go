package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"gopkg.in/go-playground/validator.v9"

	"github.com/totoval/framework/helpers/locale"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/trans"
)

type Validation struct {
}

//@todo Deprecated, for compatible with v0.10.0
func (v *Validation) Validate(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return v.ValidateJSON(c, requestDataPtr, onlyFirstError)
}

// not consume request body
func (v *Validation) ValidateJSONMulti(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(func(rdp interface{}) error {
		return c.ShouldBindBodyWith(rdp, binding.JSON)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateXMLMulti(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(func(rdp interface{}) error {
		return c.ShouldBindBodyWith(rdp, binding.XML)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateYAMLMulti(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(func(rdp interface{}) error {
		return c.ShouldBindBodyWith(rdp, binding.YAML)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateMsgPackMulti(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(func(rdp interface{}) error {
		return c.ShouldBindBodyWith(rdp, binding.MsgPack)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateProtoBufMulti(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(func(rdp interface{}) error {
		return c.ShouldBindBodyWith(rdp, binding.ProtoBuf)
	}, c, requestDataPtr, onlyFirstError)
}

// consume request body
func (v *Validation) ValidateJSON(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if hasBodyBytesKey(c) {
		return v.ValidateJSONMulti(c, requestDataPtr, onlyFirstError)
	}

	return validate(func(rdp interface{}) error {
		return c.ShouldBindWith(rdp, binding.JSON)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateXML(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if hasBodyBytesKey(c) {
		return v.ValidateXMLMulti(c, requestDataPtr, onlyFirstError)
	}

	return validate(func(rdp interface{}) error {
		return c.ShouldBindWith(rdp, binding.XML)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateYAML(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if hasBodyBytesKey(c) {
		return v.ValidateYAMLMulti(c, requestDataPtr, onlyFirstError)
	}

	return validate(func(rdp interface{}) error {
		return c.ShouldBindWith(rdp, binding.YAML)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateMsgPack(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if hasBodyBytesKey(c) {
		return v.ValidateMsgPackMulti(c, requestDataPtr, onlyFirstError)
	}

	return validate(func(rdp interface{}) error {
		return c.ShouldBindWith(rdp, binding.MsgPack)
	}, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateProtoBuf(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if hasBodyBytesKey(c) {
		return v.ValidateProtoBufMulti(c, requestDataPtr, onlyFirstError)
	}

	return validate(func(rdp interface{}) error {
		return c.ShouldBindWith(rdp, binding.ProtoBuf)
	}, c, requestDataPtr, onlyFirstError)
}
func hasBodyBytesKey(c Context) bool {
	if _, ok := c.Get(gin.BodyBytesKey); ok {
		return true
	}
	return false
}

func (v *Validation) ValidateUri(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(c.ShouldBindUri, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateQuery(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(c.ShouldBindQuery, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateForm(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(c.ShouldBind, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateFormPost(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(c.ShouldBind, c, requestDataPtr, onlyFirstError)
}
func (v *Validation) ValidateFormMultipart(c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	return validate(c.ShouldBind, c, requestDataPtr, onlyFirstError)
}

func validate(shouldBindFunc func(rdp interface{}) error, c Context, requestDataPtr interface{}, onlyFirstError bool) (isAbort bool) {
	if err := shouldBindFunc(requestDataPtr); err != nil {

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
