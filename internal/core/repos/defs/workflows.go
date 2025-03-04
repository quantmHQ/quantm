// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
