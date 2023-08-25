package user

import "github.com/google/uuid"

type User struct {
	Name                   string
	LearningLang           string
	UserLang               string
	UserID                 uuid.UUID
	AnsweredCorrectlyToday int
	DailyGoal              int
}

func (u *User) IsDailyGoalReached() bool {
	return u.AnsweredCorrectlyToday >= u.DailyGoal
}

// TODO
type UserHistory struct {
	UserID string
	// History map[string]mc_user_history.UserHistoryForDeck
}

// TODO
func (uh *UserHistory) InitializeUserHistory() {
}
