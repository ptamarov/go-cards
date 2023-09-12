package algorithm

import (
	"fmt"
	"log"
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// NaiveAlgorithm computes a new card status by only looking at the last action
// performed by the user.
type NaiveAlgorithm struct {
	judges.LevenshsteinJudge
}

func (na *NaiveAlgorithm) ComputeNewCardStatus(
	memCard card.MemoryCard,
	actions []history.UserAction,
	judge judges.Judge,
	oldStatus card.CardStatus,
) card.CardStatus {

	oneMinute := 60 * time.Second
	oneHour := 60 * oneMinute
	oneDay := 24 * oneHour

	if len(actions) == 0 {
		fmt.Println("ALGORITHM: no actions to evaluate")
		return oldStatus
	}
	action := actions[len(actions)-1]
	score := na.EvaluateUserAction(memCard, action)

	log.Println("score is", score)
	var newCardStatus card.CardStatus

	timeNowUTC := time.Now().UTC()
	if score == 1.0 {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(oneDay)
		newCardStatus.CardProgress = oldStatus.CardProgress + 1
		if newCardStatus.CardProgress == 5 {
			newCardStatus.CardLearned = true
		}
		return newCardStatus
	} else if score >= 0.8 {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(oneHour)
	} else if score >= 0.5 {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(5 * oneMinute)
		newCardStatus.CardProgress--
	} else if score >= 0.2 {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(oneMinute)
		newCardStatus.CardProgress--
	} else {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(1 * time.Second)
		newCardStatus.CardProgress--
	}

	// make sure progress does not go below zero
	if newCardStatus.CardProgress < 0 {
		newCardStatus.CardProgress = 0
	}

	return newCardStatus
}
