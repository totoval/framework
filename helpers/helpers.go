package helpers

import (
	"github.com/totoval/framework/helpers/trans"
	"github.com/totoval/framework/request"

	"github.com/totoval/framework/helpers/locale"
)

func L(c *request.Context, messageID string, dataNlocale ...interface{}) string {
	l := locale.Locale(c)
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

	return trans.CustomTranslate(messageID, data, l)
}

//func Encrypt(secret string) {
//	//@todo
//}
//
//func Decrypt() {
//	//@todo
//}
