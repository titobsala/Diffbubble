package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CommitInfo represents metadata about a single commit.
type CommitInfo struct {
	Hash      string    // Short hash (7 chars)
	Message   string    // First line of commit message
	Author    string    // Author name
	Date      time.Time // Commit date
	DateStr   string    // Relative date string (e.g., "2 hours ago")
	Additions int       // Total lines added in commit
	Deletions int       // Total lines deleted in commit
}

// GetBranchCommits retrieves commits in current branch that are not in baseBranch.
// Uses git log baseBranch...HEAD to get commits unique to current branch.
// Returns empty slice if no commits found (not an error).
func GetBranchCommits(n int, baseBranch string) ([]CommitInfo, error) {
	if n <= 0 {
		return []CommitInfo{}, nil
	}

	// Get commits in HEAD that are not in baseBranch (PR-style)
	args := []string{"log", fmt.Sprintf("%s...HEAD", baseBranch), "-n", strconv.Itoa(n), "--pretty=format:%h|%s|%an|%ar|%at"}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		// Check if it's because there are no commits
		if exitErr, ok := err.(*exec.ExitError); ok {
			// exit code 128 typically means no commits or invalid revision
			if exitErr.ExitCode() == 128 {
				return []CommitInfo{}, nil
			}
		}
		return nil, fmt.Errorf("running git log: %w", err)
	}

	// Parse commit list
	var commits []CommitInfo
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse format: hash|message|author|relative_date|unix_timestamp
		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}

		hash := parts[0]
		message := parts[1]
		author := parts[2]
		dateStr := parts[3]
		timestamp, _ := strconv.ParseInt(parts[4], 10, 64)

		// Get stats for this commit
		additions, deletions, err := GetCommitStats(hash)
		if err != nil {
			// If we can't get stats, use 0 but don't fail the whole operation
			additions = 0
			deletions = 0
		}

		commits = append(commits, CommitInfo{
			Hash:      hash,
			Message:   message,
			Author:    author,
			Date:      time.Unix(timestamp, 0),
			DateStr:   dateStr,
			Additions: additions,
			Deletions: deletions,
		})
	}

	return commits, nil
}

// GetCommitStats retrieves additions and deletions for a specific commit.
// Returns (additions, deletions, error).
func GetCommitStats(commitHash string) (int, int, error) {
	// Use git show --shortstat to get compact stats
	// Format: "1 file changed, 42 insertions(+), 18 deletions(-)"
	cmd := exec.Command("git", "show", "--shortstat", "--format=", commitHash)
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("running git show --shortstat for %s: %w", commitHash, err)
	}

	// Parse the shortstat line
	line := strings.TrimSpace(string(out))
	if line == "" {
		// Commit with no file changes (e.g., merge commit or empty commit)
		return 0, 0, nil
	}

	additions := 0
	deletions := 0

	// Use regex to extract numbers
	// Pattern: "X insertion(s)(+)" or "Y deletion(s)(-)"
	insertionRegex := regexp.MustCompile(`(\d+) insertion`)
	deletionRegex := regexp.MustCompile(`(\d+) deletion`)

	if matches := insertionRegex.FindStringSubmatch(line); len(matches) > 1 {
		additions, _ = strconv.Atoi(matches[1])
	}

	if matches := deletionRegex.FindStringSubmatch(line); len(matches) > 1 {
		deletions, _ = strconv.Atoi(matches[1])
	}

	return additions, deletions, nil
}
