package hub

import "github.com/totoval/framework/queue"

// totoval do the broadcast it self, for compatible with the queue driver which doesn't support topic broadcasting
func topicName(e Eventer, l Listener, supportBroadCasting func() bool) string {
	if supportBroadCasting() {
		return "EVENT-" + EventName(e)
	}
	return "EVENT-" + EventName(e) + "-" + channelName(l)
}
func channelName(l Listener) string {
	return l.Name()
}
func RegisterQueue() {
	for e, llist := range hub {
		for _, l := range llist {
			if err := queue.Queue().Register(topicName(event(e), l, queue.Queue().SupportBroadCasting), channelName(l)); err != nil {
				panic(err)
			}
		}
	}
}
