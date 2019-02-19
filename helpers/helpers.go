package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
)

func InSlice(needle interface{}, slice []interface{}) bool {
	for _, value := range slice {
		if value == needle {
			return true
		}
	}
	return false
}

func Dump(v ...interface{}){
	fmt.Println("########### Totoval Dump ###########")
	for _, value := range v {
		spew.Dump(value)
	}
	fmt.Println("########### Totoval Dump ###########")
}

func DD(v ...interface{}){
	fmt.Println("########### Totoval DD ###########")
	for _, value := range v {
		spew.Dump(value)
	}
	fmt.Println("########### Totoval DD ###########")
	os.Exit(1)
}