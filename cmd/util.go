package cmd

import "github.com/fatih/color"

type TermLog struct {
	Message string
	Code Attribute
}

type Attribute interface {}

func (t *TermLog) Print (){
	if attribute, ok := t.Code.(color.Attribute); ok {
		d := color.New(attribute)
		d.Println(t.Message)
	}
}


const (
	CODE_SUCCESS color.Attribute = color.FgGreen
	CODE_WARNING color.Attribute = color.FgYellow
	CODE_INFO color.Attribute = color.FgBlue
	CODE_ERROR color.Attribute = color.FgRed

	MESSAGE_FINISHED = "Finished!"
)