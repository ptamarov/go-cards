package judges

import "testing"

func TestRemoveUmlaut(t *testing.T) {
	s1 := []rune("Übung")
	s2 := []rune("Ubung")

	for i, r := range s1 {
		newRune := removeUmlauts(r)

		if newRune != s2[i] {
			t.Errorf("expected %d but got %d instead", s2[i], newRune)
		}
	}

}

type matchTestCase struct {
	r1   rune
	r2   rune
	want bool
}

func Test_matchUmlautAndCaseInsensitive(t *testing.T) {
	var tests = []matchTestCase{
		{'a', 'A', true},
		{'Ä', 'A', true},
		{'ä', 'A', true},
		{'a', 'Ä', true},
		{'a', 'ä', true},
		{'a', 'e', false},
	}

	for _, test := range tests {
		got := matchUmlautAndCaseInsensitive(test.r1, test.r2)

		if got != test.want {
			t.Errorf("wanted %t but got %t for %v against %v", test.want, got, test.r1, test.r2)
		}
	}
}

func Test_matchUmlautInsensitive(t *testing.T) {
	var tests = []matchTestCase{
		{'a', 'A', false},
		{'Ä', 'A', true},
		{'ä', 'A', false},
		{'a', 'Ä', false},
		{'a', 'ä', true},
		{'a', 'e', false},
	}
	for _, test := range tests {
		got := matchUmlautInsensitive(test.r1, test.r2)

		if got != test.want {
			t.Errorf("wanted %t but got %t for %v against %v", test.want, got, test.r1, test.r2)
		}
	}
}

func Test_matchCaseInsensitive(t *testing.T) {
	var tests = []matchTestCase{
		{'a', 'A', true},
		{'Ä', 'A', false},
		{'ä', 'A', false},
		{'a', 'Ä', false},
		{'a', 'ä', false},
		{'a', 'e', false},
	}

	for _, test := range tests {
		got := matchCaseInsensitive(test.r1, test.r2)
		if got != test.want {
			t.Errorf("wanted %t but got %t for %v against %v", test.want, got, test.r1, test.r2)
		}
	}
}
