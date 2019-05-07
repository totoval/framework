package driver

import (
	"fmt"
	"sync"
	"time"

	_nsq "github.com/nsqio/go-nsq"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/hash"
)

func NewNsq(connection string) *nsq {
	n := new(nsq)
	n.consumerList = make(map[hashTopicChannel]*consumer)
	n.setConnection(connection)
	addr := n.connectionArgs()

	n.producer = new(producer)
	n.producer.cfg = _nsq.NewConfig()

	var err error
	if n.producer.p, err = _nsq.NewProducer(addr, n.producer.cfg); err != nil {
		panic(err)
	}

	return n
}

type hashTopicChannel = string

type producer struct {
	cfg *_nsq.Config
	p   *_nsq.Producer
}
type consumer struct {
	topicName        string
	channelName      string
	hashTopicChannel hashTopicChannel
	cfg              *_nsq.Config
	c                *_nsq.Consumer
}

type nsq struct {
	producer     *producer
	conn         string
	consumerList map[hashTopicChannel]*consumer
	Lock         sync.RWMutex
}

func (n *nsq) Push(topicName string, channelName string, delay time.Duration, body []byte) (err error) {
	if delay > 0 {
		return n.producer.p.DeferredPublish(topicName, delay, body)
	}
	return n.producer.p.Publish(topicName, body)
}

func (n *nsq) Pop(topicName string, channelName string, handler func(hash string, body []byte) error) (err error) {
	if err := n.connect(topicName, channelName); err != nil {
		return err
	}
	h := n.HashTopicChannel(topicName, channelName)
	// n.consumerList[h].c.AddHandler(_nsq.HandlerFunc(func(message *_nsq.Message) error {
	n.consumer(h).c.AddHandler(_nsq.HandlerFunc(func(message *_nsq.Message) error {
		return handler(string(message.ID[:]), message.Body)
	}))
	return n.consumer(h).c.ConnectToNSQD(n.connectionArgs())
	//return n.consumerList[h].c.ConnectToNSQD(n.connectionArgs())
}

func (n *nsq) HashTopicChannel(topicName string, channelName string) hashTopicChannel {
	return hash.Md5(fmt.Sprintf("%s||||||%s", topicName, channelName))
}

func (n *nsq) connect(topicName string, channelName string) (err error) {

	c := new(consumer)
	c.topicName = topicName
	c.channelName = channelName
	c.hashTopicChannel = n.HashTopicChannel(c.topicName, c.channelName)
	c.cfg = _nsq.NewConfig()
	c.c, err = _nsq.NewConsumer(c.topicName, c.channelName, c.cfg)
	if err != nil {
		return err
	}

	n.setConsumer(c.hashTopicChannel, c) // concurrent map writes
	// n.consumerList[c.hashTopicChannel] = c
	return nil
}

func (n *nsq) consumer(h hashTopicChannel) *consumer {
	n.Lock.RLock()
	defer n.Lock.RUnlock()
	return n.consumerList[h]
}
func (n *nsq) setConsumer(h hashTopicChannel, c *consumer) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	n.consumerList[h] = c
}

func (n *nsq) setConnection(connection string) {
	n.conn = connection
}
func (n *nsq) connection() string {
	return n.conn
}
func (n *nsq) host() string {
	return n.config("host")
}
func (n *nsq) port() string {
	return n.config("port")
}
func (n *nsq) config(key string) string {
	value := config.GetString("queue.connections." + n.connection() + "." + key)
	if value == "" {
		panic("queue " + key + " parse error")
	}
	return value
}
func (n *nsq) connectionArgs() string {
	return n.host() + ":" + n.port()
}
