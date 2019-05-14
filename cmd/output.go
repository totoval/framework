package cmd

import (
	"log"

	"github.com/fatih/color"
)

func Println(code Attribute, msg string) {
	if _, err := color.New(code.(color.Attribute)).Println(msg); err != nil {
		log.Println(err)
	}
}

type Attribute interface{}

const (
	CODE_SUCCESS color.Attribute = color.FgGreen
	CODE_WARNING color.Attribute = color.FgYellow
	CODE_INFO    color.Attribute = color.FgBlue
	CODE_ERROR   color.Attribute = color.FgRed

	MESSAGE_FINISHED = "Finished!"
)
