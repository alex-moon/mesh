package column

import (
	"log/slog"
	"mesh/src/components/base"
	"mesh/src/components/card"
	"mesh/src/services"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type Handler struct {
	*base.BaseHandler
	CardHandler *card.Handler
	*services.CardService
	SSEService *services.SSEService
}

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
	eventService.SubscribeCardDeleted(h.OnCardDeleted)
	eventService.SubscribeCardChanged(h.OnCardChanged)
	eventService.SubscribeCardMoved(h.OnCardMoved)
	return h
}

func (h *Handler) OnCardDeleted(event *services.CardDeletedEvent) {
	column, err := h.CardService.GetColumn(event.ColumnID)
	if err == nil {
		component := h.RenderComponent(column, true)
		h.SSEService.BroadcastOOBUpdate(component)
	} else {
		h.Log.Error("Failed to get to-column for SSE broadcast", "columnID", event.ColumnID, "error", err)
	}
}

func (h *Handler) OnCardChanged(event *services.CardChangedEvent) {
	card, err := h.CardService.GetCard(event.CardID)
	if err != nil {
		h.Log.Error("Failed to get card for card changed event", "cardID", event.CardID, "error", err)
		return
	}

	column, err := h.CardService.GetColumn(card.ColumnID)
	if err == nil {
		component := h.RenderComponent(column, true)
		h.SSEService.BroadcastOOBUpdate(component)
	} else {
		h.Log.Error("Failed to get to-column for SSE broadcast", "columnID", card.ColumnID, "error", err)
	}
}

func (h *Handler) OnCardMoved(event *services.CardMovedEvent) {
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, map[string]http.HandlerFunc{
		http.MethodGet: h.Get,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	columnIDString := r.FormValue("columnID")
	if columnIDString == "" {
		http.Error(w, "Missing column ID", http.StatusNotFound)
		return
	}

	columnID, err := strconv.Atoi(columnIDString)
	if err != nil {
		http.Error(w, "Invalid column ID", http.StatusNotFound)
		return
	}

	column, err := h.CardService.GetColumn(columnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	h.RenderTemplate(r.Context(), w, h.RenderComponent(column, false))
}

func (h *Handler) RenderComponent(column *services.ColumnWithCards, oob bool) templ.Component {
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
