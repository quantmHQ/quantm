# Git Implementation Specification

## 1. Executive Summary
This document specifies the replacement of the existing CGO-based git implementation (`git2go`) with a pure Go wrapper around the system's `git` CLI (`os/exec`). This decision stems from a comprehensive investigation into the project's "broken" state caused by complex build dependencies and fragmented git library usage.

**Primary Goal:** Eliminate CGO dependencies to fix the build/deployment pipeline.
**Critical Constraint:** The **Contract** (public interfaces, struct definitions, JSON signatures) is the **Holy Grail**. It is non-negotiable. The surface area and signatures of Activities and Definition structs must remain exactly as they are.

## 2. Investigation & Findings

### 2.1. Current State Analysis
We analyzed `internal/core/repos` and found a fragmented and problematic git architecture:
1.  **`git2go` (libgit2):** The production code in `internal/core/repos/activities/branch.go` relies heavily on `github.com/jeffwelling/git2go/v37`. This introduces a heavy CGO dependency, which is the root cause of the current "fucked" state regarding builds and stability.
2.  **`go-git`:** An alternative implementation existed in `internal/core/repos/git` but was confirmed to be dead code, unused by any active workflow (`branch.go`, `repo.go`, `trunk.go`). It has since been removed.
3.  **Missing CLI:** References were found to a previous "command line" implementation, but the code itself was missing.

### 2.2. The "Better Path"
The investigation concluded that the most robust path forward is to re-implement the Git layer using `os/exec`.
*   **Stability:** Relies on the battle-tested system `git` binary.
*   **Portability:** Removes CGO, making cross-compilation and containerization trivial.
*   **Feature Completeness:** Support for complex operations like "ahead-of-line testing" (rebasing multiple branches onto a shadow branch) is significantly easier with the CLI than with library bindings.

### 2.3. Artifact: Command Mapping
As part of the investigation, we mapped every existing `git2go` operation in `activities/branch.go` to its CLI equivalent. This mapping is the authoritative guide for the implementation.

| Operation       | Current `git2go` Implementation                                      | Proposed CLI Command                                                             | Notes                                                                                                                                                    |
| :-------------- | :------------------------------------------------------------------- | :------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Clone**       | `git.Clone(url, path, opts)` (opts.CheckoutBranch = payload.Branch)  | `git clone --branch <branch> <url> <path>`                                       | Clones the repository into a specific path and checks out the target branch immediately.                                                                   |
| **Fetch**       | `remote.Fetch(...)`                                                  | `git fetch origin <branch>:<branch>`                                             | Fetches the specific branch from origin and updates the local reference.                                                                                 |
| **Diff (Files)**| `repo.DiffTreeToTree(...)`; `diff.NumDeltas()`                       | `git diff --name-status <base>...<head>`                                         | Used to identify which files were Added, Modified, Deleted, or Renamed. `...` (triple dot) is safer for feature branches as it diffs from the merge base. |
| **Diff (Lines)**| `diff.Stats()`; `stats.Insertions()`, `stats.Deletions()`            | `git diff --shortstat <base>...<head>` OR `git diff --numstat <base>...<head>`   | `--shortstat` gives a summary ("X files changed, Y insertions(+), Z deletions(-)"), while `--numstat` gives per-file stats which can be summed.          |
| **Rebase**      | `repo.InitRebase(...)`; `rebase.Next()`, `rebase.Commit()`           | `git rebase <base>`                                                              | The CLI handles the rebase process automatically. We will likely not need to step through commits manually unless we want fine-grained progress updates. |
| **Check Conflicts**| `idx.HasConflicts()`; `idx.ConflictIterator()`                    | `git status --porcelain` OR `git diff --name-only --diff-filter=U`               | If `git rebase` fails (exit code != 0), we check for unmerged files (status 'U') to identify conflicts.                                                  |
| **Abort Rebase**| `rebase.Abort()`                                                     | `git rebase --abort`                                                             | Cleans up the state if a rebase fails due to conflicts or other errors.                                                                                  |
| **Resolve Commit**| `repo.LookupCommit(oid)`                                           | `git rev-parse <sha>`                                                            | Verifies a commit SHA exists.                                                                                                                            |
| **Tree from Ref**| `repo.References.Lookup(...)` -> `Commit.Tree()`                    | `git rev-parse <branch>^{tree}`                                                  | Internally used by `git2go` for diffing, but CLI `diff` handles branch names/SHAs directly.                                                              |

