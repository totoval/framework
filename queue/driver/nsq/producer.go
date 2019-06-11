package nsq

import (
	_nsq "github.com/nsqio/go-nsq"
)

type producer struct {
	cfg *_nsq.Config
	p   *_nsq.Producer
}
