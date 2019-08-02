package websocket

import "github.com/totoval/framework/request"

var hubs []*connectionHub

type connectionHub struct {
	msgChan  chan *Msg
	isClosed bool
	request.Context
}

func (ch *connectionHub) init(c request.Context) {
	ch.msgChan = make(chan *Msg, 256)
	ch.isClosed = false
	ch.Context = c
	hubs = append(hubs, ch)
}
func (ch *connectionHub) Send(msg *Msg) {
	ch.msgChan <- msg
}
func (ch *connectionHub) Broadcast(msg *Msg) {
	for i, hub := range hubs {
		if !hub.available() {
			hubs = append(hubs[:i], hubs[i+1:]...) // remove hubs[i]
			continue
		}

		hub.Send(msg)
	}
}
func (ch *connectionHub) getChan() chan *Msg {
	return ch.msgChan
}
func (ch *connectionHub) close() {
	ch.isClosed = true
}
func (ch *connectionHub) available() bool {
	return !ch.isClosed
}
