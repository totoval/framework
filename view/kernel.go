package view

import (
	"html/template"
	"sync"

	"github.com/totoval/framework/request"
)

func Initialize(r *request.Engine) {
	t := template.New("")
	for _, tmpl := range engineTemplateMap.Get() {
		t, _ = t.New(tmpl.name).Parse(tmpl.content)
	}
	r.SetHTMLTemplate(t)
}

func AddView(name string, content string) {
	engineTemplateMap.Set(&tmpl{
		name:    name,
		content: content,
	})
}

type tmpl struct {
	name    string
	content string
}
type engineTemplate struct {
	Lock sync.RWMutex
	data []*tmpl
}

func newEngineTemplate() *engineTemplate {
	return &engineTemplate{
		data: []*tmpl{},
	}
}
func (et *engineTemplate) Get() []*tmpl {
	et.Lock.RLock()
	defer et.Lock.RUnlock()
	return et.data
}
func (et *engineTemplate) Set(tmpl *tmpl) {
	et.Lock.Lock()
	defer et.Lock.Unlock()
	et.data = append(et.data, tmpl)
}

var engineTemplateMap *engineTemplate

func init() {
	engineTemplateMap = newEngineTemplate()
}