## 3. Technical Specification

### 3.1. Package Design: `internal/core/repos/git`
A new package will be created to encapsulate the shell execution logic. It will **not** contain business logic or domain types; it is purely an execution engine.

```go
package git

// Run executes a git command in the specified directory.
// It returns the combined stdout/stderr and an error if the command fails.
func Run(ctx context.Context, dir string, args ...string) (string, error)

// Typed wrappers for safety and clarity:
func Clone(ctx context.Context, dir, url, branch string) (string, error)
func Fetch(ctx context.Context, dir, remote, branch string) (string, error)
func DiffNameStatus(ctx context.Context, dir, base, head string) (string, error)
func DiffNumStat(ctx context.Context, dir, base, head string) (string, error)
func Rebase(ctx context.Context, dir, base string) (string, error)
func AbortRebase(ctx context.Context, dir string) error
func Status(ctx context.Context, dir string) (string, error)
```

### 3.2. Activity Integration (`activities/branch.go`)
The `Branch` activity will be refactored to use the new `git` package.
**Crucial:** The method signatures `Clone`, `Diff`, and `Rebase` must **NOT** change. They must return the exact same structs populated with data parsed from the CLI output.

#### Contract: `Clone`
*   **Input:** `defs.ClonePayload`
*   **Output:** `string` (path), `error`
*   **Implementation:** Call `git.Clone`. Verify directory existence.

#### Contract: `Diff`
*   **Input:** `defs.DiffPayload`
*   **Output:** `*eventsv1.Diff`, `error`
*   **Implementation:**
    1.  Call `git.DiffNameStatus`. Parse output line-by-line (`A`, `M`, `D`, `R` prefixes) to populate `eventsv1.DiffFiles`.
    2.  Call `git.DiffNumStat`. Parse output (tab-separated numbers) to calculate total `Added` and `Removed` lines for `eventsv1.DiffLines`.

#### Contract: `Rebase`
*   **Input:** `defs.RebasePayload`
*   **Output:** `*defs.RebaseResult`, `error`
*   **Implementation:**
    1.  Call `git.Rebase`.
    2.  **Success:** Return `RebaseStatusSuccess`.
    3.  **Failure (Exit Code != 0):**
        *   Call `git.Status` or `git diff --name-only --diff-filter=U` to identify conflicting files.
        *   Populate `RebaseResult.Conflicts`.
        *   Set status to `RebaseStatusConflicts`.
        *   Call `git.AbortRebase` to clean up.
    4.  **Other Errors:** Return `RebaseStatusFailure`.

### 3.3. Definition Adjustments (`defs/rebase.go`)
**Warning:** The `defs` package defines the contract. The structs `RebaseResult` and `RebaseOperation` must remain identical in structure and JSON serialization.
*   **The Problem:** `defs/rebase.go` currently imports `git2go` to use `git.RebaseOperationType` in a helper method (`AddOperation`).
*   **The Fix:**
    1.  Remove the `import "github.com/jeffwelling/git2go/v37"`.
    2.  Change the `op` parameter in `AddOperation` from `git.RebaseOperationType` to a native type (e.g., `string` or `RebaseOperationKind`) that matches the enum already defined in the file.
    3.  **Do NOT touch** the fields of `RebaseResult` or `RebaseOperation`.

## 4. Implementation Plan
1.  **Create Wrapper:** Implement `internal/core/repos/git/cli.go`.
2.  **Refactor Activities:** Rewrite `Branch.Clone`, `Branch.Diff`, and `Branch.Rebase` in `internal/core/repos/activities/branch.go` to use the CLI wrapper and parse output.
3.  **Sanitize Defs:** Remove the `git2go` import from `internal/core/repos/defs/rebase.go` and update the `AddOperation` helper signature to use native types, strictly preserving the data contract.
4.  **Verify:** Ensure no CGO dependencies remain and the build succeeds.
