package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/resources/lang"
	"github.com/totoval/framework/utils/jwt"
	"net/http"
	"os"
	"unicode/utf8"
)

func InSlice(needle interface{}, slice interface{}) bool {
	for _, value := range slice.([]interface{}) {
		if value == needle {
			return true
		}
	}
	return false
}

func Dump(v ...interface{}) {
	fmt.Println("########### Totoval Dump ###########")
	for _, value := range v {
		spew.Dump(value)
	}
	fmt.Println("########### Totoval Dump ###########")
}

func DD(v ...interface{}) {
	fmt.Println("########### Totoval DD ###########")
	for _, value := range v {
		spew.Dump(value)
	}
	fmt.Println("########### Totoval DD ###########")
	os.Exit(1)
}

func AuthClaimsID(c *gin.Context) uint {
	claims, exist := c.Get("claims")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not login"})
		return 0
	}

	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r)
}


func L(c *gin.Context, messageID string, locale ...string) string {
	l := lang.Locale(c)
	if len(locale) > 0{
		l = locale[0]
	}
	return lang.Translate(messageID, l)
}
