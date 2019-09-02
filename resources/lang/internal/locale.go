package internal

import (
	"errors"

	"github.com/go-playground/locales"
	"github.com/go-playground/universal-translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/totoval/framework/resources/lang"
)

type UnmarshalFunc = i18n.UnmarshalFunc

var LocalizerMap map[string]*i18n.Localizer

var localeMap map[string]*locale

type locale struct {
	languageName             string
	localizer                *i18n.Localizer
	localesTranslator        *locales.Translator
	validationTranslation    *lang.ValidationTranslation
	universalTranslator      ut.Translator
	validationRegisterStatus bool
}

func init() {
	localeMap = make(map[string]*locale)
	LocalizerMap = make(map[string]*i18n.Localizer)
}

func (l *locale) setValidationRegistered() *locale {
	l.validationRegisterStatus = true
	return l
}
func (l *locale) ValidationRegistered() bool {
	return l.validationRegisterStatus
}

// UnmarshalFunc  func(data []byte, v interface{}) error

func (l *locale) SetCustomTranslation(customTranslation *lang.CustomTranslation) *locale {
	bundle := i18n.NewBundle(language.English)

	for id, value := range *customTranslation {
		m := i18n.Message{
			ID:    id,
			Other: value,
		}
		if err := bundle.AddMessages(language.English, &m); err != nil {
			panic(errors.New("add Message error"))
		}
	}

	l.localizer = i18n.NewLocalizer(bundle, l.languageName)
	return l
}
func (l *locale) CustomTranslation() *i18n.Localizer {
	return l.localizer
}

func (l *locale) SetValidationTranslation(validationTranslation *lang.ValidationTranslation) *locale {
	localesTranslator := NewCommonLanguage(l.languageName)
	l.localesTranslator = &localesTranslator
	l.validationTranslation = validationTranslation
	return l
}

func (l *locale) SetLanguageName(langName string) *locale {
	l.languageName = langName
	return l
}

func (l *locale) SetUniversalTranslator() *locale {
	uttr := ut.New(NewCommonLanguage(l.languageName))
	l.universalTranslator, _ = uttr.GetTranslator(l.languageName)
	return l
}
func (l *locale) UniversalTranslator() ut.Translator {
	return l.universalTranslator
}

func AddLocale(langName string, customTranslation *lang.CustomTranslation, validationTranslation *lang.ValidationTranslation) {
	l := locale{}
	l.SetLanguageName(langName).SetCustomTranslation(customTranslation).SetValidationTranslation(validationTranslation).SetUniversalTranslator()
	localeMap[langName] = &l
}

func HasLocale(langName string) bool {
	if _, ok := localeMap[langName]; ok {
		return true
	}
	return false
}

func Locale(langName string) *locale {
	if HasLocale(langName) {
		return localeMap[langName]
	}
	return nil
}
