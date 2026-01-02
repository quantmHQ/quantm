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
// IMPORTANT: Do not use this function directly. Instead, use mutex.New to create and interact with mutex instances.
//
// The workflow consists of three main event loops:
//  1. Main loop: Handles acquiring, releasing, and timing out of locks.
//  2. Prepare loop: Listens for and handles preparation of lock requests.
//  3. Cleanup loop: Manages the cleanup process and potential workflow shutdown.
//
// It operates as a state machine, transitioning between MutexStatus states:
//
// Acquiring -> Locked -> Releasing -> Released (or Timeout)
//
// Uses two pools to manage lock requests:
//   - Main pool: Tracks active lock requests and currently held locks.
//   - Orphans pool: Tracks locks that have timed out.
//
// Responds to several signals:
//   - WorkflowSignalPrepare: Prepares a new lock request.
//   - WorkflowSignalAcquire: Attempts to acquire a lock.
//   - WorkflowSignalRelease: Releases a held lock.
//   - WorkflowSignalCleanup: Initiates the cleanup process.
func MutexWorkflow(ctx workflow.Context, state *MutexState) error {
	state.restore(ctx)

	shutdown, shutdownfn := workflow.NewFuture(ctx)

	_ = state.set_query_state(ctx)

	workflow.Go(ctx, state.on_prepare(ctx))
	workflow.Go(ctx, state.on_cleanup(ctx, shutdownfn))

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

	for state.Persist {
		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main", "waiting for lock request ...")

		// Reset idle timer to normal timeout while waiting
		idle.Restart(ctx, IdleTimeout)

		acquirer := workflow.NewSelector(ctx)

		// 1. Wait for Acquire
		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), state.on_acquire(ctx))

		// 2. Wait for Manual Shutdown/Cleanup
		acquirer.AddFuture(shutdown, state.on_terminate(ctx))

		// 3. Wait for Idle Timeout
		acquirer.AddReceive(expired, func(c workflow.ReceiveChannel, m bool) {
			state.logger.info(state.Handler.Info.WorkflowExecution.ID, "timeout", "shutting down due to inactivity")
			state.stop_persisting(ctx)
		})

		acquirer.Select(ctx)

		if !state.Persist {
			break // Shutdown signal received or idle timeout
		}

		// If we are here, we acquired the lock (since shutdown/idle didn't happen).
		// "Pause" the idle timer while we process the lock.
		idle.Restart(ctx, LongTimeout)

		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main", "lock acquired!")
		state.to_locked(ctx)

		for {
			state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main", "waiting for release or timeout ...")

			releaser := workflow.NewSelector(ctx)

			releaser.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()),
				state.on_release(ctx),
			)
			releaser.AddFuture(workflow.NewTimer(ctx, state.Timeout), state.on_abort(ctx))

			releaser.Select(ctx)

			if state.Status == MutexStatusReleased || state.Status == MutexStatusTimeout {
				state.to_acquiring(ctx)
				break
			}
		}
	}

	_ = workflow.Sleep(ctx, 500*time.Millisecond) // Wait for cleanup ack if needed

	state.logger.info(state.Handler.Info.WorkflowExecution.ID, "shutdown", "shutdown!")

	// Ensure idle timer is stopped to clean up goroutine
	idle.Stop(ctx)

	return nil
}
