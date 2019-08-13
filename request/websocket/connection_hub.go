package websocket

import (
	"sync"

	"github.com/totoval/framework/request"
)

var hubs *connectionHubSlice

type connectionHubSlice struct {
	lock sync.RWMutex
	hubs []*connectionHub
}

func (chs *connectionHubSlice) Append(hub *connectionHub) {
	chs.lock.Lock()
	defer chs.lock.Unlock()
	chs.hubs = append(chs.hubs, hub)
}
func (chs *connectionHubSlice) Get() []*connectionHub {
	chs.lock.RLock()
	defer chs.lock.RUnlock()
	return chs.hubs
}
func (chs *connectionHubSlice) Drop(index int) {
	chs.lock.Lock()
	defer chs.lock.Unlock()
	chs.hubs = append(chs.hubs[:index], chs.hubs[index+1:]...) // remove hubs[i]
}

type connectionHub struct {
	msgChan  chan *Msg
	isClosed bool
	request.Context
}

func (ch *connectionHub) init(c request.Context) {
	ch.msgChan = make(chan *Msg, 256)
	ch.isClosed = false
	ch.Context = c
	hubs.Append(ch)
}
func (ch *connectionHub) Send(msg *Msg) {
	ch.msgChan <- msg
}
func (ch *connectionHub) Broadcast(msg *Msg) {
	for _, hub := range hubs.Get() {
		if !hub.available() {
			continue
		}

		hub.Send(msg)
	}
}
func (ch *connectionHub) getChan() chan *Msg {
	return ch.msgChan
}
func (ch *connectionHub) close() {
	close(ch.msgChan)
	ch.isClosed = true
	for i, hub := range hubs.Get() {
		if hub == ch {
			hubs.Drop(i) // remove hubs[i]
		}
	}
}
func (ch *connectionHub) available() bool {
	return !ch.isClosed
}
