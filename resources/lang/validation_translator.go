package lang

import (
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

func AddLocale(langName string, validationTranslation *ValidationTranslation, customTranslation *CustomTranslation) {
	l := locale{}
	l.setLanguageName(langName).setCustomTranslation(customTranslation).setValidationTranslation(validationTranslation).setUniversalTranslator()
	localeMap[langName] = &l
}

func Translator(v *validator.Validate, langName string) (ut.Translator, error) {

	locale, ok := localeMap[langName]
	if !ok {
		return locale.universalTranslator, errors.New("validation translation not found")
	}

	if !locale.validationRegistered() {
		if err := registerDefaultTranslations(v, locale); err != nil {
			return locale.universalTranslator, err
		}
	}

	return locale.universalTranslator, nil
}
