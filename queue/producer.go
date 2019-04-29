package queue

import (
	"github.com/gogo/protobuf/proto"

	message "github.com/totoval/framework/queue/protocol_buffers"
)

type producer struct {
	topicName   string
	channelName string
	param       proto.Message
	retries     uint32
}

func NewProducer(topicName string, channelName string, param proto.Message, retries uint32) *producer {
	return &producer{
		topicName:   topicName,
		channelName: channelName,
		param:       param,
		retries:     retries,
	}
}

func (p *producer) Push() error {
	// compress param
	paramPb, err := proto.Marshal(p.param)
	if err != nil {
		return err
	}

	// compress message
	return push(p.topicName, p.channelName, &message.Message{
		//Param: &any.Any{
		//	TypeUrl: "github.com/totoval/framework/queue/" + proto.MessageName(p.param), // p.param should be pointer
		//	Value:   paramPb,
		//},
		Param:   paramPb,
		Retries: p.retries,
	})
}

func push(topicName string, channelName string, msg *message.Message) error {
	// compress message
	messagePb, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return Queue().Push(topicName, channelName, messagePb)
}
