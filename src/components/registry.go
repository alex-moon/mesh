// components/registry.go
package components

import (
	"log/slog"

	"mesh/src/components/app"
	"mesh/src/components/board"
	"mesh/src/components/card"
	"mesh/src/components/column"
	"mesh/src/services"
)

// Registry holds references to all component handlers
type Registry struct {
	AppHandler    *app.Handler
	BoardHandler  *board.Handler
	ColumnHandler *column.Handler
	CardHandler   *card.Handler
	CardService   *services.CardService
	EventService  *services.EventService
	SSEService    *services.SSEService
	WordService   *services.WordService
}

// NewRegistry creates a new registry with all handlers properly initialized
func NewRegistry(logger *slog.Logger) *Registry {
	// Create services
	eventService := services.NewEventService(logger)
	sseService := services.NewSSEService(logger)
	wordService, err := services.NewWordService(logger, "blacklist.txt")
	if err != nil {
		panic("Failed to create WordService: missing blacklist.txt")
		return nil
	}
	cardService := services.NewCardService(logger, eventService, wordService)

	// Create handlers with proper dependencies
	cardHandler := card.New(logger, eventService, cardService, wordService)
	columnHandler := column.New(logger, cardService, eventService, cardHandler, sseService)
	boardHandler := board.New(logger, eventService, cardService, columnHandler)
	appHandler := app.New(logger, eventService, boardHandler)

	return &Registry{
		AppHandler:    appHandler,
		BoardHandler:  boardHandler,
		ColumnHandler: columnHandler,
		CardHandler:   cardHandler,
		CardService:   cardService,
		EventService:  eventService,
		SSEService:    sseService,
		WordService:   wordService,
	}
}
