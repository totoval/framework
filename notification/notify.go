package notification

type Notify struct {
	driver Driver
}
func (n *Notify)Prepare (prepareFunc func() (Messager)) Driver {
	n.driver.SetMessager(prepareFunc())
	return n.driver
}
func (n *Notify)SetDriver(driver Driver){
	n.driver = driver
}