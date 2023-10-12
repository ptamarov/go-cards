package models

import "github.com/ptamarov/go-cards/app/card"

// TemplateDate holds data sent from handlers to templates
type TemplateData struct {
	Answer    string
	CSRFToken string

	DailyGoal     int
	AnsweredToday int
	NotSeen       int
	Learned       int
	InProgress    int
	Total         int
	CardsReady    int
	Minutes       int
	Seconds       int

	InProgressPerc  float64
	BarProgressPerc float64
	NotSeenPerc     float64
	LearnedPerc     float64
	CorrectRate     float64

	DailyGoalReached bool

	CardData card.MemoryCard
}
