package helper

import (
    "github.com/gin-gonic/gin"
    "github.com/totoval/framework/resources/lang"

    "github.com/totoval/framework/config"

    "github.com/totoval/framework/resources/lang/internal"
)

func AddLocale(langName string, customTranslation *lang.CustomTranslation, validationTranslation *lang.ValidationTranslation) {
    internal.AddLocale(langName, customTranslation, validationTranslation)
}

func SetLocale(c *gin.Context, langName string) {
    c.Set("locale", langName)
}
func Locale(c *gin.Context) string {
    if contextLocale, exist := c.Get("locale"); exist {
        l := contextLocale.(string)
        return fallbackLocale(l)
    }
    configLocale := config.GetString("app.locale")
    return fallbackLocale(configLocale)
}

func fallbackLocale(langName string) string {
    if !internal.HasLocale(langName) {
        return config.GetString("app.fallback_locale", "en")
    }
    return langName
}
