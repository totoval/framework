package hub

import (
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/totoval/framework/queue"
)

func Emit(e Eventer) (errs map[ListenerName]error) {
	// For listener which has not start before emit, it will not receive the event. in SupportBroadCasting mode
	if queue.Queue().SupportBroadCasting() {
		eventListenerList := eventListener(e)
		if len(eventListenerList) <= 0 {
			errs = make(map[ListenerName]error)
			errs["nil"] = errors.New("listener doesnt't exist")
			return errs
		}

		l := eventListenerList[0]
		if err := queue.NewProducer(topicName(e, l, queue.Queue().SupportBroadCasting), channelName(l), e.paramData(), l.Retries(), l.Delay()).Push(); err != nil {
			errs = make(map[ListenerName]error)
			errs[channelName(l)] = err
		}
		return errs
	}

	// push to multi Listener
	for _, l := range eventListener(e) {
		if err := queue.NewProducer(topicName(e, l, queue.Queue().SupportBroadCasting), channelName(l), e.paramData(), l.Retries(), l.Delay()).Push(); err != nil {
			if errs == nil {
				errs = make(map[ListenerName]error)
			}

			errs[channelName(l)] = err
		}
	}
	return errs
}

func On(listenerName ListenerName) {

	l := listenerMap[listenerName]
	if l == nil {
		panic(errors.New("listener " + listenerName + " doesn't exist"))
	}

	for _, e := range l.Subscribe() {
		err := queue.NewConsumer(topicName(e, l, queue.Queue().SupportBroadCasting), channelName(l), e.ParamProto(), func(paramPtr proto.Message) error {
			if err := l.Construct(paramPtr); err != nil {
				return err
			}

			if err := l.Handle(); err != nil {
				return err
			}

			return nil
		}).Pop()
		if err != nil {
			panic(err)
		}
	}
}
