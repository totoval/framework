package queue

import (
	"github.com/golang/protobuf/ptypes"

	message "github.com/totoval/framework/queue/protocol_buffers"
)

type queueRetry struct {
	p              producer
	paramProtoByte []byte
	hash           string
}

func newQueryRetry(fq FailedQueuer) *queueRetry {
	return &queueRetry{
		p: producer{
			topicName:   fq.RetryTopic(),
			channelName: fq.RetryChannel(),
			retries:     fq.RetryRetries(),
			delay:       fq.RetryDelay(),
		},
		paramProtoByte: fq.RetryParamProtoBytes(),
		hash:           fq.RetryHash(),
	}
}
func (r *queueRetry) retry() error {
	return push(r.p.topicName, r.p.channelName, &message.Message{

		Hash:     r.hash,
		Param:    r.paramProtoByte,
		Retries:  r.p.retries,
		PushedAt: ptypes.TimestampNow(),
		Delay:    ptypes.DurationProto(r.p.delay),
		Tried:    0,
	})
}

func Forget(id uint) error {
	// delete the retried failed queue
	if err := failedProcessor.DeleteQueueById(id); err != nil {
		return err
	}
	return nil
}
func Flush() error {
	//@todo to be done
	panic("need implements")
}
func Retry(id uint) error {

	// find the failed queue
	fq, err := failedProcessor.FailedQueueById(id)
	if err != nil {
		return err
	}

	// retry the failed queue
	if err := newQueryRetry(fq).retry(); err != nil {
		return err
	}

	// delete the retried failed queue
	if err := failedProcessor.DeleteQueueById(id); err != nil {
		return err
	}

	return nil
}
