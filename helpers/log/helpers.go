package log

import (
	"github.com/totoval/framework/errors"
	"github.com/totoval/framework/logs"
)

func Error(err error, v ...logs.Field) error {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	errors.ErrPrintln(err, fields)
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
func ErrorStr(err error, v ...logs.Field) string {
	var fields logs.Field
	if len(v) > 0 {
		fields = v[0]
	}
	return errors.ErrPrint(err, 2, fields)
}
