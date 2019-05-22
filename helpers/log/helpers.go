package log

import (
	"github.com/totoval/framework/errors"
	"github.com/totoval/framework/logs"
)

func Error(err error) error {
	errors.ErrPrint(err, 2)
	return err
}

func Info(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.INFO, msg, fields)
}
func Warn(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.WARN, msg, fields)
}
func Fatal(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.FATAL, msg, fields)
}
func Debug(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.DEBUG, msg, fields)
}
func Panic(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.PANIC, msg, fields)
}
func Trace(msg string, v ...logs.Field) {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	logs.Println(logs.TRACE, msg, fields)
}
