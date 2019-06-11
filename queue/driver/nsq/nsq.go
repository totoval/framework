package nsq

import (
	"fmt"

	_nsq "github.com/nsqio/go-nsq"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/hash"
)

func NewNsq(connection string) *nsq {
	n := new(nsq)
	n.consumerList = make(map[hashTopicChannel]*consumer)
	n.setConnection(connection)

	nsqlookupdHttpConnectionArgsList := n.nsqlookupdHttpConnectionArgsList()
	if len(nsqlookupdHttpConnectionArgsList) < 0 {
		panic("nsqlookupd http connection is not set")
	}

	nsqdTcpConnectionArgsList := n.nsqdTcpConnectionArgsList()
	if len(nsqdTcpConnectionArgsList) < 0 {
		panic("nsqd tcp connection is not set")
	}

	if err := n.producerConnect(); err != nil {
		panic(err)
	}

	return n
}

type hashTopicChannel = string

func (n *nsq) addConnectedProducer(p *producer) {
	n.connectedProducerList = append(n.connectedProducerList, p)
}
func (n *nsq) addConnectedConsumer(c *consumer) {
	n.connectedConsumerList = append(n.connectedConsumerList, c)
}

func (n *nsq) hashTopicChannel(topicName string, channelName string) hashTopicChannel {
	return hash.Md5(fmt.Sprintf("%s||||||%s", topicName, channelName))
}

func (n *nsq) consumerConnect(topicName string, channelName string) (hashTopicChannel hashTopicChannel, err error) {

	c := new(consumer)
	c.topicName = topicName
	c.channelName = channelName
	c.hashTopicChannel = n.hashTopicChannel(c.topicName, c.channelName)
	c.cfg = _nsq.NewConfig()
	c.c, err = _nsq.NewConsumer(c.topicName, c.channelName, c.cfg)
	if err != nil {
		return "", err
	}

	n.setConsumer(c.hashTopicChannel, c) // concurrent map writes
	// n.consumerList[c.hashTopicChannel] = c
	return c.hashTopicChannel, nil
}
func (n *nsq) producerConnect() (err error) {
	addr := n.nsqdTcpConnectionArgsList()[0] //@todo default producer to nsqd index 0

	p := new(producer)
	p.cfg = _nsq.NewConfig()

	if p.p, err = _nsq.NewProducer(addr, p.cfg); err != nil {
		return err
	}

	n.producer = p

	n.addConnectedProducer(n.producer)

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
func (n *nsq) httpHost() string {
	return n.config("lookupd.http.host")
}
func (n *nsq) httpPort() string {
	return n.config("lookupd.http.port")
}
func (n *nsq) tcpHost() string {
	return n.config("lookupd.tcp.host")
}
func (n *nsq) tcpPort() string {
	return n.config("lookupd.tcp.port")
}
func (n *nsq) config(key string) string {
	value := config.GetString("queue.connections." + n.connection() + "." + key)
	if value == "" {
		panic("queue " + key + " parse error")
	}
	return value
}
func (n *nsq) configMapList(key string) []map[string]interface{} { // nsqd
	value := config.Get("queue.connections." + n.connection() + "." + key)
	if value == "" {
		panic("queue " + key + " parse error")
	}
	return value.([]map[string]interface{})
}
func (n *nsq) nsqdTcpConnectionArgsList() (argsList []string) {
	for _, conf := range n.configMapList("nsqd") {
		confTcp := conf["tcp"].(map[string]interface{})
		connectionArgs := confTcp["host"].(string) + ":" + confTcp["port"].(string)
		argsList = append(argsList, connectionArgs)
	}

	return argsList
}
func (n *nsq) nsqlookupdHttpConnectionArgsList() (argsList []string) {
	for _, conf := range n.configMapList("nsqlookupd") {
		confTcp := conf["http"].(map[string]interface{})
		connectionArgs := confTcp["host"].(string) + ":" + confTcp["port"].(string)
		argsList = append(argsList, connectionArgs)
	}

	return argsList
}
