package lang

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"strings"
)

type UnmarshalFunc = i18n.UnmarshalFunc

var LocalizerMap map[string]*i18n.Localizer

func init(){
	LocalizerMap = make(map[string]*i18n.Localizer)

	langFileFormat := "json"
	langUnmarshalFunc := json.Unmarshal
	dirName := "resources/lang"

	initializeLangFiles(dirName, langFileFormat, langUnmarshalFunc)
}

func initializeLangFiles(dirName string, langFileFormat string, langUnmarshalFunc UnmarshalFunc){
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc(langFileFormat, langUnmarshalFunc)

	if dirName[len(dirName)-1:] != "/" {
		dirName = dirName + "/"
	}
	fileArr, err := ioutil.ReadDir(dirName)
	if err != nil{
		panic(err)
	}
	for _, value := range fileArr {
		bundle.MustLoadMessageFile(dirName + value.Name())
		langName := strings.Replace(value.Name(), "."+langFileFormat, "", 1) // if file name = "test.json.json", there may be a bug
		LocalizerMap[langName] = i18n.NewLocalizer(bundle, langName)
	}
}

func Translate(messageID string, langName string) string {
	return LocalizerMap[langName].MustLocalize(&i18n.LocalizeConfig{MessageID: messageID/*"auth.register.failed_existed"*/})
}