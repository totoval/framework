package websocket

import (
	"sync"
)

var channelMap *channels

type channelName = string

const totovalDefaultChannelName = "totovalDefaultChannel"

func init() {
	channelMap = &channels{
		channels: make(map[channelName]*connectionHubMap),
	}
}

type channels struct {
	lock     sync.RWMutex
	channels map[channelName]*connectionHubMap
}

func (c *channels) LeaveAll(hub Hub) {
	for _, hubMap := range c.channels {
		hubMap.Remove(hub.name())
	}
}

func (c *channels) Hubs(channelName channelName) *connectionHubMap {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if hubMap, ok := c.channels[channelName]; ok {
		if hubMap == nil {
			return newConnectionHubMap()
		}
		return hubMap
	}

	return newConnectionHubMap()
}

func (c *channels) Leave(channelName channelName, hub Hub) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if hubMap, ok := c.channels[channelName]; ok {
		if hubMap == nil {
			return
		}
		c.channels[channelName].Remove(hub.name())
	}
	return
}

func (c *channels) Join(channelName string, hub Hub) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if hubMap, ok := c.channels[channelName]; ok {
		if hubMap == nil {
			c.channels[channelName] = newConnectionHubMap()
		}
	} else {
		c.channels[channelName] = newConnectionHubMap()
	}

	c.channels[channelName].Add(hub)
}
