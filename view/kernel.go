package view

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

var templateList []*template.Template

func Initialize(r *gin.Engine) {
	for _, tmpl := range templateList {
		r.SetHTMLTemplate(tmpl)
	}
}

func AddView(name string, content string) {
	tmpl := template.Must(template.New(name).Parse(content))
	templateList = append(templateList, tmpl)
}
