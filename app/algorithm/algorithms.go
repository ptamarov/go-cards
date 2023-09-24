package algorithm

import (
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// NextCardAlgorithm is an interface that computes the status of a card from a history of user actions.
type NextCardAlgorithm interface {
	ComputeNewCardStatus(card.MemoryCard, []history.UserAction) CardStatus
	GetJudge() judges.Judge
}

// CardStatus tracks the status of a card.
type CardStatus struct {
	NextAvailableDate time.Time // The next timestamp when the user can be shown the card.
	TimesSeen         int       // The numer of times the card has been shown to the user. Counts repetitions due to incorrect answers.
	CardLearned       bool      // True if and only if the user has mastered the card contents.
	CardProgress      int       // Must be an integer between -1 and 5. By convention, a card with progress -1 is inactive.
}
