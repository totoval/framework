package queue

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/logs"
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
	return Queue().Pop(c.topicName, c.channelName, func(hash string, body []byte) error {
		// exact message
		msg := message.Message{}
		if err := proto.Unmarshal(body, &msg); err != nil {
			return err
		}

		// increase tried
		msg.Tried = msg.Tried + 1

		// log hash
		msg.Hash = hash

		log.Info("queue msg received", logs.Field{
			"msg": msg,
		})

		// exact param
		if err := proto.Unmarshal(msg.Param, c.paramPtr); err != nil {
			return err
		}

		defer c.Failed(msg)

		if err := c.handler(c.paramPtr); err != nil {
			panic(err)
		}
		return nil
	}, config.GetInt("queue.max_in_flight"))

}

func (c *consumer) Failed(msg message.Message) {
	if err := recover(); err != nil {
		//fmt.Println(err)

		newMsg := msg
		newMsg.Retries = newMsg.Retries - 1

		//fmt.Println(msg.Retries)

		_ = log.Error(errors.New("queue msg processed error"), logs.Field{
			"msg": msg,
		})

		if msg.Retries <= 0 {
			// if database save failed, then push into queue again? or log?
			if err := c.failedToDatabase(c.topicName, c.channelName, &msg, err); err != nil {

				_ = log.Error(errors.New("failedtodatabase processed failed"), logs.Field{
					"new_msg": newMsg,
				})
				newMsg.Retries = 1
				goto DB_FAILED
			}
			return
		}

	DB_FAILED:
		if err := c.failedToQueue(&newMsg); err != nil {
			if err := c.failedToDatabase(c.topicName, c.channelName, &newMsg, err); err != nil {
				// error!!!! processed failed
				_ = log.Error(errors.New("failedtoqueue processed failed"), logs.Field{
					"new_msg": newMsg,
				})
			}
		}
		return
	}
}

func (c *consumer) failedToQueue(msg *message.Message) error {
	return push(c.topicName, c.channelName, msg)
}

func (c *consumer) failedToDatabase(topicName string, channelName string, msg *message.Message, err interface{}) error {
	errStr := fmt.Sprint(err)
	if _err, ok := err.(error); ok {
		errStr = log.ErrorStr(_err)
	}
	return failedProcessor.FailedToDatabase(topicName, channelName, msg, errStr)
}
