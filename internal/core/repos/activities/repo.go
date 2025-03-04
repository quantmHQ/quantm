// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024, 2025.
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

package activities

import (
	"context"
	"log/slog"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/durable"
)

type (
	Repo struct{}
)

const (
	WorkflowBranch = "Branch" // WorkflowBranch is string representation of workflows.Branch
	WorkflowTrunk  = "Trunk"  // WorkflowTrunk is string representation of workflows.Trunk
)

// ForwardToBranch sends a signal to a branch workflow, starting it if it doesn't exist.
func (a *Repo) ForwardToBranch(ctx context.Context, payload *defs.SignalBranchPayload, event, state any) error {
	id := defs.BranchWorkflowOptions(payload.Repo, payload.Branch)
	run, err := durable.
		OnCore().
		SignalWithStartWorkflow(ctx, id, payload.Signal, event, WorkflowBranch, state)

	if err != nil {
		slog.Warn("fwd_to_branch: unable to signal", "id", id.IDSuffix(), "error", err.Error())
	} else {
		slog.Info("fwd_to_branch: signaled", "id", id.IDSuffix(), "run_id", run.GetRunID())
	}

	return err
}

// ForwardToTrunk sends a signal to the trunk workflow, starting it if it doesn't exist.
func (a *Repo) ForwardToTrunk(ctx context.Context, payload *defs.SignalTrunkPayload, event, state any) error {
	_, err := durable.
		OnCore().
		SignalWithStartWorkflow(ctx, defs.TrunkWorkflowOptions(payload.Repo), payload.Signal, event, WorkflowTrunk, state)

	return err
}

// ForwardToQueue is a no-op for now, but is reserved for a queueing mechanism.
func (a *Repo) ForwardToQueue(ctx context.Context, payload *defs.SignalQueuePayload, event, state any) error {
	return nil
}
