package lang

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/locales"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)


type ValidationTranslator interface {
	RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error)
	Locale() string
	LocalesTranslator() locales.Translator
}

func AddValidationTranslation(langName string, translation *ValidationTranslation) {
	//@todo change this config
	langFileFormat := "json"
	langUnmarshalFunc := json.Unmarshal
	langFileDirName := "resources/lang"

	// add locale
	l := locale{}
	l.setLanguageName(langName).setCustomTranslation(langFileDirName, langFileFormat, langUnmarshalFunc).setValidationTranslation(translation).setUniversalTranslator()
	localeMap[langName] = &l

}

func Translator(v *validator.Validate, langName string) (ut.Translator, error) {

	locale, ok := localeMap[langName]
	if !ok {
		return locale.universalTranslator, errors.New("validation translation not found")
	}


	//@todo 每次register会报错，如果已经注册了就不register了
	if !locale.validationRegistered() {
		if err := registerDefaultTranslations(v, locale); err != nil {
			return locale.universalTranslator, err
		}
	}

	return locale.universalTranslator, nil
}
