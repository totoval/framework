package job

import "github.com/totoval/framework/queue"

// totoval do the broadcast it self, for compatible with the queue driver which doesn't support topic broadcasting
func topicName(j jobber) string {
	return "JOB-" + j.Name()
}
func channelName(j jobber) string {
	return j.Name()
}
func RegisterQueue() {
	for _, j := range jobMap {
		if err := queue.Queue().Register(topicName(j), channelName(j)); err != nil {
			panic(err)
		}
	}
}
