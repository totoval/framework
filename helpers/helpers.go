package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/resources/lang"
	"github.com/totoval/framework/utils/jwt"
	"math/rand"
	"net/http"
	"os"
	"time"
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


func L(c *gin.Context, messageID string, dataNlocale ...interface{}) string {
	l := lang.Locale(c)
	data := make(map[string]interface{})
	switch len(dataNlocale) {
	case 1:
		data = dataNlocale[0].(map[string]interface{})
		break
	case 2:
		l = dataNlocale[1].(string)
		break
	default:
	}

	return lang.Translate(messageID, data, l)
}

func Encrypt(secret string){
	//@todo
}

func Decrypt(){
	//@todo
}


const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
func RandNumberString(length uint) string {
	return random(int(length), numberBytes)
}
func RandString(length uint) string {
	return random(int(length), letterBytes)
}
func random(n int, bytes string)string{

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(bytes) {
			b[i] = bytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}