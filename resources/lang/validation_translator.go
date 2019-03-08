package lang

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/nl"
	"github.com/go-playground/locales/pt_BR"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"

	en_translation "gopkg.in/go-playground/validator.v9/translations/en"
	fr_translation "gopkg.in/go-playground/validator.v9/translations/fr"
	id_translation "gopkg.in/go-playground/validator.v9/translations/id"
	ja_translation "gopkg.in/go-playground/validator.v9/translations/ja"
	nl_translation "gopkg.in/go-playground/validator.v9/translations/nl"
	pt_BR_translation "gopkg.in/go-playground/validator.v9/translations/pt_BR"
	zh_translation "gopkg.in/go-playground/validator.v9/translations/zh"
	zh_tw_translation "gopkg.in/go-playground/validator.v9/translations/zh_tw"
)

type ValidationTranslator interface {
	RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error)
	Locale() string
	LocalesTranslator() locales.Translator
}

var UniversalTranslator *ut.UniversalTranslator
var validationTranslatorMap map[string]ValidationTranslator

func init(){
	_en := en.New()
	_zh := zh.New()
	_ja := ja.New()
	_fr := fr.New()
	_id := id.New()
	_nl := nl.New()
	_pt_BR := pt_BR.New()
	_zh_tw := zh_Hant_TW.New()

	UniversalTranslator = ut.New(_en, _en, _zh, _ja, _fr, _id, _nl, _pt_BR, _zh_tw)

	validationTranslatorMap = make(map[string]ValidationTranslator)
}

func AddLocale(validationTranslator ValidationTranslator) error {
	validationTranslatorMap[validationTranslator.Locale()] = validationTranslator

	return UniversalTranslator.AddTranslator(validationTranslator.LocalesTranslator(), true)
}

type RegisterDefaultTranslations func (v *validator.Validate, trans ut.Translator) (err error)
func Translator(v *validator.Validate, translator string) (trans ut.Translator, found bool) {
	trans, found = UniversalTranslator.GetTranslator(translator)

	var rdt RegisterDefaultTranslations

	switch translator {
	case "en":
		rdt = en_translation.RegisterDefaultTranslations
		break
	case "zh":
		rdt = zh_translation.RegisterDefaultTranslations
		break
	case "ja":
		rdt = ja_translation.RegisterDefaultTranslations
		break
	case "fr":
		rdt = fr_translation.RegisterDefaultTranslations
		break
	case "id":
		rdt = id_translation.RegisterDefaultTranslations
		break
	case "nl":
		rdt = nl_translation.RegisterDefaultTranslations
		break
	case "pt_BR":
		rdt = pt_BR_translation.RegisterDefaultTranslations
		break
	case "zh_tw":
		rdt = zh_tw_translation.RegisterDefaultTranslations
		break
	}

	if validationTranslator, ok := validationTranslatorMap[translator]; ok {
		rdt = validationTranslator.RegisterDefaultTranslations
	}

	if err := rdt(v, trans); err != nil{
		return
	}

	return
}
