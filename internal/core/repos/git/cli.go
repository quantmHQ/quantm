package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Run executes a git command in the specified directory.
func Run(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("git command failed: %w: %s", err, out.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// Clone clones a repository into a specific directory.
func Clone(ctx context.Context, dir, url, branch string) (string, error) {
	return Run(ctx, ".", "clone", "--branch", branch, url, dir)
}

// Fetch fetches the specific branch from origin.
func Fetch(ctx context.Context, dir, branch string) (string, error) {
	return Run(ctx, dir, "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
}

// StatusPorcelain returns the git status in porcelain format.
func StatusPorcelain(ctx context.Context, dir string) (string, error) {
	return Run(ctx, dir, "status", "--porcelain")
}

// DiffStatus returns the diff with --name-status.
func DiffStatus(ctx context.Context, dir, base, head string) (string, error) {
	return Run(ctx, dir, "diff", "--name-status", fmt.Sprintf("%s...%s", base, head))
}

// DiffStat returns the diff with --numstat.
func DiffStat(ctx context.Context, dir, base, head string) (string, error) {
	return Run(ctx, dir, "diff", "--numstat", fmt.Sprintf("%s...%s", base, head))
}

// Rebase attempts to rebase onto the base branch.
func Rebase(ctx context.Context, dir, base string) (string, error) {
	return Run(ctx, dir, "rebase", base)
}

// AbortRebase aborts an in-progress rebase.
func AbortRebase(ctx context.Context, dir string) error {
	_, err := Run(ctx, dir, "rebase", "--abort")
	return err
}

// RevParse verifies a commit SHA or ref exists.
func RevParse(ctx context.Context, dir, ref string) (string, error) {
	return Run(ctx, dir, "rev-parse", ref)
}
