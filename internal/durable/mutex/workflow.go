// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
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

package mutex

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/durable/periodic"
)

const (
	IdleTimeout = 10 * time.Minute
	LongTimeout = 365 * 24 * time.Hour // Effectively infinite for "pausing" the idle timer
)

// MutexWorkflow is the mutex workflow. It controls access to a resource.
//
// It operates as a serialized state machine:
// 1. Wait for Acquire (or Idle Timeout).
// 2. Lock Resource.
// 3. Wait for Release (or Lease Timeout).
// 4. Repeat.
func MutexWorkflow(ctx workflow.Context, state *MutexState) error {
	state.restore(ctx)

	// Setup Query Handler
	_ = state.set_query_state(ctx)

	// Idle Timer Setup
	// We use the periodic package to manage an idle timeout.
	// If no activity occurs within IdleTimeout, the workflow shuts down.
	idle := periodic.New(ctx, IdleTimeout)
	tick := workflow.NewBufferedChannel(ctx, 1)

	workflow.Go(ctx, func(ctx workflow.Context) {
		idle.Tick(ctx)
		tick.Send(ctx, true)
	})

	for state.Persist {
		state.ready(ctx)
		idle.Restart(ctx, IdleTimeout)

		workflow.NewSelector(ctx).
			AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), state.on_aquire(ctx)).
			AddReceive(tick, state.stop(ctx)).
			Select(ctx)

		if state.Status == MutexStatusLocked {
			idle.Restart(ctx, LongTimeout)
			lease := workflow.NewTimer(ctx, state.Timeout)

			for state.Status == MutexStatusLocked {
				workflow.NewSelector(ctx).
					AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()), state.on_release(ctx)).
					AddFuture(lease, state.expire(ctx)).
					Select(ctx)
			}
		}

		if !state.Persist {
			break
		}
	}

	idle.Stop(ctx)

	return nil
}
