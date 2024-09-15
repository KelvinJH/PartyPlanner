package bus

type EventBus struct {
	Events chan []byte
}

func NewEventBus() *EventBus {
	return &EventBus{
		Events: make(chan []byte),
	}
}

func (bus *EventBus) Publish(event []byte) {
	bus.Events <- event
}

func (bus *EventBus) Close() {
	close(bus.Events)
}