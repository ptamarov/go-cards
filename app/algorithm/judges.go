package algorithm

import (
	"strings"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Judges for cards
type GuessJudge interface {
	EvaluateUserAction(card card.MemoryCard, h history.UserAction) float64
}

type NaiveJudge struct {
}

func (nj NaiveJudge) EvaluateGuess(card card.MemoryCard, g history.UserAction) float64 {

	if strings.EqualFold(card.WordToLearn, g.Guess) {
		return 1.0
	} else {
		return 0.0
	}
}

type LevenshsteinJudge struct {
	Options levenshtein.Options
}

func (lj LevenshsteinJudge) EvaluateGuess(card card.MemoryCard, g history.UserAction) float64 {
	wordToLearn := strings.ToLower(card.WordToLearn)
	guess := strings.ToLower(g.Guess)

	return levenshtein.RatioForStrings([]rune(wordToLearn), []rune(guess), lj.Options)
}

// Options for Levenshtein judge.
var DefaultOptions levenshtein.Options = levenshtein.Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 2,
	Matches: levenshtein.IdenticalRunes,
}
