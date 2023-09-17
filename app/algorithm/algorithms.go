package algorithm

import (
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

// NextCardAlgorithm is an interface that computes the
// status of a card from a history of user actions.
type NextCardAlgorithm interface {
	ComputeNewCardStatus(card.MemoryCard, []history.UserAction) CardStatus
}

// CardStatus tracks the status of a card.
type CardStatus struct {
	NextAvailableDate time.Time
	TimesSeen         int
	CardLearned       bool
	CardProgress      int // integer between -1 and 5, -1: inactive
}
