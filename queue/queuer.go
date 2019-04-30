package queue

import (
	"time"

	"github.com/totoval/framework/queue/driver"
)

func Initialize() {
	//@todo use different config to new different queuer
	//@todo memory, nsq, rabbitmq
	//@todo failure controll
	setQueue(driver.NewNsq("nsq"))
}

var queue queuer

func Queue() queuer {
	return queue
}
func setQueue(q queuer) {
	queue = q
}

type queuer interface {
	producerer
	consumerer
}

type producerer interface {
	Push(topicName string, channelName string, delay time.Duration, body []byte) (err error)
}
type consumerer interface {
	Pop(topicName string, channelName string, handler func(body []byte) error) (err error)
}
