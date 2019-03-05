package email

import (
	"github.com/totoval/framework/notification"
)
type Email struct {
	driver notification.Driver
}
func (e *Email)Prepare (prepareFunc func() (notification.Messager)) notification.Driver {
	e.driver.SetMessager(prepareFunc())
	return e.driver
}