package translate

import "testing"

type TestCaseGetMarkedWord struct {
	prompt string
	want   string
}

func TestGetMarkedWordFromPrompt(t *testing.T) {
	tests := []TestCaseGetMarkedWord{
		{"He said \"You should tell this *person* about our project.\"",
			"person"},
	}

	for _, test := range tests {
		got, _ := GetMarkedWordFromPrompt(test.prompt)

		if got != test.want {
			t.Errorf("wanted %s but got %s when processing %s", test.want, got, test.prompt)
		}

	}
}
