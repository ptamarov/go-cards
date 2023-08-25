package application

import (
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/deck"
	"github.com/ptamarov/go-cards/app/user"
)

type ApplicationConfig struct {
	User      user.User
	Deck      deck.MemoryCardDeck
	Algorithm algorithm.NextCardAlgorithm
	Judge     algorithm.GuessJudge
}
