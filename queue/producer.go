package queue

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	message "github.com/totoval/framework/queue/protocol_buffers"
)

type producer struct {
	topicName   string
	channelName string
	param       proto.Message
	retries     uint32
	delay       time.Duration
}

func NewProducer(topicName string, channelName string, param proto.Message, retries uint32, delay time.Duration) *producer {
	return &producer{
		topicName:   topicName,
		channelName: channelName,
		param:       param,
		retries:     retries,
		delay:       delay,
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
		Param:     paramPb,
		Retries:   p.retries,
		CreatedAt: ptypes.TimestampNow(),
		Delay:     ptypes.DurationProto(p.delay),
		Tried:     0,
	})
}

func push(topicName string, channelName string, msg *message.Message) error {
	// compress message
	messagePb, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return Queue().Push(topicName, channelName, time.Duration(msg.Delay.GetSeconds())*time.Second+time.Duration(msg.Delay.GetNanos())*time.Nanosecond, messagePb)
}
