package hub

// totoval do the broadcast it self, for compatible with the queue driver which doesn't support topic broadcasting
func topicName(e Eventer, l Listener) string {
	return "EVENT-" + EventName(e) + "-" + l.Name()
}
func channelName(l Listener) string {
	return l.Name()
}
