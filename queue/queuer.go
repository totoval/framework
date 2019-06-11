package queue

import (
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/queue/driver/nsq"
)

func Initialize() {
	//@todo use different config to new different queuer
	//@todo memory, nsq, rabbitmq
	setQueue(nsq.NewNsq("nsq"))
	initializeFailedProcessor()
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
	registerer
	SupportBroadCasting() bool // For listener which has not start before emit, it will not receive the event. in SupportBroadCasting mode
	Close() (err error)
}

type registerer interface {
	Register(topicName string, channelName string) (err error)
	Unregister(topicName string, channelName string) (err error)
}

type producerer interface {
	Push(topicName string, channelName string, delay zone.Duration, body []byte) (err error)
}
type consumerer interface {
	// when handlerErr not nil, queue should send REQ (re-queue), when nil, send FIN (finish)
	Pop(topicName string, channelName string, handler func(hash string, body []byte) (handlerErr error), maxInFlight int) (err error)
}
