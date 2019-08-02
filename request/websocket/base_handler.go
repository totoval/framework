package websocket

import (
	"net/http"

	"github.com/totoval/framework/helpers/zone"
)

type BaseHandler struct {
}

func (bh *BaseHandler) OnPing(hub Hub, appData string) {

}

func (bh *BaseHandler) OnPong(hub Hub, appData string) {

}

func (bh *BaseHandler) OnClose(hub Hub, code int, text string) {

}

func (bh *BaseHandler) ReadBufferSize() int {
	return 1024
}
func (bh *BaseHandler) WriteBufferSize() int {
	return 1024
}
func (bh *BaseHandler) CheckOrigin(r *http.Request) bool {
	return true
}
func (bh *BaseHandler) WriteTimeout() zone.Duration {
	return 10 * zone.Second
}
func (bh *BaseHandler) ReadTimeout() zone.Duration {
	return 60 * zone.Second
}
func (bh *BaseHandler) MaxMessageSize() int64 {
	return 512
}
