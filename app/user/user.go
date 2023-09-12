package user

import (
	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
)

type User struct {
	UserID           uuid.UUID
	DeckID           uuid.UUID
	CurrentCard      card.MemoryCard
	LastAnswer       string
	UserName         string
	TargetLanguage   string
	UserLang         string
	CorrectToday     int
	DailyGoal        int
	DailyGoalReached bool
}

func (u *User) IsDailyGoalReached() bool {
	return u.CorrectToday >= u.DailyGoal
}
