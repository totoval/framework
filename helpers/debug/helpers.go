package debug

import (
	"os"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"

	"github.com/totoval/framework/console"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/log"
)

func Dump(v ...interface{}) {
	log.Println("DUMP", console.Sprintf(console.CODE_INFO, "%s - %s", zone.Now().String(), spew.Sdump(v...)), console.Sprintf(console.CODE_WARNING, string(debug.Stack())))
}

func DD(v ...interface{}) {
	Dump(v...)
	os.Exit(1)
}
