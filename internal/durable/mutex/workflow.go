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
	expired := workflow.NewBufferedChannel(ctx, 1)

	workflow.Go(ctx, func(ctx workflow.Context) {
		idle.Tick(ctx)
		// If Tick returns, it means the timer fired without being restarted/stopped.
		expired.Send(ctx, true)
	})

	state.logger.info(state.Handler.Info.WorkflowExecution.ID, "start", "mutex workflow started")

	for state.Persist {
		// --- Phase 1: Wait for Acquire ---

		state.to_acquiring(ctx)
		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main", "waiting for lock request ...")

		// Reset idle timer to normal timeout while waiting
		idle.Restart(ctx, IdleTimeout)

		var acquireHandler Handler
		acquired := false
		timedout := false

		selector := workflow.NewSelector(ctx)

		// 1. Wait for Acquire Signal
		selector.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), func(c workflow.ReceiveChannel, _ bool) {
			c.Receive(ctx, &acquireHandler)
			acquired = true
		})

		// 2. Wait for Idle Timeout
		selector.AddReceive(expired, func(_ workflow.ReceiveChannel, _ bool) {
			timedout = true
		})

		selector.Select(ctx)

		if timedout {
			state.logger.info(state.Handler.Info.WorkflowExecution.ID, "timeout", "shutting down due to inactivity")
			state.stop_persisting(ctx)

			break
		}

		if !acquired {
			continue
		}

		// --- Phase 2: Lock Acquired ---

		// Pause idle timer
		idle.Restart(ctx, LongTimeout)

		state.set_handler(ctx, &acquireHandler)
		state.set_timeout(ctx, acquireHandler.Timeout)
		state.to_locked(ctx)

		// Signal client that they have the lock
		_ = workflow.
			SignalExternalWorkflow(ctx, acquireHandler.Info.WorkflowExecution.ID, acquireHandler.Info.WorkflowExecution.RunID, WorkflowSignalLocked.String(), true).
			Get(ctx, nil)

		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main", "lock acquired", "holder", acquireHandler.Info.WorkflowExecution.ID)

		// --- Phase 3: Wait for Release or Lease Timeout ---

		// Setup Lease Timer
		leaseTimer := workflow.NewTimer(ctx, state.Timeout)

		releaseSelector := workflow.NewSelector(ctx)
		released := false
		leaseExpired := false

		// 1. Wait for Release Signal
		releaseSelector.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()), func(c workflow.ReceiveChannel, _ bool) {
			var releaseHandler Handler
			c.Receive(ctx, &releaseHandler)

			// Only accept release from the current holder
			if releaseHandler.Info.WorkflowExecution.ID == state.Handler.Info.WorkflowExecution.ID {
				released = true
			} else {
				state.logger.warn(state.Handler.Info.WorkflowExecution.ID, "release", "ignored release from non-holder", "sender", releaseHandler.Info.WorkflowExecution.ID)
			}
		})

		// 2. Wait for Lease Timeout
		releaseSelector.AddFuture(leaseTimer, func(_ workflow.Future) {
			leaseExpired = true
		})

		releaseSelector.Select(ctx)

		if released {
			state.to_releasing(ctx)
			_ = workflow.SignalExternalWorkflow(ctx, state.Handler.Info.WorkflowExecution.ID, state.Handler.Info.WorkflowExecution.RunID, WorkflowSignalReleased.String(), true).Get(ctx, nil)
			state.to_released(ctx)
		} else if leaseExpired {
			state.to_timeout(ctx)
			state.logger.warn(state.Handler.Info.WorkflowExecution.ID, "lease", "lock lease expired", "holder", state.Handler.Info.WorkflowExecution.ID)
		}
	}

	state.logger.info(state.Handler.Info.WorkflowExecution.ID, "shutdown", "workflow completed")
	idle.Stop(ctx)

	return nil
}
