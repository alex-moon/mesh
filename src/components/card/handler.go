package card

import (
	"fmt"
	"log/slog"
	"mesh/src/services"
	"net/http"
	"strconv"
	"strings"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

type Handler struct {
	*base.BaseHandler
	*services.CardService
}

func New(log *slog.Logger, eventService *services.EventService, cardService *services.CardService) *Handler {
	return &Handler{
		BaseHandler: base.NewBaseHandler(log, "card", eventService),
		CardService: cardService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, map[string]http.HandlerFunc{
		http.MethodGet:    h.Get,
		http.MethodPost:   h.Post,
		http.MethodPatch:  h.Patch,
		http.MethodPut:    h.Put,
		http.MethodDelete: h.Delete,
	})
}

func (h *Handler) getCardFromRequest(r *http.Request) (*services.Card, error) {
	cardIDString := r.FormValue("cardID")
	if cardIDString == "" {
		return nil, fmt.Errorf("missing card ID")
	}

	cardID, err := strconv.Atoi(cardIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid card ID %s", cardIDString)
	}

	card, err := h.CardService.GetCard(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found %d", cardID)
	}

	return card, nil
}

func (h *Handler) getColumnFromRequest(r *http.Request) (*services.Column, error) {
	columnIDString := r.FormValue("columnID")
	if columnIDString == "" {
		return nil, fmt.Errorf("missing column ID")
	}

	columnID, err := strconv.Atoi(columnIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid column ID %s", columnIDString)
	}

	column, err := h.CardService.GetColumn(columnID)
	if err != nil {
		return nil, fmt.Errorf("column not found %d", columnID)
	}

	return &column.Column, nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.RenderTemplate(r.Context(), w, h.RenderComponent(card))
}

func (h *Handler) validate(r *http.Request) (Data, Errors) {
	errors := Errors{}

	var data = Data{
		Title:   strings.TrimSpace(r.FormValue("title")),
		Content: strings.TrimSpace(r.FormValue("content")),
	}

	if data.Title == "" {
		errors.Title = "Title is required"
	}

	if len(data.Title) > 100 {
		errors.Title = "Title must be less than 100 characters"
	}

	if len(data.Content) > 1000 {
		errors.Content = "Content must be less than 1000 characters"
	}

	if r.FormValue("columnID") != "" {
		var column, err = h.getColumnFromRequest(r)
		if err != nil {
			errors.ColumnID = err.Error()
		} else {
			data.ColumnID = column.ID
		}
	}

	return data, errors
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.CardService.DeleteCard(card.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.EventService.PublishCardDeleted(card.ColumnID)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var data, errors = h.validate(r)
	if errors.Any() {
		var props = h.getPropsWithData(&services.Card{}, data, errors)
		h.RenderTemplate(r.Context(), w, Card(props))
		return
	}

	card, err := h.CardService.AddCard(
		data.Title,
		data.Content,
		data.ColumnID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.RenderTemplate(r.Context(), w, h.RenderComponent(card))
	h.RenderTemplate(r.Context(), w, h.RenderComponentForNew(card.ColumnID))

	h.EventService.PublishCardChanged(card.ID)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var data, errors = h.validate(r)

	if errors.Any() {
		var props = h.getPropsWithData(card, data, errors)
		h.RenderTemplate(r.Context(), w, Card(props))
		return
	}

	err = h.CardService.UpdateCard(
		card.ID,
		data.Title,
		data.Content,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.RenderTemplate(r.Context(), w, h.RenderComponent(card))

	h.EventService.PublishCardChanged(card.ID)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	action := r.FormValue("action")
	switch action {
	case PutActionDemote:
		fromColumn, toColumn, err := h.CardService.Demote(card.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updatedCard, err := h.CardService.GetCard(card.ID)
		if err == nil {
			props := h.getProps(updatedCard)
			props.OOB = true
			h.RenderTemplate(r.Context(), w, Card(props))
		}
		h.EventService.PublishCardMoved(card.ID, fromColumn.ID, toColumn.ID)
		break
	case PutActionPromote:
		fromColumn, toColumn, err := h.CardService.Promote(card.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Get the updated card after the move
		updatedCard, err := h.CardService.GetCard(card.ID)
		if err == nil {
			// Render the moved card with OOB to provide immediate visual feedback
			props := h.getProps(updatedCard)
			props.OOB = true
			h.RenderTemplate(r.Context(), w, Card(props))
		}
		h.EventService.PublishCardMoved(card.ID, fromColumn.ID, toColumn.ID)
		break
	case PutActionMove:
		columnID, err := strconv.Atoi(r.FormValue("columnID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		position, err := strconv.Atoi(r.FormValue("position"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fromColumn, toColumn, err := h.CardService.MoveCard(card.ID, columnID, position)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Get the updated card after the move
		updatedCard, err := h.CardService.GetCard(card.ID)
		if err == nil {
			// Render the moved card with OOB to provide immediate visual feedback
			props := h.getProps(updatedCard)
			props.OOB = true
			h.RenderTemplate(r.Context(), w, Card(props))
		}
		h.EventService.PublishCardMoved(card.ID, fromColumn.ID, toColumn.ID)
	}
}

func (h *Handler) RenderComponent(card *services.Card) templ.Component {
	props := h.getProps(card)
	return Card(props)
}

func (h *Handler) RenderComponentForNew(columnID int) templ.Component {
	props := h.getPropsForNew(columnID)
	return Card(props)
}

func (h *Handler) getPropsForNew(columnID int) CardProps {
	return h.getPropsWithData(
		&services.Card{ColumnID: columnID},
		Data{},
		Errors{},
	)
}

func (h *Handler) getProps(card *services.Card) CardProps {
	return h.getPropsWithData(card, Data{
		ID:      card.ID,
		Title:   card.Title,
		Content: card.Content,
	}, Errors{})
}

func (h *Handler) getPropsWithData(card *services.Card, data Data, errors Errors) CardProps {
	return CardProps{
		Card:       card,
		Data:       data,
		Errors:     errors,
		IsEditing:  errors.Any(),
		CanDemote:  h.CardService.CanDemote(card.ID),
		CanPromote: h.CardService.CanPromote(card.ID),
	}
}
