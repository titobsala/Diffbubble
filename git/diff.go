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
	DiffBranch                   // Compare current branch to target branch
)

// FileStatus represents the status of a modified file.
type FileStatus int

const (
	StatusModified FileStatus = iota
	StatusAdded
	StatusDeleted
	StatusRenamed
	StatusUntracked
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
func GetModifiedFiles(mode DiffMode, compareBranch string, includeUntracked bool) ([]FileStat, error) {
	// Build git diff arguments based on mode
	var diffArgs []string
	switch mode {
	case DiffBranch:
		if compareBranch == "" {
			return nil, fmt.Errorf("compare branch not specified")
		}
		// Three-dot diff: changes in HEAD since diverging from compareBranch
		diffArgs = []string{"diff", fmt.Sprintf("%s...HEAD", compareBranch)}
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

	// Fetch untracked files if requested
	if includeUntracked {
		untrackedCmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
		untrackedOut, err := untrackedCmd.Output()
		if err != nil {
			return nil, fmt.Errorf("running git ls-files: %w", err)
		}

		scanner = bufio.NewScanner(bytes.NewReader(untrackedOut))
		for scanner.Scan() {
			path := scanner.Text()
			if path == "" {
				continue
			}

			// For untracked files, we need to count lines manually or via wc -l
			// to approximate "additions"
			lineCount := 0
			if countOut, err := exec.Command("wc", "-l", path).Output(); err == nil {
				parts := strings.Fields(string(countOut))
				if len(parts) > 0 {
					lineCount, _ = strconv.Atoi(parts[0])
				}
			}

			files = append(files, FileStat{
				Path:      path,
				Status:    StatusUntracked,
				Additions: lineCount,
				Deletions: 0,
			})
		}
	}

	return files, nil
}

// GetBranches returns a list of all available branches (local and remote).
func GetBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-a", "--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running git branch: %w", err)
	}

	var branches []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// GetCurrentBranch returns the name of the current branch.
// Returns empty string if in detached HEAD state.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running git branch --show-current: %w", err)
	}

	branch := strings.TrimSpace(string(out))
	return branch, nil
}

// ValidateBranch checks if a branch exists.
func ValidateBranch(branch string) error {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("branch '%s' not found", branch)
	}
	return nil
}

// GetFileDiff returns the unified diff for a specific file.
// contextLines specifies how many context lines to show (0 for default, -1 for full file)
// mode specifies which changes to show (staged, unstaged, or all)
// compareBranch specifies the target branch for comparison (only used when mode is DiffBranch)
func GetFileDiff(filepath string, contextLines int, mode DiffMode, compareBranch string) ([]byte, error) {
	// Check if file is untracked
	isUntracked := false
	checkCmd := exec.Command("git", "ls-files", "--error-unmatch", filepath)
	if err := checkCmd.Run(); err != nil {
		// Error means it's not tracked (exit code 1)
		isUntracked = true
	}

	var args []string
	if isUntracked {
		// Use --no-index to compare against /dev/null
		args = []string{"diff", "--no-index"}
		// Add context argument
		if contextLines == -1 {
			args = append(args, "-U999999")
		} else if contextLines > 0 {
			args = append(args, fmt.Sprintf("-U%d", contextLines))
		}
		args = append(args, "--", "/dev/null", filepath)
	} else {
		// Normal diff
		switch mode {
		case DiffBranch:
			if compareBranch == "" {
				return nil, fmt.Errorf("compare branch not specified")
			}
			// Three-dot diff: changes in HEAD since diverging from compareBranch
			args = []string{"diff", fmt.Sprintf("%s...HEAD", compareBranch)}
		case DiffStaged:
			args = []string{"diff", "--cached"}
		case DiffUnstaged:
			args = []string{"diff"}
		default: // DiffAll
			args = []string{"diff", "HEAD"}
		}

		// Add context argument
		if contextLines == -1 {
			args = append(args, "-U999999")
		} else if contextLines > 0 {
			args = append(args, fmt.Sprintf("-U%d", contextLines))
		}

		args = append(args, "--", filepath)
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()

	// git diff exits with 1 if there are differences, which is not an error for us
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return out, nil
			}
		}
		return nil, fmt.Errorf("running git diff for %s: %w", filepath, err)
	}
	return out, nil
}
