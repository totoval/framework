package debug

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/davecgh/go-spew/spew"

	"github.com/totoval/framework/console"
	"github.com/totoval/framework/helpers/zone"
)

func Dump(v ...interface{}) {
	log.Println(console.Sprintf(console.CODE_INFO, "%s - %s", zone.Now().String(), spew.Sdump(v...)), console.Sprintf(console.CODE_WARNING, string(debug.Stack())))
}

func DD(v ...interface{}) {
	log.Println(console.Sprintf(console.CODE_ERROR, "%s - %s%s", zone.Now().String(), spew.Sdump(v...), console.Sprintf(console.CODE_WARNING, string(debug.Stack()))))
	os.Exit(1)
}
