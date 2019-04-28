package queue

import "github.com/gogo/protobuf/proto"

type producer struct {
	topicName   string
	channelName string
	param       proto.Message
}

func NewProducer(topicName string, channelName string, param proto.Message) *producer {
	return &producer{
		topicName:   topicName,
		channelName: channelName,
		param:       param,
	}
}

func (p *producer) Push() error {
	pb, err := proto.Marshal(p.param)
	if err != nil {
		return err
	}
	return Queue().Push(p.topicName, pb)
}
