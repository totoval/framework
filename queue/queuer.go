package queue

func init() {
	//@todo use different config to new different queuer
	//@todo memory, nsq, rabbitmq
	//@todo failure controll
}

var queue queuer

func Queue() queuer {
	return queue
}

type queuer interface {
	producerer
	consumerer
}

type producerer interface {
	Push(topicName string, body []byte) (err error)
}
type consumerer interface {
	Pop(topicName string, channelName string, handler func(body []byte) error) (err error)
}
