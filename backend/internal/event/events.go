package event

import (
	"time"
)

type EventType string

type Event interface {
	Type() EventType
	Timestamp() time.Time
	AggregateID() string
}

type BaseEvent struct {
	EventType EventType
	CreatedAt time.Time
	ID        string
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}

func (e BaseEvent) AggregateID() string {
	return e.ID
}

func (e BaseEvent) Timestamp() time.Time {
	return e.CreatedAt
}
