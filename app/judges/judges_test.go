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
