package services

import (
	"log/slog"

	"github.com/a-h/templ"
)

const (
	CardMovedEventKey = "card-moved"
)

type Event interface {
	Key() string
}

type CardMovedEvent struct {
	CardID       int
	FromColumnID int
	ToColumnID   int
	oob          []templ.Component
}

func (e *CardMovedEvent) Key() string {
	return CardMovedEventKey
}

func (e *CardMovedEvent) UpdateOob(component templ.Component) []templ.Component {
	e.oob = append(e.oob, component)
	return e.oob
}

func (e *CardMovedEvent) GetOob() []templ.Component {
	return e.oob
}

type EventService struct {
	log         *slog.Logger
	subscribers map[string][]func(event Event)
}

func (e *EventService) Publish(event Event) {
	for _, subscriber := range e.subscribers[event.Key()] {
		subscriber(event)
	}
}

func (e *EventService) Subscribe(key string, subscriber func(event Event)) {
	e.subscribers[key] = append(e.subscribers[key], subscriber)
}

func NewEventService(log *slog.Logger) *EventService {
	return &EventService{
		log:         log,
		subscribers: make(map[string][]func(event Event)),
	}
}

func (e *EventService) PublishCardMoved(cardID int, fromColumnID int, toColumnID int) *CardMovedEvent {
	event := &CardMovedEvent{
		CardID:       cardID,
		FromColumnID: fromColumnID,
		ToColumnID:   toColumnID,
	}
	e.Publish(event)
	return event
}

func (e *EventService) SubscribeCardMoved(subscriber func(event *CardMovedEvent)) {
	e.Subscribe(CardMovedEventKey, func(event Event) {
		subscriber(event.(*CardMovedEvent))
	})
}
