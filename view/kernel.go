package view

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/helpers/debug"
)

var templateList []*template.Template

func Initialize(r *gin.Engine) {
	for _, tmpl := range templateList {
		r.SetHTMLTemplate(tmpl)
	}
}

func AddView(name string, content string) {
	debug.Dump(name)
	tmpl := template.Must(template.New(name).Parse(content))
	templateList = append(templateList, tmpl)
}
