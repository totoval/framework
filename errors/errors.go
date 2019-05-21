package errors

import (
	"github.com/ztrue/tracerr"

	"github.com/totoval/framework/log"
)

func ErrPrint(prefix string, err error, startFrom int) {
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	if startFrom > len(frameList) {
		log.Println(prefix, err)
	}
	traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	log.Println(prefix, tracerr.SprintSourceColor(traceErr))
}
