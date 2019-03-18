package lang

import (
	"encoding/json"
	"io/ioutil"
	"strings"

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
	languageName string
	localizer *i18n.Localizer
	localesTranslator *locales.Translator
	validationTranslation *ValidationTranslation
	universalTranslator ut.Translator
	validationRegisterStatus bool
}

func init(){
	localeMap = make(map[string]*locale)
	LocalizerMap = make(map[string]*i18n.Localizer)

	// init lang files
	langFileFormat := "json"
	langUnmarshalFunc := json.Unmarshal
	langFileDirName := "resources/lang"
	initializeLangFiles(langFileDirName, langFileFormat, langUnmarshalFunc)

}

func initializeLangFiles(dirName string, langFileFormat string, langUnmarshalFunc UnmarshalFunc){
	bundle := &i18n.Bundle{DefaultLanguage: language.English} //@todo xxx there's maybe a bug
	bundle.RegisterUnmarshalFunc(langFileFormat, langUnmarshalFunc)

	if dirName[len(dirName)-1:] != "/" {
		dirName = dirName + "/"
	}
	fileArr, err := ioutil.ReadDir(dirName)
	if err != nil{
		panic(err)
	}
	for _, value := range fileArr {
		if value.IsDir() {
			continue
		}
		fileName := value.Name()
		if fileName[len(fileName)-len("."+langFileFormat):] != ("."+langFileFormat) {
			continue
		}
		//@todo bundle.AddMessage to support custom language!!!


		bundle.MustLoadMessageFile(dirName + fileName)
		langName := strings.Replace(fileName, "."+langFileFormat, "", 1) //@todo if file name = "test.json.json", there may be a bug
		LocalizerMap[langName] = i18n.NewLocalizer(bundle, langName)

		//addLocale(langName)
	}
}

func (l *locale)setValidationRegistered() *locale {
	l.validationRegisterStatus = true
	return l
}
func (l *locale)validationRegistered() bool {
	return l.validationRegisterStatus
}

func (l *locale)setCustomTranslation(dirName string, langFileFormat string, langUnmarshalFunc UnmarshalFunc) *locale {
	bundle := &i18n.Bundle{DefaultLanguage: language.English} //@todo xxx there's maybe a default language bug
	bundle.RegisterUnmarshalFunc(langFileFormat, langUnmarshalFunc)

	if dirName[len(dirName)-1:] != "/" {
		dirName = dirName + "/"
	}
	bundle.MustLoadMessageFile(dirName + l.languageName + "." + langFileFormat)
	l.localizer = i18n.NewLocalizer(bundle, l.languageName)
	return l
}

func (l *locale)setValidationTranslation(validationTranslation *ValidationTranslation) *locale {
	localesTranslator := NewCommonLanguage(l.languageName)
	l.localesTranslator = &localesTranslator
	l.validationTranslation = validationTranslation
	return l
}

func (l *locale)setLanguageName(langName string) *locale {
	l.languageName = langName
	return l
}

func (l *locale)setUniversalTranslator() *locale {
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
func supportedLocale() []string {
	return []string{}
}

func fallbackLocale(locale string) string {
	if !hasLocale(locale){
		return config.GetString("app.fallback_locale", "en")
	}
	return locale
}

func hasLocale(langName string)bool{
	if _, ok := localeMap[langName]; ok {
		return true
	}
	return false
}

func Translate(messageID string, data map[string]interface{}, langName string) string {
	return LocalizerMap[langName].MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData:data})
}

