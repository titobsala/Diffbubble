package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DiffMode represents which changes to show
type DiffMode int

const (
	DiffAll      DiffMode = iota // Both staged and unstaged (default)
	DiffStaged                   // Only staged changes (--cached)
	DiffUnstaged                 // Only unstaged changes
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
func GetModifiedFiles(mode DiffMode) ([]FileStat, error) {
	// Build git diff arguments based on mode
	var diffArgs []string
	switch mode {
	case DiffStaged:
		diffArgs = []string{"diff", "--cached"}
	case DiffUnstaged:
		diffArgs = []string{"diff"}
	default: // DiffAll
		diffArgs = []string{"diff", "HEAD"}
	}

	// Get file stats (additions/deletions)
	numstatArgs := append(diffArgs, "--numstat")
	numstatCmd := exec.Command("git", numstatArgs...)
	numstatOut, err := numstatCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff --numstat: %w", err)
	}

	// Get file status (M/A/D/R)
	statusArgs := append(diffArgs, "--name-status")
	statusCmd := exec.Command("git", statusArgs...)
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
// contextLines specifies how many context lines to show (0 for default, -1 for full file)
// mode specifies which changes to show (staged, unstaged, or all)
func GetFileDiff(filepath string, contextLines int, mode DiffMode) ([]byte, error) {
	// Build base command arguments based on mode
	var args []string
	switch mode {
	case DiffStaged:
		args = []string{"diff", "--cached"}
	case DiffUnstaged:
		args = []string{"diff"}
	default: // DiffAll
		args = []string{"diff", "HEAD"}
	}

	// Add context argument
	if contextLines == -1 {
		// Full context mode - show entire file
		args = append(args, "-U999999")
	} else if contextLines > 0 {
		// Custom context lines
		args = append(args, fmt.Sprintf("-U%d", contextLines))
	}
	// else use default context (usually 3 lines)

	// Add filepath
	args = append(args, "--", filepath)

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git diff for %s: %w", filepath, err)
	}
	return out, nil
}
