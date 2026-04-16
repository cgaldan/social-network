package event

import (
	"sync"

	"social-network/packages/logger"
)

type EventHandler func(Event) error

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType EventType, handler EventHandler) error
	Unsubscribe(eventType EventType, handler EventHandler) error
}

type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
	logger   *logger.Logger
}

func NewInMemoryBus(logger *logger.Logger) *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[EventType][]EventHandler),
		logger:   logger,
	}
}

func (bus *InMemoryBus) Publish(event Event) error {
	bus.mu.RLock()
	handlers, exists := bus.handlers[event.Type()]
	bus.mu.RUnlock()

	if !exists {
		bus.logger.Debug("No handlers for event", "eventType", event.Type())
		return nil
	}

	for _, handler := range handlers {
		go func(h EventHandler, e Event) {
			if err := h(e); err != nil {
				bus.logger.Error("Error handling event", "eventType", e.Type(), "error", err)
			}
		}(handler, event)
	}

	return nil
}

func (bus *InMemoryBus) Subscribe(eventType EventType, handler EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	bus.logger.Debug("Handler subscribed to event", "eventType", eventType)
	return nil
}
