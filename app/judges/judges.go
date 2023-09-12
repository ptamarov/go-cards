package judges

import (
	"strings"
	"unicode"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Judges for cards
type Judge interface {
	EvaluateUserAction(card card.MemoryCard, h history.UserAction) float64
}

// NaiveJudge evaluates a guess by checking if the two strings are equal.
type NaiveJudge struct {
	CaseInsensitive   bool
	UmlautInsensitive bool
}

func (nj *NaiveJudge) EvaluateUserAction(card card.MemoryCard, g history.UserAction) float64 {
	if strings.EqualFold(card.Answer, g.Guess) {
		return 1.0
	} else {
		return 0.0
	}
}

type LevenshsteinJudge struct {
	CaseInsensitive   bool
	UmlautInsensitive bool
}

func (lj *LevenshsteinJudge) EvaluateUserAction(card card.MemoryCard, g history.UserAction) float64 {
	answer := card.Answer
	guess := g.Guess

	var options = levenshtein.Options{
		InsCost: 1,
		DelCost: 1,
		SubCost: 2,
		Matches: levenshtein.IdenticalRunes,
	}

	if lj.CaseInsensitive && lj.UmlautInsensitive {
		options.Matches = MatchUmlautAndCaseInsensitive
	}
	if lj.CaseInsensitive && !lj.UmlautInsensitive {
		options.Matches = MatchLowercaseInsensitive
	}
	if !lj.CaseInsensitive && lj.UmlautInsensitive {
		options.Matches = MatchUmlautInsensitive
	}
	return levenshtein.RatioForStrings([]rune(answer), []rune(guess), options)
}

// Options for Levenshtein judge.
var DefaultOptions = levenshtein.DefaultOptions

func MatchLowercaseInsensitive(a rune, b rune) bool {
	return unicode.ToLower(a) == unicode.ToLower(b)
}

func MatchUmlautInsensitive(a rune, b rune) bool {
	return removeUmlauts(a) == removeUmlauts(b)
}

func MatchUmlautAndCaseInsensitive(a rune, b rune) bool {
	a, b = removeUmlauts(a), removeUmlauts(b)
	return unicode.ToLower(a) == unicode.ToLower(b)
}

func removeUmlauts(a rune) rune {
	if newRune, ok := umlautToVowelMap[a]; ok {
		return newRune
	} else {
		return a
	}
}

var umlautToVowelMap map[rune]rune = map[rune]rune{
	196: 65,
	214: 79,
	220: 85,
	228: 97,
	246: 111,
	252: 117,
}
