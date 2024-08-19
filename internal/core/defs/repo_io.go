package defs

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/shared"
)

type (
	PullRequestAction string
)

// RepoIO signals.
const (
	RepoIOSignalPush             shared.WorkflowSignal = "repo_io__push"
	RepoIOSignalCreateOrDelete   shared.WorkflowSignal = "repo_io__create_or_delete"
	RepoIOSignalRebase           shared.WorkflowSignal = "repo_io__rebase"
	RepoIOSignalPullRequest      shared.WorkflowSignal = "repo_io__pull_request"
	RepoIOSignalLabel            shared.WorkflowSignal = "repo_io__label"
	RepoIOSignalQueueAdd         shared.WorkflowSignal = "repo_io__queue__add"
	RepoIOSignalQueueAddPriority shared.WorkflowSignal = "repo_io__queue__add__priority"
	RepoIOSignalQueuePromote     shared.WorkflowSignal = "repo_io__queue__promote"
	RepoIOSignalQueueDemote      shared.WorkflowSignal = "repo_io__queue__demote"
)

const (
	PullRequestActionCreated       PullRequestAction = "created"
	PullRequestActionLabeled       PullRequestAction = "labeled"
	PullRequestActionClosed        PullRequestAction = "closed"
	PullRequestActionMerged        PullRequestAction = "merged"
	PullRequestActionReviewRequest PullRequestAction = "review_request"
	PullRequestActionApproved      PullRequestAction = "approved"
)

// signal payloads.
type (
	RepoIOSignalPushPayload struct {
		BranchRef      string         `json:"branch_ref"`
		Before         string         `json:"before"`
		After          string         `json:"after"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		CtrlID         string         `json:"ctrl_id"` // CtrlID represents the id of the provider repo in the quantm DB. Should be UUID.
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		Commits        RepoIOCommits  `json:"commits"`
		User           *auth.TeamUser `json:"user"`
		Author         string         `json:"author"`
	}

	RepoIOSignalCreateOrDeletePayload struct {
		IsCreated      bool           `json:"is_created"`
		Ref            string         `json:"ref"`
		RefType        string         `json:"ref_type"`
		DefaultBranch  string         `json:"default_branch"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		CtrlID         string         `json:"ctrl_id"` // CtrlID represents the id of the provider repo in the quantm DB. Should be UUID.
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		User           *auth.TeamUser `json:"user"`
	}

	RepoIOSignalPullRequestPayload struct {
		Action         string         `json:"action"`
		Number         shared.Int64   `json:"number"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		BaseBranch     string         `json:"base_branch"`
		HeadBranch     string         `json:"head_branch"`
		CtrlID         string         `json:"ctrl_id"`
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		User           *auth.TeamUser `json:"user"` // TODO: need to find more optimze way
	}

	RepoIOSignalLabelPayload struct {
		Action         string         `json:"action"`
		Number         shared.Int64   `json:"number"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		BaseBranch     string         `json:"base_branch"`
		HeadBranch     string         `json:"head_branch"`
		CtrlID         string         `json:"ctrl_id"`
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		User           *auth.TeamUser `json:"user"` // TODO: need to find more optimze way
	}
)

// RepoIO types.
type (
	RepoIOProviderInfo struct {
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"owner"`
		DefaultBranch  string       `json:"default_branch"`
		ProviderID     string       `json:"provider_id"`
		InstallationID shared.Int64 `json:"installation_id"`
	}

	RepoIOClonePayload struct {
		Repo   *Repo                    `json:"repo"`   // Repo is the db record of the repo
		Push   *RepoIOSignalPushPayload `json:"push"`   // Push event payload
		Branch string                   `json:"branch"` // Branch to clone
		Path   string                   `json:"path"`   // Path to clone to
	}

	RepoIOPushBranchPayload struct {
		Branch string `json:"branch"`
		Path   string `json:"path"`
		Force  bool   `json:"force"`
	}

	RepoIODetectChangesPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"defualt_branch"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIOChanges struct {
		Added      shared.Int64 `json:"added"`
		Removed    shared.Int64 `json:"removed"`
		Modified   []string     `json:"modified"`
		Delta      shared.Int64 `json:"delta"`
		CompareUrl string       `json:"compare_url"`
		RepoUrl    string       `json:"repo_url"`
	}

	RepoIOCommit struct {
		SHA       string        `json:"sha"`
		Message   string        `json:"message"`
		Author    string        `json:"author"`
		Timestamp time.Time     `json:"timestamp"`
		Changes   RepoIOChanges `json:"changes"`
	}

	RepoIOMergePRPayload struct {
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"owner"`
		DefaultBranch  string       `json:"defualt_branch"`
		TargetBranch   string       `json:"target_branch"`
		InstallationID shared.Int64 `json:"installation_id"`
	}

	RepoIOCommits []RepoIOCommit

	RepoIOPullRequest struct {
		Number     shared.Int64 `json:"number"`
		HeadBranch string       `json:"head_branch"`
		BaseBranch string       `json:"base_branch"`
	}

	RepoIORebaseAtCommitResponse struct {
		SHA        string `json:"sha"`
		Message    string `json:"message"`
		InProgress bool   `json:"in_progress"`
	}

	RepoIOGetRepoByProviderIDPayload struct {
		ProviderID string `json:"provider_id"`
	}

	RepoIOSignalBranchCtrlPayload struct {
		Repo    *Repo                 `json:"repo"`    // Repo is the db record of the repo
		Branch  string                `json:"branch"`  // Branch to signal
		Signal  shared.WorkflowSignal `json:"signal"`  // Signal to send
		Payload any                   `json:"payload"` // Payload to send
	}

	RepoIOSignalQueueCtrlPayload struct {
		Repo    *Repo                 `json:"repo"`    // Repo is the db record of the repo
		Branch  string                `json:"branch"`  // Branch to signal
		Signal  shared.WorkflowSignal `json:"signal"`  // Signal to send
		Payload any                   `json:"payload"` // Payload to send
	}
)

func (commits RepoIOCommits) Size() int {
	return len(commits)
}

// Latest returns the most recent RepoIOCommit from the provided slice of commits.
// If the slice is empty, it returns nil.
//
// FIXME: it should iterate over the commits and return the most recent commit based on the timestamp.
func (commits RepoIOCommits) Latest() *RepoIOCommit {
	if len(commits) == 0 {
		return nil
	}

	return &commits[0]
}

func (signal *RepoIOSignalCreateOrDeletePayload) ForBranch(ctx workflow.Context) bool {
	return signal.RefType == "branch"
}

func (signal *RepoIOSignalCreateOrDeletePayload) ForTag(ctx workflow.Context) bool {
	return signal.RefType == "tag"
}
