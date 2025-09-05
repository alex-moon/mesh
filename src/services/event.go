package services

import (
	"log/slog"
)

const (
	CardMovedEventKey   = "card-moved"
	CardChangedEventKey = "card-changed"
	CardDeletedEventKey = "card-deleted"
)

type Event interface {
	Key() string
}

type CardDeletedEvent struct {
	ColumnID int
}

func (e *CardDeletedEvent) Key() string {
	return CardDeletedEventKey
}

type CardChangedEvent struct {
	CardID int
}

func (e *CardChangedEvent) Key() string {
	return CardChangedEventKey
}

type CardMovedEvent struct {
	CardID       int
	FromColumnID int
	ToColumnID   int
}

func (e *CardMovedEvent) Key() string {
	return CardMovedEventKey
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

func (e *EventService) PublishCardMoved(
	cardID int,
	fromColumnID int,
	toColumnID int,
) *CardMovedEvent {
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

func (e *EventService) PublishCardChanged(cardID int) *CardChangedEvent {
	event := &CardChangedEvent{
		CardID: cardID,
	}
	e.Publish(event)
	return event
}

func (e *EventService) SubscribeCardChanged(subscriber func(event *CardChangedEvent)) {
	e.Subscribe(CardChangedEventKey, func(event Event) {
		subscriber(event.(*CardChangedEvent))
	})
}

func (e *EventService) PublishCardDeleted(columnID int) *CardDeletedEvent {
	event := &CardDeletedEvent{
		ColumnID: columnID,
	}
	e.Publish(event)
	return event
}

func (e *EventService) SubscribeCardDeleted(subscriber func(event *CardDeletedEvent)) {
	e.Subscribe(CardDeletedEventKey, func(event Event) {
		subscriber(event.(*CardDeletedEvent))
	})
}
