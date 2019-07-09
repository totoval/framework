package queue

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
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
	return Queue().Pop(c.topicName, c.channelName, func(hash string, body []byte) (handlerErr error) {
		// exact message
		msg := message.Message{}
		if err := proto.Unmarshal(body, &msg); err != nil {
			return err
		}

		// increase tried
		msg.Tried = msg.Tried + 1

		// log hash
		msg.Hash = hash

		log.Info("queue msg received", toto.V{
			"msg": msg,
		})

		// exact param
		if err := proto.Unmarshal(msg.Param, c.paramPtr); err != nil {
			return err
		}

		defer c.Failed(msg, &handlerErr)

		if err := c.handler(c.paramPtr); err != nil {
			log.Info(err.Error())
			panic(err)
		}

		// if handler panic or return err, will not return nil
		return nil
	}, config.GetInt("queue.max_in_flight"))

}

func (c *consumer) Failed(msg message.Message, handlerErrPtr *error) {
	if hErr := recover(); hErr != nil {
		//fmt.Println(err)

		newMsg := msg
		newMsg.Retries = newMsg.Retries - 1
		// delay the every retries more 3 minutes
		newMsg.Delay = ptypes.DurationProto(zone.Duration(msg.Tried) * 3 * zone.Minute)

		//fmt.Println(msg.Retries)

		_ = log.Error(errors.New("queue msg processed error"), toto.V{
			"msg":   msg,
			"error": hErr,
		})

		if msg.Retries <= 0 {
			// if database save failed, then push into queue again? or log?
			if err := c.failedToDatabase(c.topicName, c.channelName, &msg, hErr); err != nil {

				_ = log.Error(errors.New("failedtodatabase processed failed"), toto.V{
					"new_msg": newMsg,
				})
				newMsg.Retries = 1
				goto DB_FAILED
			}
			return
		}

	DB_FAILED:
		if err := c.failedToQueue(&newMsg, hErr, handlerErrPtr); err != nil {
			if err := c.failedToDatabase(c.topicName, c.channelName, &newMsg, hErr); err != nil {
				// error!!!! processed failed
				_ = log.Error(errors.New("failedtoqueue processed failed"), toto.V{
					"new_msg": newMsg,
				})
			}
		}
		return
	}
}

func (c *consumer) failedToQueue(msg *message.Message, handlerErr interface{}, handlerErrPtr *error) error {
	// repush the failed message and add retries
	*handlerErrPtr = nil
	return push(c.topicName, c.channelName, msg)

	//// when handlerErr not nil, queue should send REQ (re-queue), when nil, send FIN (finish)
	//err := convertInterfaceErr(handlerErr)
	//*handlerErrPtr = err
	//return nil
}

func (c *consumer) failedToDatabase(topicName string, channelName string, msg *message.Message, err interface{}) error {
	return failedProcessor.FailedToDatabase(topicName, channelName, msg, convertInterfaceErr(err).Error())
}

func convertInterfaceErr(err interface{}) error {
	errStr := fmt.Sprint(err)
	if _err, ok := err.(error); ok {
		errStr = log.ErrorStr(_err)
	}
	return errors.New(errStr)
}
