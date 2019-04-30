package queue

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/helpers/str"
	message "github.com/totoval/framework/queue/protocol_buffers"
)

type consumer struct {
	topicName   string
	channelName string
	paramPtr    proto.Message //for param retrieve
	handler     func(paramPtr proto.Message) error
}

func NewConsumer(topicName string, channelName string, paramPtr proto.Message, handler func(paramPtr proto.Message) error) *consumer {
	return &consumer{
		topicName:   topicName,
		channelName: channelName,
		paramPtr:    paramPtr,
		handler:     handler,
	}
}
func (c *consumer) Pop() error {
	return Queue().Pop(c.topicName, c.channelName, func(body []byte) error {
		// exact message
		msg := message.Message{}
		if err := proto.Unmarshal(body, &msg); err != nil {
			return err
		}

		// increase tried
		msg.Tried = msg.Tried + 1

		//@todo here should proceed the hash， write it into the message
		// msg.Hash = "xxxxxxx"
		msg.Hash = str.RandString(40)

		log.Println(msg)

		// exact param
		if err := proto.Unmarshal(msg.Param, c.paramPtr); err != nil {
			return err
		}

		defer c.Failed(msg)

		if err := c.handler(c.paramPtr); err != nil {
			return err
		}
		return nil
	})

}

func (c *consumer) Failed(msg message.Message) {
	if err := recover(); err != nil {
		log.Println(err)

		newMsg := msg
		newMsg.Retries = newMsg.Retries - 1

		fmt.Println(msg.Retries)

		if msg.Retries <= 0 {
			// if database save failed, then push into queue again? or log?
			if err := c.failedToDatabase(c.topicName, c.channelName, &msg); err != nil {
				log.Println(msg)
				log.Println(newMsg, "failedtodatabase processed failed")
				newMsg.Retries = 1
				goto DB_FAILED
			}
			return
		}

	DB_FAILED:
		if err := c.failedToQueue(&newMsg); err != nil {
			if err := c.failedToDatabase(c.topicName, c.channelName, &newMsg); err != nil {
				// error!!!! processed failed
				log.Println(newMsg, "failedtoqueue processed failed")
			}
		}
		return
	}
}

func (c *consumer) failedToQueue(msg *message.Message) error {
	return push(c.topicName, c.channelName, msg)
}

func (c *consumer) failedToDatabase(topicName string, channelName string, msg *message.Message) error {
	return failedProcessor.FailedToDatabase(topicName, channelName, msg)
}
