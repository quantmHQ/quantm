package defs

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

const (
	LabelMerge    = "quantm-merge"
	LabelPriority = "quantm-priority"
)

// signals.
const (
	SignalPush                     queues.Signal = "push"              // signals a push event.
	SignalRef                      queues.Signal = "ref"               // signals a branch event.
	SignalPullRequest              queues.Signal = "pr"                // signals a pull request event.
	SignalRebase                   queues.Signal = "rebase"            // signals a rebase event.
	SignalPullRequestLabel         queues.Signal = "pr_label"          // signals a pull request label event.
	SignalPullRequestReview        queues.Signal = "pr_review"         // signals a pull request review event.
	SignalPullRequestReviewComment queues.Signal = "pr_review_comment" // signals a pull request review comment event.
	SignalMergeQueue               queues.Signal = "merge_queue"       // signals a pull request queue event.
)

const (
	QueryRepoForEventParent queues.Query = "event_parent" // query to find the parent event for the given event
)

type (
	ClonePayload struct {
		Repo   *entities.Repo    `json:"repo"`
		Hook   eventsv1.RepoHook `json:"hook"`
		Branch string            `json:"branch"`
		Path   string            `json:"path"`
		SHA    string            `json:"at"`
	}

	DiffPayload struct {
		Path string `json:"path"`
		Base string `json:"base"`
		SHA  string `json:"sha"`
	}

	DiffFiles struct {
		Added      []string `json:"added"`
		Deleted    []string `json:"deleted"`
		Modified   []string `json:"modified"`
		Renamed    []string `json:"renamed"`
		Copied     []string `json:"copied"`
		TypeChange []string `json:"typechange"`
		Unreadable []string `json:"unreadable"`
		Ignored    []string `json:"ignored"`
		Untracked  []string `json:"untracked"`
		Conflicted []string `json:"conflicted"`
	}

	DiffLines struct {
		Added   int `json:"added"`
		Removed int `json:"removed"`
	}

	DiffResult struct {
		Files DiffFiles `json:"files"`
		Lines DiffLines `json:"lines"`
	}

	SignalBranchPayload struct {
		Signal queues.Signal  `json:"signal"`
		Repo   *entities.Repo `json:"repo"`
		Branch string         `json:"branch"`
	}

	SignalTrunkPayload struct {
		Signal queues.Signal  `json:"signal"`
		Repo   *entities.Repo `json:"repo"`
	}

	SignalQueuePayload struct{}
)

// Sum returns the sum of added and removed lines.
func (d *DiffLines) Sum() int {
	return d.Added + d.Removed
}
