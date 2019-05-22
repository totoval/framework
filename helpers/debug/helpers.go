package debug

import (
	"errors"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ztrue/tracerr"

	"github.com/totoval/framework/console"
	"github.com/totoval/framework/helpers/log"
)

func Dump(v ...interface{}) {
	console.Println(console.CODE_ERROR, spew.Sdump(v...))
	debugPrint(errors.New("====== Totoval Debug ======"), 2)
}

func DD(v ...interface{}) {
	Dump(v...)
	os.Exit(1)
}

func debugPrint(err error, startFrom int) {
	traceErr := tracerr.Wrap(err)
	frameList := tracerr.StackTrace(traceErr)
	if startFrom > len(frameList) {
		_ = log.Error(err)
	}
	traceErr = tracerr.CustomError(err, frameList[startFrom:len(frameList)-2])
	tracerr.PrintSourceColor(traceErr)
}
