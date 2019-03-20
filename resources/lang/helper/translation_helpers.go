package helper

import (
    "errors"

    "github.com/nicksnyder/go-i18n/v2/i18n"

    "github.com/totoval/framework/resources/lang/internal"
)

func CustomTranslate(messageID string, data map[string]interface{}, langName string) string {

    if l := internal.Locale(langName); l != nil {
        return l.CustomTranslation().MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData: data})
    }
    return errors.New("translation not found").Error()
}