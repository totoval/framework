package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/totoval/framework/cmd"
)

func Dump(v ...interface{}) {
	cmd.Println(cmd.CODE_WARNING, spew.Sdump(v...))
}

func DD(v ...interface{}) {
	cmd.Println(cmd.CODE_ERROR, spew.Sdump(v...))
	os.Exit(1)
}
