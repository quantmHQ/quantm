# Git Command Mappings

This document maps the current `git2go` operations used in `internal/core/repos/activities/branch.go` to their equivalent `git` CLI commands. This will serve as a specification for the new CLI-based implementation.

|Operation|Current `git2go` Implementation|Proposed CLI Command|Notes|
|:---|:---|:---|:---|
|**Clone**|`git.Clone(url, path, opts)` (opts.CheckoutBranch = payload.Branch)|`git clone --branch <branch> <url> <path>`|Clones the repository into a specific path and checks out the target branch immediately.|
|**Fetch**|`remote.Fetch(...)`|`git fetch origin <branch>:<branch>`|Fetches the specific branch from origin and updates the local reference.|
|**Diff (Files)**|`repo.DiffTreeToTree(...)`; `diff.NumDeltas()`|`git diff --name-status <base>...<head>`|Used to identify which files were Added, Modified, Deleted, or Renamed. `...` (triple dot) is safer for feature branches as it diffs from the merge base.|
|**Diff (Lines)**|`diff.Stats()`; `stats.Insertions()`, `stats.Deletions()`|`git diff --shortstat <base>...<head>` OR `git diff --numstat <base>...<head>`|`--shortstat` gives a summary ("X files changed, Y insertions(+), Z deletions(-)"), while `--numstat` gives per-file stats which can be summed.|
|**Rebase**|`repo.InitRebase(...)`; `rebase.Next()`, `rebase.Commit()`|`git rebase <base>`|The CLI handles the rebase process automatically. We will likely not need to step through commits manually unless we want fine-grained progress updates.|
|**Check Conflicts**|`idx.HasConflicts()`; `idx.ConflictIterator()`|`git status --porcelain` OR `git diff --name-only --diff-filter=U`|If `git rebase` fails (exit code != 0), we check for unmerged files (status 'U') to identify conflicts.|
|**Abort Rebase**|`rebase.Abort()`|`git rebase --abort`|Cleans up the state if a rebase fails due to conflicts or other errors.|
|**Resolve Commit**|`repo.LookupCommit(oid)`|`git rev-parse <sha>`|Verifies a commit SHA exists.|
|**Tree from Ref**|`repo.References.Lookup(...)` -> `Commit.Tree()`|`git rev-parse <branch>^{tree}`|Internally used by `git2go` for diffing, but CLI `diff` handles branch names/SHAs directly.|
