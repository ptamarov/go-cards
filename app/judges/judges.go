package judges

import (
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

// Judge is an interface that can evaluate an action on a card. Should return a float in the range [0..1].
type Judge interface {
	EvaluateUserAction(card card.MemoryCard, h history.UserAction) float64
}
