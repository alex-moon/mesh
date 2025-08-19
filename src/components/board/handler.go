package board

import (
	"context"
	"log/slog"
	"mesh/src/components/base"
	"mesh/src/components/column"
	"mesh/src/services"
	"net/http"

	"github.com/a-h/templ"
)

// Handler represents the board component's HTTP handler
type Handler struct {
	*base.BaseHandler
	CardService   *services.CardService
	ColumnHandler *column.Handler
}

// New creates a new board handler with dependencies
func New(log *slog.Logger, cardService *services.CardService, columnHandler *column.Handler) *Handler {
	return &Handler{
		BaseHandler:   base.NewBaseHandler(log),
		CardService:   cardService,
		ColumnHandler: columnHandler,
	}
}

// ServeHTTP handles HTTP requests for the board component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, h.Get, h.Post)
}

// Get renders the board component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isHTMX := h.IsHTMXRequest(r)

	// Load all columns from service
	columnsWithCards := h.CardService.GetColumns()

	// Convert to column components
	var columnComponents []templ.Component
	for _, columnWithCards := range columnsWithCards {
		columnComponent := h.ColumnHandler.RenderComponent(&columnWithCards)
		columnComponents = append(columnComponents, columnComponent)
	}

	// Create props for the board template
	props := BoardProps{
		Columns: columnComponents,
		IsHTMX:  isHTMX,
	}

	// Render using the base handler
	c := Board(props)
	h.RenderTemplate(ctx, w, c, "board")
}

// Post handles form submissions and actions
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	// Board-specific POST actions can be added here
	http.Error(w, "No actions implemented", http.StatusNotImplemented)
}

// RenderComponent provides a way for other components to render this one
func (h *Handler) RenderComponent(ctx context.Context) templ.Component {
	columnsWithCards := h.CardService.GetColumns()

	// Convert to column components
	var columnComponents []templ.Component
	for _, columnWithCards := range columnsWithCards {
		columnComponent := h.ColumnHandler.RenderComponent(&columnWithCards)
		columnComponents = append(columnComponents, columnComponent)
	}

	props := BoardProps{
		Columns: columnComponents,
		IsHTMX:  false,
	}

	return Board(props)
}
