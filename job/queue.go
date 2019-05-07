package job

// totoval do the broadcast it self, for compatible with the queue driver which doesn't support topic broadcasting
func topicName(j jobber) string {
	return "JOB-" + j.Name()
}
func channelName(j jobber) string {
	return j.Name()
}
