package fuzzy

import (
	"testing"

	"github.com/rodrigo462003/FlickMeter/model"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s, t   string
		result int
	}{
		{"ola", "OLA", 0},
		{"ola", "ola", 0},
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"gumbo", "gambit", 3},
		{"", "", 0},
		{"hello", "hello", 0},
		{"test", "test123", 3},
		{"google", "facebook", 8},
		{"fsdjkalfjaslkd;fjdsalk", "fkjlsdafjklasdfjlsakdfj", 16},
	}

	for _, tt := range tests {
		t.Run(tt.s+" -> "+tt.t, func(t *testing.T) {
			got := levenshtein(tt.s, tt.t)
			if got != tt.result {
				t.Errorf("LevenshteinDistance(%v, %v) = %v; want %v", tt.s, tt.t, got, tt.result)
			}
		})
	}
}

func TestTree(t *testing.T) {
	movieTitles := []string{"cat", "bat", "rat", "hat", "caterpillar", "apple"}
	movies := make([]Stringer, len(movieTitles))

	for i, title := range movieTitles {
		movies[i] = model.MovieIndex{
			ID:            i + 1,
			OriginalTitle: title,
			Popularity:    0.0,
			Video:         false,
			Adult:         false,
		}
	}

	tree := NewTree(movies)

	tests := []struct {
		word     string
		expected string
	}{
		{"cat", "cat"},
		{"caterpillar", "caterpillar"},
		{"banana", "bat"},
		{"apple", "apple"},
		{"", "cat"},
		{"abcdefgh", "apple"},
		{"xyz", "rat"},
	}

	for _, test := range tests {
		t.Run("Insert and Lookup "+test.word, func(t *testing.T) {
			word := tree.Lookup(test.word)

			if word.String() != test.expected {
				expectedLev := levenshtein(test.word, test.expected)
				gotLev := levenshtein(word.String(), test.word)
				if expectedLev < gotLev {
					t.Errorf("Expected word %s, but got %s", test.expected, word)
				}
				if expectedLev > gotLev {
					t.Errorf("Test is wrong")
				}
			}
		})
	}
}
