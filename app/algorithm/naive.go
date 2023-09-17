package algorithm

import (
	"time"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// NaiveAlgorithm computes a new card status by only looking at the last action
// performed by the user. It wraps a judge that uses the Levenshtein distance to
// judge actions.
type NaiveAlgorithm struct {
	judge judges.LevenshteinJudge
}

func (na *NaiveAlgorithm) GetJudge() judges.Judge {
	return &na.judge
}

func (na *NaiveAlgorithm) ComputeNewCardStatus(c card.MemoryCard, a []history.UserAction, s CardStatus) CardStatus {
	oneMinute := 60 * time.Second
	oneHour := 60 * oneMinute
	oneDay := 24 * oneHour

	if len(a) == 0 { // if no actions then status is unchanged
		return s
	}

	action := a[len(a)-1]
	score := na.judge.EvaluateUserAction(c, action)

	var newCardStatus CardStatus
	timeNowUTC := time.Now().UTC()

	if score == 1.0 {
		newCardStatus.NextAvailableDate = timeNowUTC.Add(oneDay)
		newCardStatus.CardProgress = s.CardProgress + 1
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
