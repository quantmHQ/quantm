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
	git "github.com/jeffwelling/git2go/v37"

	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	RebaseOperationKind string
	RebaseStatus        string

	RebasePayload struct {
		Rebase *eventsv1.Rebase `json:"rebase"`
		Path   string           `json:"path"`
	}

	RebaseOperation struct {
		Kind    RebaseOperationKind `json:"kind"`
		Status  RebaseStatus        `json:"status"`
		Head    string              `json:"head"`
		Message string              `json:"message"`
		Error   string              `json:"error,omitempty"`
	}

	RebaseResult struct {
		Head         string            `json:"head"`
		Status       RebaseStatus      `json:"status"`
		Operations   []RebaseOperation `json:"operations"`
		TotalCommits uint              `json:"count"`
		Conflicts    []string          `json:"conflicts"`
		Error        string            `json:"error,omitempty"`
	}
)

const (
	RebaseStatusSuccess   RebaseStatus = "success"
	RebaseStatusFailure   RebaseStatus = "failure"
	RebaseStatusConflicts RebaseStatus = "conflicts"
	RebaseStatusUpToDate  RebaseStatus = "up_to-date"
	RebaseStatusAborted   RebaseStatus = "aborted"
	RebaseStatusPartial   RebaseStatus = "partial"
)

const (
	RebaseOperationKindPick   RebaseOperationKind = "pick"
	RebaseOperationKindReword RebaseOperationKind = "reword"
	RebaseOperationKindEdit   RebaseOperationKind = "edit"
	RebaseOperationKindSquash RebaseOperationKind = "squash"
	RebaseOperationKindFixup  RebaseOperationKind = "fixup"
)

var (
	gitOpTypeMap = map[git.RebaseOperationType]RebaseOperationKind{
		git.RebaseOperationPick:   RebaseOperationKindPick,
		git.RebaseOperationReword: RebaseOperationKindReword,
		git.RebaseOperationEdit:   RebaseOperationKindEdit,
		git.RebaseOperationSquash: RebaseOperationKindSquash,
		git.RebaseOperationFixup:  RebaseOperationKindFixup,
	}
)

func (r *RebaseResult) HasConflicts() bool {
	return len(r.Conflicts) > 0
}

func (r *RebaseResult) AppliedCommit() int {
	return len(r.Operations)
}

func (r *RebaseResult) AddOperation(op git.RebaseOperationType, status RebaseStatus, head, message string, err error) {
	err_ := ""

	if err != nil {
		status = RebaseStatusFailure
		err_ = err.Error()
	}

	r.Operations = append(
		r.Operations,
		RebaseOperation{
			Kind:    gitOpTypeMap[op],
			Status:  status,
			Head:    head,
			Message: message,
			Error:   err_,
		})
}

func NewRebaseResult() *RebaseResult {
	return &RebaseResult{
		Status:     RebaseStatusFailure,
		Conflicts:  []string{},
		Operations: []RebaseOperation{},
	}
}

func (r *RebaseResult) SetStatusSuccess() {
	r.Status = RebaseStatusSuccess
}

func (r *RebaseResult) SetStatusFailure(err error) {
	r.Status = RebaseStatusFailure
	r.Error = err.Error()
}

func (r *RebaseResult) SetStatusUpToDate() {
	r.Status = RebaseStatusUpToDate
}

func (r *RebaseResult) SetStatusConflicts() {
	r.Status = RebaseStatusConflicts
}

func (r *RebaseResult) SetStatusAborted() {
	r.Status = RebaseStatusAborted
}

func (r *RebaseResult) SetStatusPartial() {
	r.Status = RebaseStatusPartial
}
