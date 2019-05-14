package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/totoval/framework/console"
)

func Dump(v ...interface{}) {
	console.Println(console.CODE_INFO, spew.Sdump(v...))
}

func DD(v ...interface{}) {
	console.Println(console.CODE_ERROR, spew.Sdump(v...))
	os.Exit(1)
}
