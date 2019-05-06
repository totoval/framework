package hub

var hub map[EventerName][]Listener
var listenerMap map[ListenerName]Listener

func init() {
	hub = make(map[EventerName][]Listener)
	listenerMap = make(map[ListenerName]Listener)
}

func Make(e Eventer) {
	e.SetParam(e.ParamProto()) // for init
	hub[EventName(e)] = nil
}

// register Listener to hub
func Register(l Listener) {
	//if len(l.Listen()) <= 0 {
	//	panic("Listener must listen at least one event, or will not works")
	//}

	// add listener map
	listenerMap[l.Name()] = l

	// pair listener to event
	for _, event := range l.Subscribe() {
		registerListenerToEvent(event, l)
	}
}

func eventListener(e Eventer) []Listener {
	return hub[EventName(e)]
}

func registerListenerToEvent(e Eventer, l Listener) bool {
	// check if Listener has exist
	if isListenerRegisteredAtEvent(e, l) {
		return false
	}

	hub[EventName(e)] = append(hub[EventName(e)], l)

	return true
}

func isListenerRegisteredAtEvent(e Eventer, l Listener) bool {
	for _, Listener := range hub[EventName(e)] {
		if Listener.Name() == l.Name() {
			return true
		}
	}
	return false
}
