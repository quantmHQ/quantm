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
		// If Tick returns, it means the timer fired without being restarted/stopped.
		tick.Send(ctx, true)
	})

	state.start(ctx)

	for state.Persist {
		// --- Phase 1: Wait for Acquire ---

		state.wait_acquire(ctx)

		// Reset idle timer to normal timeout while waiting
		idle.Restart(ctx, IdleTimeout)

		var rx Handler

		acquired := false
		timeout := false

		acquirer := workflow.NewSelector(ctx)

		// 1. Wait for Acquire Signal
		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), func(c workflow.ReceiveChannel, _ bool) {
			c.Receive(ctx, &rx)
			acquired = true
		})

		// 2. Wait for Idle Timeout
		acquirer.AddReceive(tick, func(_ workflow.ReceiveChannel, _ bool) {
			timeout = true
		})

		acquirer.Select(ctx)

		if timeout {
			state.idle_timeout(ctx)
			break
		}

		if !acquired {
			continue
		}

		// --- Phase 2: Lock Acquired ---

		// Pause idle timer
		idle.Restart(ctx, LongTimeout)

		state.acquired(ctx, &rx)

		// Signal client that they have the lock
		_ = workflow.
			SignalExternalWorkflow(ctx, rx.WorkflowExecutionID(), rx.WorkflowRunID(), WorkflowSignalLocked.String(), true).
			Get(ctx, nil)

		// --- Phase 3: Wait for Release or Lease Timeout ---

		// Setup Lease Timer
		timer := workflow.NewTimer(ctx, state.Timeout)

		releaser := workflow.NewSelector(ctx)
		done := false
		expired := false

		// 1. Wait for Release Signal
		releaser.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()), func(c workflow.ReceiveChannel, _ bool) {
			var rx Handler
			c.Receive(ctx, &rx)

			// Only accept release from the current holder
			if rx.WorkflowExecutionID() == state.Handler.WorkflowExecutionID() {
				done = true
			} else {
				state.ignore_release(ctx, rx.WorkflowExecutionID())
			}
		})

		// 2. Wait for Lease Timeout
		releaser.AddFuture(timer, func(_ workflow.Future) {
			expired = true
		})

		releaser.Select(ctx)

		if done {
			state.to_releasing(ctx)
			_ = workflow.
				SignalExternalWorkflow(ctx, state.Handler.WorkflowExecutionID(), state.Handler.WorkflowRunID(), WorkflowSignalReleased.String(), true).
				Get(ctx, nil)
			state.to_released(ctx)
		} else if expired {
			state.lease_expired(ctx)
		}
	}

	state.shutdown(ctx)
	idle.Stop(ctx)

	return nil
}
