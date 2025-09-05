package services

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
)

type Card struct {
	ID       int
	Title    string
	Content  string
	ColumnID int
}

type Column struct {
	ID    int
	Title string
	Order int
}

type CardService struct {
	mu          sync.RWMutex
	cards       map[int]*Card   // cardID -> Card
	columns     map[int]*Column // columnID -> Column
	columnCards map[int][]int   // columnID -> []cardID (ordered)

	nextCardID   int
	nextColumnID int

	log          *slog.Logger
	eventService *EventService
}

func NewCardService(log *slog.Logger, eventService *EventService) *CardService {
	service := &CardService{
		mu:           sync.RWMutex{},
		cards:        make(map[int]*Card),
		columns:      make(map[int]*Column),
		columnCards:  make(map[int][]int),
		log:          log,
		eventService: eventService,
	}
	service.seedData()
	return service
}

func (c *CardService) seedData() {
	// Create columns
	c.columns[1] = &Column{ID: 1, Title: "To Do", Order: 0}
	c.columns[2] = &Column{ID: 2, Title: "In Progress", Order: 1}
	c.columns[3] = &Column{ID: 3, Title: "Done", Order: 2}
	c.nextColumnID = 4

	// Create cards
	c.cards[1] = &Card{ID: 1, Title: "Blog post", Content: "Once the app is working and looking good, write it up", ColumnID: 1}
	c.cards[2] = &Card{ID: 2, Title: "Post to HN", Content: "", ColumnID: 1}
	c.cards[3] = &Card{ID: 3, Title: "Build app", Content: "Implement minimal Kanban Board with columns and draggable/editable cards", ColumnID: 2}
	c.nextCardID = 4

	// Set up column ordering
	c.columnCards[1] = []int{1, 2}
	c.columnCards[2] = []int{3}
	c.columnCards[3] = []int{}
}

func removeFromSlice(slice []int, element int) []int {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (c *CardService) getColumnByOrder(order int) *Column {
	for _, column := range c.columns {
		if column.Order == order {
			return column
		}
	}
	return nil
}

func (c *CardService) getSortedColumns() []*Column {
	columns := make([]*Column, 0, len(c.columns))
	for _, column := range c.columns {
		columns = append(columns, column)
	}
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Order < columns[j].Order
	})
	return columns
}

func (c *CardService) CanPromote(cardID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	card, exists := c.cards[cardID]
	if !exists {
		return false
	}

	currentColumn := c.columns[card.ColumnID]
	if currentColumn == nil {
		return false
	}

	nextColumn := c.getColumnByOrder(currentColumn.Order + 1)
	return nextColumn != nil
}

func (c *CardService) CanDemote(cardID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	card, exists := c.cards[cardID]
	if !exists {
		return false
	}

	currentColumn := c.columns[card.ColumnID]
	if currentColumn == nil {
		return false
	}

	return currentColumn.Order > 0
}

func (c *CardService) Promote(cardID int) (*Column, *Column, error) {
	card, exists := c.cards[cardID]
	if !exists {
		return nil, nil, fmt.Errorf("card with ID %d not found", cardID)
	}

	currentColumn := c.columns[card.ColumnID]
	if currentColumn == nil {
		return nil, nil, fmt.Errorf("current column not found for card %d", cardID)
	}

	targetColumn := c.getColumnByOrder(currentColumn.Order + 1)
	if targetColumn == nil {
		return nil, nil, fmt.Errorf("card %d cannot be promoted further", cardID)
	}

	return c.MoveCard(cardID, targetColumn.ID, -1)
}

func (c *CardService) Demote(cardID int) (*Column, *Column, error) {
	card, exists := c.cards[cardID]
	if !exists {
		return nil, nil, fmt.Errorf("card with ID %d not found", cardID)
	}

	currentColumn := c.columns[card.ColumnID]
	if currentColumn == nil {
		return nil, nil, fmt.Errorf("current column not found for card %d", cardID)
	}

	targetColumn := c.getColumnByOrder(currentColumn.Order - 1)
	if targetColumn == nil {
		return nil, nil, fmt.Errorf("card %d cannot be demoted further", cardID)
	}

	return c.MoveCard(cardID, targetColumn.ID, -1)
}

func (c *CardService) GetColumn(id int) (*ColumnWithCards, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	column, exists := c.columns[id]
	if !exists {
		return nil, fmt.Errorf("column with id %d not found", id)
	}

	cards := c.getCardsForColumn(id)
	return &ColumnWithCards{
		Column: *column,
		Cards:  cards,
	}, nil
}

func (c *CardService) GetColumns() []ColumnWithCards {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []ColumnWithCards

	// Sort columns by order
	for i := 0; i < len(c.columns); i++ {
		for _, column := range c.columns {
			if column.Order == i {
				cards := c.getCardsForColumn(column.ID)
				result = append(result, ColumnWithCards{
					Column: *column,
					Cards:  cards,
				})
				break
			}
		}
	}

	return result
}

type ColumnWithCards struct {
	Column Column
	Cards  []Card
}

func (c *CardService) getCardsForColumn(columnID int) []Card {
	cardIDs := c.columnCards[columnID]
	cards := make([]Card, 0, len(cardIDs))

	for _, cardID := range cardIDs {
		if card, exists := c.cards[cardID]; exists {
			cards = append(cards, *card)
		}
	}

	return cards
}

func (c *CardService) GetCard(cardID int) (*Card, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if card, exists := c.cards[cardID]; exists {
		return card, nil
	}
	return nil, fmt.Errorf("card with ID %d not found", cardID)
}

func (c *CardService) AddCard(title, content string, columnID int) (*Card, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.columns[columnID]; !exists {
		return nil, fmt.Errorf("column with ID %d not found", columnID)
	}

	card := &Card{
		ID:       c.nextCardID,
		Title:    title,
		Content:  content,
		ColumnID: columnID,
	}

	c.cards[c.nextCardID] = card
	c.columnCards[columnID] = append(c.columnCards[columnID], c.nextCardID)
	c.nextCardID++

	return card, nil
}

func (c *CardService) UpdateCard(cardID int, title, content string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	card.Title = title
	card.Content = content
	return nil
}

func (c *CardService) MoveCard(cardID, newColumnID, newPosition int) (*Column, *Column, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return nil, nil, fmt.Errorf("card with ID %d not found", cardID)
	}

	if _, exists := c.columns[newColumnID]; !exists {
		return nil, nil, fmt.Errorf("column with ID %d not found", newColumnID)
	}

	oldColumnID := card.ColumnID
	c.removeCardFromColumn(cardID, oldColumnID)
	c.insertCardInColumn(cardID, newColumnID, newPosition)
	card.ColumnID = newColumnID

	return c.columns[oldColumnID], c.columns[newColumnID], nil
}

func (c *CardService) DeleteCard(cardID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	c.removeCardFromColumn(cardID, card.ColumnID)
	delete(c.cards, cardID)
	return nil
}

func (c *CardService) removeCardFromColumn(cardID, columnID int) {
	c.columnCards[columnID] = removeFromSlice(c.columnCards[columnID], cardID)
}

func (c *CardService) insertCardInColumn(cardID, columnID, position int) {
	cardList := c.columnCards[columnID]

	if position == -1 || position >= len(cardList) {
		c.columnCards[columnID] = append(cardList, cardID)
	} else {
		tail := append([]int{cardID}, cardList[position:]...)
		cardList = append(cardList[:position], tail...)
		c.columnCards[columnID] = cardList
	}
}
