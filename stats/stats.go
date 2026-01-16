package stats

import (
	"fmt"

	"github.com/titobsala/diffbubble/git"
)

// Stats represents aggregated statistics for a set of changes.
type Stats struct {
	TotalAdditions int            // Sum of all additions
	TotalDeletions int            // Sum of all deletions
	TotalFiles     int            // Number of files changed
	FileStats      []git.FileStat // Individual file statistics
	MaxChanges     int            // Max changes in any single file (for bar scaling)
}

// CommitInfo represents a single commit with statistics.
type CommitInfo struct {
	Hash      string // Short hash (7 chars)
	Message   string // First line of commit message
	Author    string // Author name
	DateStr   string // Relative date string (e.g., "2 hours ago")
	Additions int    // Total lines added in commit
	Deletions int    // Total lines deleted in commit
}

// CommitHistory represents a list of recent commits.
type CommitHistory struct {
	Commits      []CommitInfo // List of commits
	TotalCommits int          // Total number of commits
	DisplayLimit int          // How many commits to display (10-20)
}

// BranchStats represents statistics for current branch vs base branch.
type BranchStats struct {
	BranchChanges Stats         // Current branch vs base (main or remote)
	Commits       CommitHistory // Commits unique to branch
	BaseBranch    string        // "main" or "origin/feature"
	CurrentBranch string        // Current branch name
	CompareMode   string        // "vs main" or "vs remote"
}

// ComputeStats aggregates statistics from a list of files.
func ComputeStats(files []git.FileStat) Stats {
	stats := Stats{
		TotalFiles: len(files),
		FileStats:  files,
	}

	for _, file := range files {
		stats.TotalAdditions += file.Additions
		stats.TotalDeletions += file.Deletions

		// Track max changes for any single file (for bar chart scaling)
		totalChanges := file.Additions + file.Deletions
		if totalChanges > stats.MaxChanges {
			stats.MaxChanges = totalChanges
		}
	}

	return stats
}

// ComputeBranchStats computes statistics for current branch vs base branch.
// baseBranch: branch to compare against ("main", "master", or remote tracking branch)
// compareRemote: if true, compare with remote tracking branch instead of main
func ComputeBranchStats(baseBranch string, compareRemote bool) (*BranchStats, error) {
	var base string
	var mode string

	if compareRemote {
		// Compare with remote tracking branch
		upstream, err := git.GetUpstreamBranch()
		if err != nil {
			return nil, fmt.Errorf("no remote tracking branch configured")
		}
		base = upstream
		mode = "vs remote"
	} else {
		// Compare with main
		if baseBranch == "" {
			var err error
			baseBranch, err = git.GetMainBranch()
			if err != nil {
				return nil, fmt.Errorf("cannot find main branch: %w", err)
			}
		}
		base = baseBranch
		mode = "vs " + baseBranch
	}

	// Get current branch name
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("cannot get current branch: %w", err)
	}

	// Check if we're on the base branch
	if currentBranch == base || currentBranch == "" {
		return nil, fmt.Errorf("already on %s - no changes to compare", base)
	}

	// Get diff between base and HEAD
	files, err := git.GetModifiedFiles(git.DiffBranch, base, false)
	if err != nil {
		return nil, fmt.Errorf("cannot get modified files: %w", err)
	}

	// Get commits unique to current branch
	commitLimit := 15
	gitCommits, err := git.GetBranchCommits(commitLimit, base)
	if err != nil {
		// Don't fail on error, just use empty commits
		gitCommits = []git.CommitInfo{}
	}

	// Convert git.CommitInfo to stats.CommitInfo
	commits := make([]CommitInfo, len(gitCommits))
	for i, gc := range gitCommits {
		commits[i] = CommitInfo{
			Hash:      gc.Hash,
			Message:   gc.Message,
			Author:    gc.Author,
			DateStr:   gc.DateStr,
			Additions: gc.Additions,
			Deletions: gc.Deletions,
		}
	}

	history := CommitHistory{
		Commits:      commits,
		TotalCommits: len(commits),
		DisplayLimit: commitLimit,
	}

	return &BranchStats{
		BranchChanges: ComputeStats(files),
		Commits:       history,
		BaseBranch:    base,
		CurrentBranch: currentBranch,
		CompareMode:   mode,
	}, nil
}
