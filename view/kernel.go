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
	lock sync.RWMutex
	data []*tmpl
}

func newEngineTemplate() *engineTemplate {
	return &engineTemplate{
		data: []*tmpl{},
	}
}
func (et *engineTemplate) Get() []*tmpl {
	et.lock.RLock()
	defer et.lock.RUnlock()
	return et.data
}
func (et *engineTemplate) Set(tmpl *tmpl) {
	et.lock.Lock()
	defer et.lock.Unlock()
	et.data = append(et.data, tmpl)
}

var engineTemplateMap *engineTemplate

func init() {
	engineTemplateMap = newEngineTemplate()
}
