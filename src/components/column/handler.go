package column

import (
	"log/slog"
	"mesh/src/components/card"
	"mesh/src/services"
	"net/http"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

// Handler represents the column component's HTTP handler
type Handler struct {
	*base.BaseHandler
	CardHandler *card.Handler
	*services.CardService
	SSEService *services.SSEService
}

// New creates a new column handler with dependencies
func New(
	log *slog.Logger,
	cardService *services.CardService,
	eventService *services.EventService,
	cardHandler *card.Handler,
	sseService *services.SSEService,
) *Handler {
	h := &Handler{
		BaseHandler: base.NewBaseHandler(log, "column", eventService),
		CardHandler: cardHandler,
		CardService: cardService,
		SSEService:  sseService,
	}
	eventService.SubscribeCardChanged(h.OnCardChanged)
	eventService.SubscribeCardMoved(h.OnCardMoved)
	eventService.SubscribeCardDeleted(h.OnCardDeleted)
	return h
}

func (h *Handler) OnCardDeleted(event *services.CardDeletedEvent, context services.EventContext) {
	// Broadcast to column updates via SSE for real-time collaboration
	column, err := h.CardService.GetColumn(event.ColumnID)
	if err == nil {
		component := h.RenderComponent(column, true)
		h.SSEService.BroadcastOOBUpdate(component)
	} else {
		h.Log.Error("Failed to get to-column for SSE broadcast", "columnID", event.ColumnID, "error", err)
	}
}

func (h *Handler) OnCardChanged(event *services.CardChangedEvent, context services.EventContext) {
	card, err := h.CardService.GetCard(event.CardID)
	if err != nil {
		h.Log.Error("Failed to get card for card changed event", "cardID", event.CardID, "error", err)
		return
	}

	// Broadcast to column updates via SSE for real-time collaboration
	column, err := h.CardService.GetColumn(card.ColumnID)
	if err == nil {
		component := h.RenderComponent(column, true)
		h.SSEService.BroadcastOOBUpdate(component)
	} else {
		h.Log.Error("Failed to get to-column for SSE broadcast", "columnID", card.ColumnID, "error", err)
	}
}

func (h *Handler) OnCardMoved(event *services.CardMovedEvent, context services.EventContext) {
	// Broadcast to column updates via SSE for real-time collaboration
	column, err := h.CardService.GetColumn(event.ToColumnID)
	if err == nil {
		component := h.RenderComponent(column, true)
		h.SSEService.BroadcastOOBUpdate(component)
	} else {
		h.Log.Error("Failed to get to-column for SSE broadcast", "columnID", event.ToColumnID, "error", err)
	}

	column, err = h.CardService.GetColumn(event.FromColumnID)
	if err == nil {
		h.SSEService.BroadcastOOBUpdate(h.RenderComponent(column, true))
	} else {
		h.Log.Error("Failed to get from-column for SSE broadcast", "columnID", event.FromColumnID, "error", err)
	}
}

// ServeHTTP handles HTTP requests for the column component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, map[string]http.HandlerFunc{
		http.MethodGet: h.Get,
	})
}

// Get renders the full column component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Convert to column components
	var columnID = ctx.Value("columnID").(int)
	var columnWithCards, err = h.CardService.GetColumn(columnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}
	var cardComponents []templ.Component
	for _, card := range columnWithCards.Cards {
		cardComponent := h.CardHandler.RenderComponent(&card)
		cardComponents = append(cardComponents, cardComponent)
	}
	newCard := h.CardHandler.RenderComponentForNew(columnWithCards.Column.ID)
	cardComponents = append(cardComponents, newCard)
	props := ColumnProps{
		Column: &columnWithCards.Column,
		Cards:  cardComponents,
	}

	// Render using the base handler
	c := Column(props)
	h.RenderTemplate(ctx, w, c)
}

func (h *Handler) RenderComponent(column *services.ColumnWithCards, oob bool) templ.Component {
	// Convert to card components
	var cardComponents []templ.Component
	for _, card := range column.Cards {
		columnComponent := h.CardHandler.RenderComponent(&card)
		cardComponents = append(cardComponents, columnComponent)
	}
	newCard := h.CardHandler.RenderComponentForNew(column.Column.ID)
	cardComponents = append(cardComponents, newCard)
	props := ColumnProps{
		Column: &column.Column,
		Cards:  cardComponents,
		OOB:    oob,
	}

	return Column(props)
}
