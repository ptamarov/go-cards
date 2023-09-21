package algorithm

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// SM2Algorithm implements the NextCardAlgorithm interface. It determines a card status
// by using the "Super Memo 2" algorithm.
type sm2algorithm struct {
	judge judges.LevenshteinJudge
}

// SM2CardStatus holds the status of a card with the sm2 parameters:
//
//   - EaseFactor: the ease with which the user remembers the card.
//
//   - Interval: unitless distance to the next time the card can be shown.
//
//   - IsLearned: true if and only if the card is marked as learned.
//
//   - Progress: an integer in the range [0...5] that stores the user progress.
//
//   - TimesSeen: times the card was shown to the user.
type SM2CardStatus struct {
	EaseFactor float64 `json:"ease_factor"`
	Interval   int     `json:"interval"`
	IsLearned  bool    `json:"is_learned"`
	Progress   int     `json:"progress"`
	TimesSeen  int     `json:"times_seen"`
}

// NewSM2 returns a new SM2 Algorithm implementing the NextCardAlgorithm interface.
// Its judge is a LevenshteinJudge, which compares guesses using the Levenshtein distance.
func NewSM2(CaseInsensitive, UmlautInsensitive bool) sm2algorithm {
	var new sm2algorithm
	new.judge.CaseInsensitive = CaseInsensitive
	new.judge.UmlautInsensitive = UmlautInsensitive
	return new
}

// GetJudge returns a LevenshteinJudge.
func (sm2 *sm2algorithm) GetJudge() judges.Judge {
	return &sm2.judge
}

func (s SM2CardStatus) String() string {
	bytes, _ := json.MarshalIndent(s, " ", "\t")
	return fmt.Sprint(string(bytes))
}

// ComputeNewCardStatus computes a card status for a card from a list of actions.
func (sm2 *sm2algorithm) ComputeNewCardStatus(c card.MemoryCard, a []history.UserAction) CardStatus {
	var newCardStatus CardStatus
	status := sm2.DetermineStatus(c, a)
	newCardStatus.CardProgress = status.Progress
	newCardStatus.CardLearned = status.IsLearned

	// computing next available date needs more thought
	var delta time.Duration
	switch status.Interval {
	case 1:
		delta = 30 * time.Second
	case 6:
		delta = (60 * 5) * time.Second
	default:
		delta = time.Duration(status.Interval) * (30 * time.Minute)
	}
	newCardStatus.NextAvailableDate = time.Now().UTC().Add(delta)
	newCardStatus.TimesSeen = status.TimesSeen
	return newCardStatus

}

// ComputeActionQuality computes the quality of an action taken for a card. The quality is an integer from 0 to 5.
func (sm2 *sm2algorithm) ComputeActionQuality(card card.MemoryCard, action history.UserAction) int {
	correctFactor := sm2.judge.EvaluateUserAction(card, action)
	// 	4. After each repetition assess the quality of repetition response in 0-5 grade scale.
	return int(math.Floor(5 * correctFactor))
}

// ComputeNewEaseFactorFromQuality computes a new ease factor from an old one and a guess quality.
func (sm2 *sm2algorithm) ComputeNewEaseFactorFromQuality(oldEaseFactor float64, quality int) float64 {
	q := float64(quality)
	newEaseFactor := oldEaseFactor - 0.8 + 0.28*q - 0.02*q*q
	newEaseFactor = float64(int(newEaseFactor*100)) / 100
	if newEaseFactor < 1.3 {
		newEaseFactor = 1.3
	}
	// f(q) = - 0.8 + 0.28*q - 0.02*q*q is zero at q = 4.
	// It is negative in [0,4)
	// It is positive and increasing in (4,7)
	// It has a global maximum at q = 7.
	return newEaseFactor
}

// ComputeNextStatus computes the a SM2 status of a card from a previous status and a guess quality.
func (sm2 *sm2algorithm) ComputeNextStatus(guessQuality int, oldStatus SM2CardStatus) SM2CardStatus {
	// 3. Repeat items using the following intervals:
	// I(1):= 1, I(2):= 6
	// for n > 2 : I(n) = I(n-1) * EaseFactor
	// If interval is a fraction, round it up to the nearest integer.

	if oldStatus.IsLearned {
		fmt.Println("Card already learned!")
		return oldStatus
	}
	var newStatus SM2CardStatus

	// update progress if quality is 5 (perfect guess)
	if guessQuality == 5 {
		if oldStatus.Progress < 10 {
			newStatus.Progress = oldStatus.Progress + 1
		}
		if newStatus.Progress >= 10 {
			newStatus.IsLearned = true // update to learned if progress reaches 10
		}
	} else {
		newStatus.Progress = oldStatus.Progress
	}

	// update ease factor
	newStatus.EaseFactor = sm2.ComputeNewEaseFactorFromQuality(oldStatus.EaseFactor, guessQuality)

	// times seen goes up by one
	newStatus.TimesSeen = oldStatus.TimesSeen + 1

	// update interval
	switch newStatus.TimesSeen {
	case 0:
		newStatus.Interval = 1
	case 1:
		newStatus.Interval = 6
	default:
		switch guessQuality {
		case 5:
			newStatus.Interval = int(math.Ceil(float64(oldStatus.Interval) * (newStatus.EaseFactor - 0.8)))
		case 4:
			newStatus.Interval = int(math.Ceil(float64(oldStatus.Interval) * (newStatus.EaseFactor - 1.3)))
		default:
			newStatus.Interval = 1 // If quality is very low, start over.
			newStatus.EaseFactor = oldStatus.EaseFactor
		}
	}
	return newStatus
}

// DetermineStatus returns an SM2 card status from a list of guesses corresponding to a card.
func (sm2 *sm2algorithm) DetermineStatus(card card.MemoryCard, actions []history.UserAction) SM2CardStatus {
	var currentStatus SM2CardStatus
	currentStatus.EaseFactor = 2.5 // 2. With all items associate an EaseFactor equal to 2.5.

	for _, action := range actions {
		quality := sm2.ComputeActionQuality(card, action)
		log.Println(quality)
		currentStatus = sm2.ComputeNextStatus(quality, currentStatus)
		log.Println(currentStatus)
	}
	return currentStatus
}
