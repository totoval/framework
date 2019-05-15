package queue

import (
	"errors"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/hub"
)

func init() {
	cmd.Add(&Listen{})
}

type Listen struct {
}

func (l *Listen) Command() string {
	return "queue:listen {listener_name}"
}

func (l *Listen) Description() string {
	return "Listen a given listener"
}

func (l *Listen) Handler(arg *cmd.Arg) error {
	listenerNamePtr, err := arg.Get("listener_name")
	if err != nil {
		return err
	}

	if listenerNamePtr == nil {
		return errors.New("listener_name is invalid")
	}

	hub.On(*listenerNamePtr)

	select {}

	return nil
}
