package card

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type MemoryCardData struct {
	Prompt            string    `json:"prompt"`      // a prompt must be a sentence with a unique *marked* word to learn
	Hint              string    `json:"hint"`        // a hint indicates what the prompt is missing (marked word)
	Grammar           string    `json:"grammar"`     // explains the grammar of the missing word (in context)
	PromptTranslation string    `json:"translation"` // translation of the prompt in the user's language
	CardID            uuid.UUID `json:"card_id"`
}

type MemoryCard struct {
	Prompt            string
	WordToLearn       string
	Hint              string
	Grammar           string
	PromptTranslation string
	CardID            uuid.UUID
}

func NewMemoryCard(mcdata MemoryCardData) (MemoryCard, error) {
	WordToLearn, err := GetWordToLearnFromPrompt(mcdata.Prompt)

	if err != nil {
		return MemoryCard{}, err
	}

	return MemoryCard{
		Prompt:            mcdata.Prompt,
		WordToLearn:       WordToLearn,
		Hint:              mcdata.Hint,
		Grammar:           mcdata.Grammar,
		PromptTranslation: mcdata.PromptTranslation,
		CardID:            mcdata.CardID,
	}, nil
}

func GetWordToLearnFromPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")

	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		start, end := 0, 0
	outer:
		for i := 0; i < len(prompt); i++ {
			if string(prompt[i]) == "*" {
				start = i + 1
				end = start
				for end < len(prompt) && string(prompt[end]) != "*" {
					end++

				}
				break outer
			}
		}
		return prompt[start:end], nil

	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}

func GetRedactedPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")

	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		start, end := 0, 0
	outer:
		for i := 0; i < len(prompt); i++ {
			if string(prompt[i]) == "*" {
				start = i + 1
				end = start
				for end < len(prompt) && string(prompt[end]) != "*" {
					end++

				}
				break outer
			}
		}
		return prompt[:start-1] + "_____" + prompt[end+1:], nil

	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}
