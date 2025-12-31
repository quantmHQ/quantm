package activities

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/git"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Branch struct{}
)

// Clone clones a repo to a temp path, fetching a specified branch.
func (a *Branch) Clone(ctx context.Context, payload *defs.ClonePayload) (string, error) {
	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		slog.Warn("clone: unable to get tokenized url", "error", err)
		return "", err
	}

	path := fmt.Sprintf("/tmp/%s", payload.Path)

	// Ensure parent directory exists (though /tmp usually does)
	if err := os.MkdirAll("/tmp", 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	if _, err := git.Clone(ctx, path, url, payload.Branch); err != nil {
		slog.Warn("clone: failed", "error", err, "url", url, "path", path)
		return "", err
	}

	return path, nil
}

// RemoveDir removes a directory and handles potential errors.
func (a *Branch) RemoveDir(ctx context.Context, path string) error {
	slog.Debug("removing directory", "path", path)

	if err := os.RemoveAll(path); err != nil {
		slog.Warn("Failed to remove directory", "error", err, "path", path)
	}

	return nil
}

// Diff computes the diff between two commits using git CLI.
func (a *Branch) Diff(ctx context.Context, payload *defs.DiffPayload) (*eventsv1.Diff, error) {
	// Ensure base and head exist
	if _, err := git.RevParse(ctx, payload.Path, payload.Base); err != nil {
		slog.Warn("diff: unable to resolve base", "base", payload.Base, "error", err)
		return nil, err
	}

	if _, err := git.RevParse(ctx, payload.Path, payload.SHA); err != nil {
		slog.Warn("diff: unable to resolve head", "head", payload.SHA, "error", err)
		return nil, err
	}

	// Fetch latest origin to ensure we have the refs (CLI equivalent of refresh_remote)
	if _, err := git.Fetch(ctx, payload.Path, payload.Base); err != nil {
		slog.Warn("diff: unable to fetch base", "base", payload.Base, "error", err)
		// Proceeding anyway as local refs might be sufficient
	}

	// 1. Get File Status (Added, Modified, Deleted, Renamed)
	names, err := git.DiffStatus(ctx, payload.Path, payload.Base, payload.SHA)
	if err != nil {
		slog.Warn("diff: failed to get name-status", "error", err)
		return nil, err
	}

	// 2. Get Line Stats (Insertions, Deletions)
	stats, err := git.DiffStat(ctx, payload.Path, payload.Base, payload.SHA)
	if err != nil {
		slog.Warn("diff: failed to get numstat", "error", err)
		return nil, err
	}

	return a.parse_diff(names, stats)
}

// Rebase performs a git rebase operation. Handles conflicts and returns result.
func (a *Branch) Rebase(ctx context.Context, payload *defs.RebasePayload) (*defs.RebaseResult, error) {
	result := defs.NewRebaseResult()

	// Ensure remote is up to date
	if _, err := git.Fetch(ctx, payload.Path, payload.Rebase.Base); err != nil {
		slog.Warn("rebase: unable to fetch base", "error", err.Error(), "branch", payload.Rebase.Base, "sha", payload.Rebase.Head)
		result.SetStatusFailure(err)

		return result, nil // Return result with error status, not the error itself, to match old behavior
	}

	// Perform Rebase
	out, err := git.Rebase(ctx, payload.Path, payload.Rebase.Base)
	if err == nil {
		// Success
		result.SetStatusSuccess()
		result.Head = payload.Rebase.Head

		return result, nil
	}

	// Rebase Failed - Check for conflicts
	slog.Debug("rebase failed, checking conflicts", "error", err, "output", out)

	status, err := git.StatusPorcelain(ctx, payload.Path)
	if err != nil {
		_ = git.AbortRebase(ctx, payload.Path)

		result.SetStatusFailure(err)
		slog.Warn(
			"rebase: failed to check status after failure",
			"error", err.Error(), "branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		return result, nil
	}

	conflicts := a.parse_conflicts(status)
	if len(conflicts) > 0 {
		result.Conflicts = conflicts
		result.SetStatusConflicts()
		// Abort to clean up
		_ = git.AbortRebase(ctx, payload.Path)

		return result, nil
	}

	// Failure was not due to conflicts (or we couldn't detect them)
	result.SetStatusFailure(fmt.Errorf("rebase failed: %s", out)) // Use output as error message

	_ = git.AbortRebase(ctx, payload.Path)

	return result, nil
}

// - Helpers -

func (a *Branch) parse_diff(names, stats string) (*eventsv1.Diff, error) {
	result := &eventsv1.Diff{Files: &eventsv1.DiffFiles{}, Lines: &eventsv1.DiffLines{}}

	// Parse Name Status
	lines := strings.Split(names, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		path := parts[1]

		switch {
		case strings.HasPrefix(status, "A"):
			result.Files.Added = append(result.Files.Added, path)
		case strings.HasPrefix(status, "D"):
			result.Files.Deleted = append(result.Files.Deleted, path)
		case strings.HasPrefix(status, "M"):
			result.Files.Modified = append(result.Files.Modified, path)
		case strings.HasPrefix(status, "R"):
			if len(parts) >= 3 {
				result.Files.Renamed = append(result.Files.Renamed, &eventsv1.RenamedFile{Old: parts[1], New: parts[2]})
			}
		}
	}

	// Parse Num Stat
	stat_lines := strings.Split(stats, "\n")
	for _, line := range stat_lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		added, _ := strconv.Atoi(parts[0])
		deleted, _ := strconv.Atoi(parts[1])

		result.Lines.Added += int32(added)     // nolint: gosec
		result.Lines.Removed += int32(deleted) // nolint: gosec
	}

	return result, nil
}

func (a *Branch) parse_conflicts(status string) []string {
	var conflicts []string

	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		// Porcelain format: XY PATH
		// Unmerged states: DD, AU, UD, UA, DU, AA, UU
		code := line[:2]
		path := strings.TrimSpace(line[3:])

		if code == "UU" || code == "AA" || code == "DU" || code == "UD" || code == "UA" || code == "AU" || code == "DD" {
			conflicts = append(conflicts, path)
		}
	}

	return conflicts
}
