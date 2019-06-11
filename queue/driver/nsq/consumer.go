package nsq

import (
	_nsq "github.com/nsqio/go-nsq"
)

type consumer struct {
	topicName        string
	channelName      string
	hashTopicChannel hashTopicChannel
	cfg              *_nsq.Config
	c                *_nsq.Consumer
}
