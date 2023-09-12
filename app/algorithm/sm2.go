package algorithm

import (
	"math"
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

type SM2Algorithm struct {
	judges.LevenshsteinJudge
}

type SM2CardStatus struct {
	easeFactor float64
	interval   int
	learned    bool
	progress   int
	timesSeen  int
}

func (sm2 *SM2Algorithm) ComputeNewCardStatus(c card.MemoryCard, a []history.UserAction, j judges.Judge, os card.CardStatus) card.CardStatus {
	var newCardStatus card.CardStatus
	status := DetermineStatus(j, c, a)
	newCardStatus.CardProgress = status.progress
	newCardStatus.CardLearned = status.learned

	delta := time.Duration(status.interval) * time.Minute
	newCardStatus.NextAvailableDate = time.Now().Add(delta)

	return newCardStatus

}

func ComputeActionQuality(j judges.Judge, card card.MemoryCard, action history.UserAction) int {
	correctFactor := j.EvaluateUserAction(card, action)
	// 	4. After each repetition assess the quality of repetition response in 0-5 grade scale.
	if action.Duration >= 10 {
		return 0
	} else {
		return int(math.Floor(5 * correctFactor))
	}
}

func ComputeNextStatus(guessQuality int, oldStatus SM2CardStatus) SM2CardStatus {
	// 3. Repeat items using the following intervals:
	// I(1):= 1
	// I(2):= 6
	// for n>2 : I(n) = I(n-1)*EaseFactor
	// If interval is a fraction, round it up to the nearest integer.

	var newStatus SM2CardStatus

	// update times seen
	newStatus.timesSeen = oldStatus.timesSeen + 1

	// update progress if quality is 5
	if guessQuality == 5 {
		newStatus.progress = oldStatus.progress + 1
		if newStatus.progress >= 5 {
			// update learned if progress reaches 5
			newStatus.learned = true
		}
	}

	// update ease factor
	q := float64(guessQuality)

	newStatus.easeFactor = oldStatus.easeFactor - 0.8 + 0.28*q - 0.02*q*q
	// f(q) = - 0.8 + 0.28*q - 0.02*q*q is zero at q = 4.
	// It is negative in [0,4)
	// It is positive and increasing in (4,7)
	// It has a global maximum at q = 7.

	if newStatus.easeFactor < 1.3 {
		newStatus.easeFactor = 1.3
	}

	// update interval
	switch newStatus.timesSeen {
	case 0:
		newStatus.interval = 1
	case 1:
		newStatus.interval = 6
	default:
		switch guessQuality >= 3 {
		case false:
			newStatus.interval = 1
			newStatus.easeFactor = oldStatus.easeFactor
		default:
			newStatus.interval = int(math.Ceil(float64(oldStatus.interval) * newStatus.easeFactor))
		}
	}
	return newStatus
}

// Processes a list of guesses corresponding to a card, ouputs the interval
// value obtained by processing the data according to the SM-2 algorithm."""
func DetermineStatus(judge judges.Judge, card card.MemoryCard, actions []history.UserAction) SM2CardStatus {
	var currentStatus SM2CardStatus
	currentStatus.easeFactor = 2.5 // 2. With all items associate an EaseFactor equal to 2.5.

	for _, action := range actions {
		quality := ComputeActionQuality(judge, card, action)
		currentStatus = ComputeNextStatus(quality, currentStatus)
	}
	return currentStatus
}
