package websocket

import "sync"

type hubName = string
type connectionHubMap struct {
	lock sync.RWMutex
	hubs map[hubName]Hub
}

func newConnectionHubMap() *connectionHubMap {
	return &connectionHubMap{
		hubs: make(map[hubName]Hub),
	}
}

func (chs *connectionHubMap) Add(hub Hub) {
	chs.lock.Lock()
	defer chs.lock.Unlock()

	chs.hubs[hub.name()] = hub
}
func (chs *connectionHubMap) All() map[hubName]Hub {
	chs.lock.RLock()
	defer chs.lock.RUnlock()

	return chs.hubs
}
func (chs *connectionHubMap) Remove(hubName hubName) {
	chs.lock.Lock()
	defer chs.lock.Unlock()

	delete(chs.hubs, hubName)
}
