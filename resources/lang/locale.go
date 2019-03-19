package lang

import (
    "errors"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/locales"
    "github.com/go-playground/universal-translator"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"

    "github.com/totoval/framework/config"
)

type UnmarshalFunc = i18n.UnmarshalFunc

var LocalizerMap map[string]*i18n.Localizer

var localeMap map[string]*locale

type locale struct {
    languageName             string
    localizer                *i18n.Localizer
    localesTranslator        *locales.Translator
    validationTranslation    *ValidationTranslation
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
func (l *locale) validationRegistered() bool {
    return l.validationRegisterStatus
}

// UnmarshalFunc  func(data []byte, v interface{}) error
type CustomTranslation map[string]string
func (l *locale) setCustomTranslation(customTranslation *CustomTranslation) *locale {
    bundle := &i18n.Bundle{DefaultLanguage: language.English}

    for id, value := range *customTranslation {
        m := i18n.Message{
            ID: id,
            Other: value,
        }
        if err := bundle.AddMessages(language.English, &m); err != nil {
            panic(errors.New("add Message error"))
        }
    }

    l.localizer = i18n.NewLocalizer(bundle, l.languageName)
    return l
}

func (l *locale) setValidationTranslation(validationTranslation *ValidationTranslation) *locale {
    localesTranslator := NewCommonLanguage(l.languageName)
    l.localesTranslator = &localesTranslator
    l.validationTranslation = validationTranslation
    return l
}

func (l *locale) setLanguageName(langName string) *locale {
    l.languageName = langName
    return l
}

func (l *locale) setUniversalTranslator() *locale {
    uttr := ut.New(NewCommonLanguage(l.languageName))
    l.universalTranslator, _ = uttr.GetTranslator(l.languageName)
    return l
}

func SetLocale(c *gin.Context, locale string) {
    c.Set("locale", locale)
}
func Locale(c *gin.Context) string {
    if contextLocale, exist := c.Get("locale"); exist {
        l := contextLocale.(string)
        return fallbackLocale(l)
    }
    configLocale := config.GetString("app.locale")
    return fallbackLocale(configLocale)
}

func fallbackLocale(locale string) string {
    if !hasLocale(locale) {
        return config.GetString("app.fallback_locale", "en")
    }
    return locale
}

func hasLocale(langName string) bool {
    if _, ok := localeMap[langName]; ok {
        return true
    }
    return false
}

func Translate(messageID string, data map[string]interface{}, langName string) string {
    return localeMap[langName].localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData: data})
}
