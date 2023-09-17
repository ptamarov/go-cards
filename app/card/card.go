package card

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// MemoryCard stores the data of a language flash card.
type MemoryCard struct {
	ID                uuid.UUID `json:"card_id"`
	Grammar           string    `json:"grammar"`            // Explains the grammar of the missing word (in context)
	Hint              string    `json:"hint"`               // Hints at the missing word or short phrase in the prompt.
	LanguageToLearn   string    `json:"lang_learn"`         // The language to learn.
	UserLanguage      string    `json:"lang_user"`          // The language the user uses to learn
	Prompt            string    `json:"prompt"`             // A sentence with a unique *marked* word to learn.
	PromptTranslation string    `json:"prompt_translation"` // A translation of the prompt in the user's language.
	Answer            string    `json:"answer"`             // The word or short phrase to learn.
}

// GetWordToLearnFromPrompt returns the word to learn from a prompt. The word must be marked with the delimiter *
// exactly once. It returns an error if there aren't exactly two delimiters.
//
// Examples:
//   - "This is *valid*." returns "valid"
//   - "No closing *star." returns "bad input: no closing * found“
//   - "No delimiters." returns "bad input: no opening * found“
//   - "*Too* many *delimiters*!" returns "bad input: too many delimiters found“
func GetWordToLearnFromPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")
	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		re := regexp.MustCompile(`\*(.*)\*`)
		return re.FindStringSubmatch(prompt)[1], nil
	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}

// GetRedactedPrompt takes a "raw" valid prompt and redacts anything within (and including)
// the delimiter *. It returns an error if there aren't exactly two delimiters.
//
// Examples:
//   - "This is *valid*." returns "This is _____."
//   - "No closing *star" returns "bad input: no closing * found“
//   - "No delimiters" returns "bad input: no opening * found“
//   - "*Too* many *delimiters*!" returns "bad input: too many delimiters found“
func GetRedactedPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")
	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		var newPrompt string
		re := regexp.MustCompile(`\*.*\*`)
		newPrompt = re.ReplaceAllString(prompt, "_____")
		return newPrompt, nil
	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}
