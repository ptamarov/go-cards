package judges

import (
	"unicode"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// levenshteinJudge implements the Judge interface. It evaluates an guess using
// the Levenshtein metric. By default, it is case insensitive and umlaut insensitive.
type levenshteinJudge struct {
	CaseInsensitive   bool
	UmlautInsensitive bool
}

func NewLevenshteinJudge(caseInsensitive, umlautInsensitive bool) levenshteinJudge {
	return levenshteinJudge{CaseInsensitive: caseInsensitive, UmlautInsensitive: umlautInsensitive}
}

// EvaluateUserAction computes the Levenshtein ratio between card.Answer and g.Guess, a float in [0..1].
func (lj *levenshteinJudge) EvaluateUserAction(card card.MemoryCard, g history.UserAction) float64 {
	answer := card.Answer
	guess := g.Guess

	var options = levenshtein.Options{
		InsCost: 1,
		DelCost: 1,
		SubCost: 2,
		Matches: levenshtein.IdenticalRunes,
	}

	if lj.CaseInsensitive && lj.UmlautInsensitive {
		options.Matches = matchUmlautAndCaseInsensitive
	}
	if lj.CaseInsensitive && !lj.UmlautInsensitive {
		options.Matches = matchCaseInsensitive
	}
	if !lj.CaseInsensitive && lj.UmlautInsensitive {
		options.Matches = matchUmlautInsensitive
	}
	return levenshtein.RatioForStrings([]rune(answer), []rune(guess), options)
}

// Matches two runes if they are equal up to case.
func matchCaseInsensitive(a rune, b rune) bool {
	return unicode.ToLower(a) == unicode.ToLower(b)
}

// Matches two runes if they are equal up to umlauts.
func matchUmlautInsensitive(a rune, b rune) bool {
	return removeUmlauts(a) == removeUmlauts(b)
}

// Matches two runes if they are equal up to case and umlauts.
func matchUmlautAndCaseInsensitive(a rune, b rune) bool {
	a, b = removeUmlauts(a), removeUmlauts(b)
	return unicode.ToLower(a) == unicode.ToLower(b)
}

// Removes umlauts from vowels.
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
