package locale

import (
	"github.com/totoval/framework/resources/lang/helper"

	"github.com/totoval/framework/resources/lang"
)

func AddLocale(langName string, customTranslation *lang.CustomTranslation, validationTranslation *lang.ValidationTranslation) {
	helper.AddLocale(langName, customTranslation, validationTranslation)
}
func SetLocale(c lang.Context, langName string) {
	helper.SetLocale(c, langName)
}
func Locale(c lang.Context) string {
	return helper.Locale(c)
}
