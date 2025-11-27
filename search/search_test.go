package search

import (
	"testing"

	"github.com/titobsala/Diffbubble/parser"
)

func TestSearchInRows_Overlap(t *testing.T) {
	// Mock rows
	rows := []parser.DiffRow{
		{
			Left: &parser.DiffLine{
				Number:  1,
				Content: "banananana",
				Kind:    parser.LineKindContext,
			},
		},
	}

	// Search for "nana"
	// In "banananana":
	// 0123456789
	// b a n a n a n a n a
	//     ^ ^ ^ ^          (match 1 at 2, len 4) -> ends at 6
	//             ^ ^ ^ ^  (match 2 at 6, len 4) -> ends at 10
	// Total 2 matches expected with non-overlapping logic.
	// With overlapping logic (advancing by 1), we would find:
	// 2: nana
	// 4: nana
	// 6: nana
	// Total 3 matches.

	matches := SearchInRows(rows, "nana", "test.txt", false)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
		for i, m := range matches {
			t.Logf("Match %d: col %d, text %s", i, m.Column, m.Content[m.Column:m.Column+m.Length])
		}
	}

	if len(matches) > 0 && matches[0].Column != 2 {
		t.Errorf("Expected first match at column 2, got %d", matches[0].Column)
	}
	if len(matches) > 1 && matches[1].Column != 6 {
		t.Errorf("Expected second match at column 6, got %d", matches[1].Column)
	}
}

func TestSearchInRows_CaseSensitivity(t *testing.T) {
	rows := []parser.DiffRow{
		{
			Left: &parser.DiffLine{
				Number:  1,
				Content: "Hello World",
				Kind:    parser.LineKindContext,
			},
		},
	}

	// Case insensitive (should match)
	matches := SearchInRows(rows, "hello", "test.txt", false)
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for case-insensitive search, got %d", len(matches))
	}

	// Case sensitive (should NOT match)
	matches = SearchInRows(rows, "hello", "test.txt", true)
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for case-sensitive search, got %d", len(matches))
	}

	// Case sensitive exact (should match)
	matches = SearchInRows(rows, "Hello", "test.txt", true)
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for case-sensitive exact search, got %d", len(matches))
	}
}
