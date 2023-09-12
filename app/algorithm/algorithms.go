package algorithm

import (
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// NextCardAlgorithm computes an updated card status from list of user actions
// on a card and a judge of correctness.
type NextCardAlgorithm interface {
	ComputeNewCardStatus(card.MemoryCard, []history.UserAction, judges.Judge) card.CardStatus
}
