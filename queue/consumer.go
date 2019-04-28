package queue

import "github.com/gogo/protobuf/proto"

type consumer struct {
	topicName   string
	channelName string
	handler     func(param proto.Message) error
}

func NewConsumer(topicName string, channelName string, handler func(param proto.Message) error) *consumer {
	return &consumer{
		topicName:   topicName,
		channelName: channelName,
		handler:     handler,
	}
}
func (c *consumer) Pop() error {
	return Queue().Pop(c.producer.topicName, c.producer.channelName, func(body []byte) error {
		var p proto.Message
		if err := proto.Unmarshal(body, p); err != nil {
			return err
		}
		if err := c.handler(p); err != nil {
			return err
		}
		return nil
	})

}
