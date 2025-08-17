package services

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

const (
	CardMovedEventKey = "card-moved"
)

type Event interface {
	Key() string
}

type EventContext struct {
	Context        context.Context
	ResponseWriter http.ResponseWriter
}

func (e *EventContext) Write(component templ.Component) {
	err := component.Render(e.Context, e.ResponseWriter)
	if err != nil {
		http.Error(e.ResponseWriter, "Failed to render OOB updates", http.StatusInternalServerError)
	}
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
	subscribers map[string][]func(event Event, context EventContext)
}

func (e *EventService) Publish(event Event, w http.ResponseWriter, ctx context.Context) {
	eventContext := EventContext{
		Context:        ctx,
		ResponseWriter: w,
	}
	for _, subscriber := range e.subscribers[event.Key()] {
		subscriber(event, eventContext)
	}
}

func (e *EventService) Subscribe(key string, subscriber func(event Event, context EventContext)) {
	e.subscribers[key] = append(e.subscribers[key], subscriber)
}

func NewEventService(log *slog.Logger) *EventService {
	return &EventService{
		log:         log,
		subscribers: make(map[string][]func(event Event, context EventContext)),
	}
}

func (e *EventService) PublishCardMoved(
	cardID int,
	fromColumnID int,
	toColumnID int,
	w http.ResponseWriter,
	ctx context.Context,
) *CardMovedEvent {
	event := &CardMovedEvent{
		CardID:       cardID,
		FromColumnID: fromColumnID,
		ToColumnID:   toColumnID,
	}
	e.Publish(event, w, ctx)
	return event
}

func (e *EventService) SubscribeCardMoved(subscriber func(event *CardMovedEvent, context EventContext)) {
	e.Subscribe(CardMovedEventKey, func(event Event, context EventContext) {
		subscriber(event.(*CardMovedEvent), context)
	})
}
