package lang

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var UniversalTranslator *ut.UniversalTranslator

func Initialize(){
	_en := en.New()
	UniversalTranslator = ut.New(_en, _en)
}

func AddLocale(translator locales.Translator) error {
	return UniversalTranslator.AddTranslator(translator, true)
}

func Translator(translator string) (ut.Translator, bool) {
	return UniversalTranslator.GetTranslator(translator)
}