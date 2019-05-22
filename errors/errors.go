package errors

import (
	"github.com/ztrue/tracerr"

	"github.com/totoval/framework/logs"
)

func ErrPrint(err error, startFrom int) {
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	if startFrom > len(frameList) {
		logs.Println(logs.ERROR, err.Error(), nil)
	}
	traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	logs.Println(logs.ERROR, tracerr.SprintSourceColor(traceErr), nil)
}
