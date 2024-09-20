package bus

type EventBus struct {
	Events chan []byte
	Ready  chan bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		Events: make(chan []byte),
		Ready:  make(chan bool, 10),
	}
}

func (bus *EventBus) Publish(event []byte) {
	bus.Events <- event
}

func (bus *EventBus) Close() {
	close(bus.Events)
	close(bus.Ready)
}
