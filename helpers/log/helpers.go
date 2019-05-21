package log

import (
	"github.com/totoval/framework/console"
	"github.com/totoval/framework/errors"
	"github.com/totoval/framework/log"
)

func Err(err error) error {
	errors.ErrPrint("ERR", err, 2)
	return err
}
func Info(v ...interface{}) {
	log.Println("INFO", console.Sprint(console.CODE_INFO, v...))
}
func Warn(v ...interface{}) {
	log.Println("WARN", console.Sprint(console.CODE_WARNING, v...))
}
func Success(v ...interface{}) {
	log.Println("SUCC", console.Sprint(console.CODE_WARNING, v...))
}
