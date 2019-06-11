package nsq

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	_nsq "github.com/nsqio/go-nsq"

	"github.com/totoval/framework/helpers/zone"
)

type nsq struct {
	producer              *producer
	conn                  string
	consumerList          map[hashTopicChannel]*consumer
	Lock                  sync.RWMutex
	connectedProducerList []*producer
	connectedConsumerList []*consumer
}

func (n *nsq) Register(topicName string, channelName string) (err error) {
	if !_nsq.IsValidTopicName(topicName) {
		return errors.New("topic name is invalid")
	}
	if !_nsq.IsValidChannelName(channelName) {
		return errors.New("channel name is invalid")
	}

	// default register to nsqd index 0
	resp, err := http.Post(fmt.Sprintf("%s/channel/create?topic=%s&channel=%s", n.nsqlookupdHttpConnectionArgsList()[0], topicName, channelName), "", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response error, code: %d", resp.StatusCode))
	}

	return nil
}
func (n *nsq) Unregister(topicName string, channelName string) (err error) { return }

func (n *nsq) Push(topicName string, channelName string, delay zone.Duration, body []byte) (err error) {
	// for SupportBroadCasting queue driver channelName is not used
	if delay > 0 {
		return n.producer.p.DeferredPublish(topicName, delay, body)
	}
	return n.producer.p.Publish(topicName, body)
}

func (n *nsq) Pop(topicName string, channelName string, handler func(hash string, body []byte) (handlerErr error), maxInFlight int) (err error) {
	htc, err := n.consumerConnect(topicName, channelName)
	if err != nil {
		return err
	}
	//h := n.hashTopicChannel(topicName, channelName)
	// n.consumerList[h].c.AddHandler(_nsq.HandlerFunc(func(message *_nsq.Message) error {
	n.consumer(htc).c.AddHandler(_nsq.HandlerFunc(func(message *_nsq.Message) error {
		return handler(string(message.ID[:]), message.Body)
	}))
	if err := n.consumer(htc).c.ConnectToNSQDs(n.nsqdTcpConnectionArgsList()); err != nil {
		return err
	}
	n.consumer(htc).c.ChangeMaxInFlight(maxInFlight)
	n.addConnectedConsumer(n.consumer(htc))
	return nil
}

func (n *nsq) SupportBroadCasting() bool {
	return true
}

func (n *nsq) Close() (err error) {
	// stop producer
	for _, p := range n.connectedProducerList {
		p.p.Stop()
	}

	// stop consumer
	for _, c := range n.connectedConsumerList {
		c.c.Stop()
	}

	return nil
}
