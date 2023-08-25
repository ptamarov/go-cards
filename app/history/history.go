package history

import (
	"time"

	"github.com/google/uuid"
)

type Guess struct {
	Guess    string
	Duration float32
}

type Action struct {
	MemoryCardID uuid.UUID
	Guess        Guess
	Date         time.Time // include 
}

type UserHistoryForDeck struct {
	UserID         uuid.UUID
	DeckID         uuid.UUID
	HistoryForDeck []Action
}
