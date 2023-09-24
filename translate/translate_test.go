package translate

import "testing"

type testCase struct {
	word     string
	sentence string
	want     string
}

func TestBlankOutWordInSentence(t *testing.T) {
	tests := []testCase{
		{word: "word", sentence: "This is a word.", want: "This is a ______."},
		{word: "bird", sentence: "This is a word.", want: "This is a word."},
		{word: "is", sentence: "This is a word.", want: "This ______ a word."},
		{word: "promise", sentence: "This is a word (promise).", want: "This is a word (______)."},
		{word: "word", sentence: "This is a word. Another word.", want: "This is a ______. Another ______."},
	}

	for _, test := range tests {
		got := BlankOutWordInSentence(test.sentence, test.word)
		if test.want != got {
			t.Errorf("processed sentence failed: wanted [%s] but got [%s]", test.want, got)
		}
	}
}

func TestProcessSentence(t *testing.T) {
	tests := []testCase{
		{word: "word", sentence: "This is a word.", want: "This is a *word*."},
		{word: "bird", sentence: "This is a word.", want: "This is a word."},
		{word: "is", sentence: "This is a word.", want: "This *is* a word."},
		{word: "promise", sentence: "This is a word (promise).", want: "This is a word (*promise*)."},
	}

	for _, test := range tests {
		got := ProcessSentence(test.sentence, test.word)
		if test.want != got {
			t.Errorf("processed sentence failed: wanted [%s] but got [%s]", test.want, got)
		}
	}
}
