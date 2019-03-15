package code_verify

import (
	"github.com/totoval/framework/helpers/str"
	"time"
)

type CodeVerify struct {
	CodeLen uint
	ExpiredDuration time.Duration
	validationList map[string]code
}
type code struct {
	value string
	createdAt time.Time
}

func (cv *CodeVerify) Generate(index string) string {
	c := code{
		value: str.RandNumberString(cv.CodeLen),
		createdAt: time.Now(),
	}
	if cv.validationList == nil{
		cv.validationList = make(map[string]code)
	}
	cv.validationList[index] = c

	return c.value
}
func (cv *CodeVerify) Verify(index string, code string) bool {
	if cv.validationList[index].createdAt.Add(cv.ExpiredDuration).Sub(time.Now()) > 0 && cv.validationList[index].value == code {
		cv.delete(index)
		return true
	}
	return false
}
func (cv *CodeVerify) delete(index string){
	delete(cv.validationList, index)
}