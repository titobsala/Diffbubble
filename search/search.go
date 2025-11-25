package search

import (
	"strings"

	"github.com/titobsala/Diffbubble/parser"
)

// Match represents a search match in the diff.
type Match struct {
	FileName   string // File containing the match
	RowIndex   int    // Index in the DiffRow slice
	Side       string // "left" or "right"
	LineNumber int    // Line number in the file
	Column     int    // Column where match starts
	Length     int    // Length of the matched text
	Content    string // The line content containing the match
}

// SearchInRows searches for a query string within a single file's diff rows.
// Returns all matches found in the current file.
func SearchInRows(rows []parser.DiffRow, query string, fileName string) []Match {
	if query == "" {
		return nil
	}

	var matches []Match
	lowerQuery := strings.ToLower(query)

	for rowIdx, row := range rows {
		// Search in left side
		if row.Left != nil && row.Left.Kind != parser.LineKindHeader {
			content := row.Left.Content
			lowerContent := strings.ToLower(content)

			// Find all occurrences in this line
			startPos := 0
			for {
				pos := strings.Index(lowerContent[startPos:], lowerQuery)
				if pos == -1 {
					break
				}
				actualPos := startPos + pos
				matches = append(matches, Match{
					FileName:   fileName,
					RowIndex:   rowIdx,
					Side:       "left",
					LineNumber: row.Left.Number,
					Column:     actualPos,
					Length:     len(query),
					Content:    content,
				})
				startPos = actualPos + 1
			}
		}

		// Search in right side
		if row.Right != nil && row.Right.Kind != parser.LineKindHeader {
			content := row.Right.Content
			lowerContent := strings.ToLower(content)

			// Find all occurrences in this line
			startPos := 0
			for {
				pos := strings.Index(lowerContent[startPos:], lowerQuery)
				if pos == -1 {
					break
				}
				actualPos := startPos + pos
				matches = append(matches, Match{
					FileName:   fileName,
					RowIndex:   rowIdx,
					Side:       "right",
					LineNumber: row.Right.Number,
					Column:     actualPos,
					Length:     len(query),
					Content:    content,
				})
				startPos = actualPos + 1
			}
		}
	}

	return matches
}

// GetMatchPosition returns the viewport scroll position for a given match.
// This helps auto-scroll to the match location.
func GetMatchPosition(match Match) int {
	return match.RowIndex
}
