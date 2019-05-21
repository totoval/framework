package log

import (
	"fmt"
	"log"
	"os"

	"github.com/totoval/framework/helpers/zone"
)

func Println(prefix string, v ...interface{}) {
	l := log.New(os.Stdout, fmt.Sprintf("%s [%s] ", zone.Now().String(), prefix), 0)
	l.Println(v...)
}
