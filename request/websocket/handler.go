package websocket

import (
	"net/http"

	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/http/controller"
)

type Handler interface {
	OnMessage(hub Hub, msg *Msg)
	Loop(hub Hub) error

	OnPing(hub Hub, appData string)
	OnPong(hub Hub, appData string)
	OnClose(hub Hub, code int, text string)

	config
	controller.Controller
}
type Hub interface {
	Send(msg *Msg)
	Broadcast(msg *Msg)
}
type config interface {
	ReadBufferSize() int
	WriteBufferSize() int
	CheckOrigin(r *http.Request) bool
	WriteTimeout() zone.Duration
	ReadTimeout() zone.Duration
	MaxMessageSize() int64
}
