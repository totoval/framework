package view

import (
	"html/template"
	"sync"

	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/request"
)

func Initialize(r *request.Engine) {
	for _, tmpl := range engineTemplateMap.Get() {
		r.SetHTMLTemplate(tmpl)
	}
	log.Info(r)
}

func AddView(name string, content string) {
	tmpl := template.Must(template.New(name).Parse(content))
	engineTemplateMap.Set(tmpl)
}

type engineTemplate struct {
	Lock sync.RWMutex
	data []*template.Template
}

func newEngineTemplate() *engineTemplate {
	return &engineTemplate{
		data: []*template.Template{},
	}
}
func (et *engineTemplate) Get() []*template.Template {
	et.Lock.RLock()
	defer et.Lock.RUnlock()
	return et.data
}
func (et *engineTemplate) Set(tmpl *template.Template) {
	et.Lock.Lock()
	defer et.Lock.Unlock()
	et.data = append(et.data, tmpl)
}

var engineTemplateMap *engineTemplate

func init() {
	engineTemplateMap = newEngineTemplate()
}
