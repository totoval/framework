package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/resources/lang"
)



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


