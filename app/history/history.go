package history

import (
	"time"

	"github.com/google/uuid"
)

// UserAction stores the action of a user.
type UserAction struct {
	UserID   uuid.UUID
	DeckID   uuid.UUID
	CardID   uuid.UUID
	Guess    string
	Duration float64
	Date     time.Time
}
