package algorithm

import (
	"math"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/deck"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/database"
)

// NextCardAlgorithm is an interface that can give the user a next card to study.
type NextCardAlgorithm interface {
	// GetNextCard must always pop a card from deck to avoid duplication.
	GetNextCard(GetNextCardInput) card.MemoryCard
}

type GetNextCardInput struct {
	UserID   uuid.UUID
	Deck     *deck.MemoryCardDeck
	Judge    GuessJudge
	Database database.Database
}

type NaiveAlgorithm struct {
}

func (na NaiveAlgorithm) GetNextCard(input GetNextCardInput) card.MemoryCard {
	return input.Deck.PopLast()
}

type SM2Algorithm struct {
}

func (sm2 SM2Algorithm) GetNextCard(in GetNextCardInput) card.MemoryCard {
	var topPriority card.MemoryCard
	smallestInterval := math.Inf(-1)

	userHistory := (in.Database).GetUserHistoryForDeck(in.UserID, in.Deck.DeckID).HistoryForDeck

	for _, card := range in.Deck.Cards {
		var cardHistory []history.Guess // store the history of guesses for a given card

		for _, action := range userHistory {
			if action.MemoryCardID == card.CardID {
				cardHistory = append(cardHistory, action.Guess)
			}
		}

		cardInterval := determineCardStatusForSM2(in.Judge, card, cardHistory)
		if cardInterval <= int(smallestInterval) {
			smallestInterval = float64(cardInterval)
			topPriority = card
		}
	}
	// TODO: add logic that removes topPriority from deck.
	return topPriority
}

func computeQualityForSM2(j GuessJudge, card card.MemoryCard, guess history.Guess) int {
	correctFactor := j.EvaluateGuess(card, guess)
	var quality int

	if guess.Duration >= 10 {
		quality = 0
	} else {
		quality = int(math.Floor(5 * correctFactor))
	}

	return quality
}

func computeNewIntervalForSM2(newEaseFactor float64, pastInterval int, timesSeen int) int {
	switch timesSeen {
	case 0:
		return 1
	case 1:
		return 6
	default:
		return int(math.Floor(float64(pastInterval) * newEaseFactor))
	}
}

// Processes a list of guesses corresponding to a card, ouputs the interval
// value obtained by processing the data according to the SM-2 algorithm."""
func determineCardStatusForSM2(j GuessJudge, card card.MemoryCard, guesses []history.Guess) int {
	var interval int
	var easeFactor float64
	var timesSeen int

	interval = 0
	easeFactor = 2.5
	timesSeen = 0

	for _, guess := range guesses {
		quality := computeQualityForSM2(j, card, guess)

		switch quality {
		case 0:
			interval = 1
		case 1:
			interval = 6
		default:
			interval = computeNewIntervalForSM2(easeFactor, interval, timesSeen)

		}

	}
	return interval
}
