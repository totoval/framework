package errors

import (
	"fmt"

	"github.com/ztrue/tracerr"

	"github.com/totoval/framework/logs"
)

func ErrPrintln(err error, fields logs.Field) {
	startFrom := 2
	if err == nil {
		return
	}
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	if startFrom > len(frameList) || len(frameList)-2 <= 0 {
		logs.Println(logs.ERROR, err.Error(), fields)
	}

	traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	traceErr = tracerr.CustomError(err, frameList)

	if fields == nil {
		fields = logs.Field{}
	}
	fields["totoval_trace"] = tracerr.SprintSource(traceErr, 0)
	logs.Println(logs.ERROR, err.Error(), fields)
}

func ErrPrint(err error, startFrom int, fields logs.Field) string {
	if err == nil {
		return ""
	}
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	//if startFrom > len(frameList) {
	//	return fmt.Sprint(err.Error(), fields)
	//}
	//
	//traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	traceErr = tracerr.CustomError(err, frameList)

	if fields == nil {
		fields = logs.Field{}
	}
	fields["totoval_trace"] = tracerr.SprintSource(traceErr)
	return fmt.Sprint(err.Error(), fields)
}
