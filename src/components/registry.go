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
}

// NewRegistry creates a new registry with all handlers properly initialized
func NewRegistry(logger *slog.Logger) *Registry {
	// Create service
	cardService := services.NewCardService(logger)

	// Create handlers with proper dependencies
	cardHandler := card.New(logger, cardService)
	columnHandler := column.New(logger, cardService, cardHandler)
	boardHandler := board.New(logger, cardService, columnHandler)
	appHandler := app.New(logger, boardHandler)

	return &Registry{
		AppHandler:    appHandler,
		BoardHandler:  boardHandler,
		ColumnHandler: columnHandler,
		CardHandler:   cardHandler,
		CardService:   cardService,
	}
}
