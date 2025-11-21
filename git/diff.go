package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FileStatus represents the status of a modified file.
type FileStatus int

const (
	StatusModified FileStatus = iota
	StatusAdded
	StatusDeleted
	StatusRenamed
	StatusUnknown
)

// FileStat contains metadata about a changed file.
type FileStat struct {
	Path      string
	Status    FileStatus
	Additions int
	Deletions int
}

// Diff executes `git diff` and returns the raw command output.
// Callers are responsible for parsing or rendering the returned bytes.
func Diff() ([]byte, error) {
	cmd := exec.Command("git", "diff")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff: %w", err)
	}
	return out, nil
}

// GetModifiedFiles returns a list of all files with changes and their stats.
func GetModifiedFiles() ([]FileStat, error) {
	// Get file stats (additions/deletions)
	numstatCmd := exec.Command("git", "diff", "--numstat")
	numstatOut, err := numstatCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff --numstat: %w", err)
	}

	// Get file status (M/A/D/R)
	statusCmd := exec.Command("git", "diff", "--name-status")
	statusOut, err := statusCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff --name-status: %w", err)
	}

	// Parse numstat output
	statsMap := make(map[string]FileStat)
	scanner := bufio.NewScanner(bytes.NewReader(numstatOut))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		additions, _ := strconv.Atoi(parts[0])
		deletions, _ := strconv.Atoi(parts[1])
		path := strings.Join(parts[2:], " ")

		statsMap[path] = FileStat{
			Path:      path,
			Additions: additions,
			Deletions: deletions,
			Status:    StatusUnknown,
		}
	}

	// Parse status output and combine
	var files []FileStat
	scanner = bufio.NewScanner(bytes.NewReader(statusOut))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		statusChar := parts[0]
		path := strings.Join(parts[1:], " ")

		stat, exists := statsMap[path]
		if !exists {
			stat = FileStat{Path: path}
		}

		switch statusChar {
		case "M":
			stat.Status = StatusModified
		case "A":
			stat.Status = StatusAdded
		case "D":
			stat.Status = StatusDeleted
		case "R":
			stat.Status = StatusRenamed
		default:
			stat.Status = StatusUnknown
		}

		files = append(files, stat)
	}

	return files, nil
}

// GetFileDiff returns the unified diff for a specific file.
func GetFileDiff(filepath string) ([]byte, error) {
	cmd := exec.Command("git", "diff", "--", filepath)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff for %s: %w", filepath, err)
	}
	return out, nil
}
