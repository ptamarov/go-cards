package card

import "testing"

type TestCaseGetWord struct {
	prompt string
	want   string
	err    bool
}

type TestCaseGetPrompt struct {
	card MemoryCard
	want string
	err  bool
}

func TestGetWordToLearnFromPrompt(t *testing.T) {
	tests := []TestCaseGetWord{
		{"Here is a *word*", "word", false},
		{"Here is a *bad input", "", true},
		{"Here is a *good* input", "good", false},
		{"Here is more bad input", "", true},
		{"Here are **too many** delimiters.", "", true},
		{"Here is ** nothing.", "", false},
	}

	for _, test := range tests {
		got, err := GetWordToLearnFromPrompt(test.prompt)

		if err != nil {
			if !test.err {
				t.Error("was not expecting an error but got", err)
			}
		} else {
			if got != test.want {
				t.Errorf("was expecting %s but got %s", test.want, got)
			}
		}
	}
}

func TestGetPromptFromCard(t *testing.T) {
	tests := []TestCaseGetPrompt{
		{MemoryCard{Prompt: "Here is a *word*."}, "Here is a _____.", false},
		{MemoryCard{Prompt: "Here is *bad input"}, "", true},
		{MemoryCard{Prompt: "Here is *good* input."}, "Here is _____ input.", false},
		{MemoryCard{Prompt: "More bad input."}, "", true},
		{MemoryCard{Prompt: "Here is ** nothing."}, "Here is _____ nothing.", false},
		{MemoryCard{Prompt: ""}, "", true},
		{MemoryCard{Prompt: "Here is *bad** input!"}, "", true},
		{MemoryCard{Prompt: "**Bad input*!"}, "", true},
	}

	for _, test := range tests {
		got, err := GetRedactedPrompt(test.card.Prompt)

		if err != nil {
			if !test.err {
				t.Error("was not expecting an error but got", err)
			}
		} else {
			if got != test.want {
				t.Errorf("was expecting %s but got %s", test.want, got)
			}
		}
	}
}
