package history

import (
	"time"

	"github.com/google/uuid"
)

type UserAction struct {
	UserID   uuid.UUID
	DeckID   uuid.UUID
	CardID   uuid.UUID
	Guess    string
	Duration float64
	Date     time.Time // include
}

type UserHistoryForDeck struct {
	UserID         uuid.UUID
	DeckID         uuid.UUID
	HistoryForDeck []UserAction
}
