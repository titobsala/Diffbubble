package git

import (
	"fmt"
	"os/exec"
)

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
