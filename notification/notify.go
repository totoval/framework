package notification

type Notify struct {
	driver Driver
}

func (n *Notify) Prepare(prepareFunc func(m Messager) Messager) Driver {
	message := &Message{}
	n.driver.SetMessager(prepareFunc(message))
	return n.driver
}
func (n *Notify) SetDriver(driver Driver) {
	n.driver = driver
}
