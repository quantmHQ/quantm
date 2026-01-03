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

package mutex

import (
	"time"

	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/workflow"
)

type (
	// MutexStatus represents the current state of the mutex.
	MutexStatus string

	// MutexState encapsulates the state of the mutex workflow.
	MutexState struct {
		Status  MutexStatus   `json:"status"`
		Handler *Handler      `json:"handler"`
		Timeout time.Duration `json:"timeout"`
		Persist bool          `json:"persist"`

		logger *MutexLogger
	}
)

const (
	MutexStatusAcquiring MutexStatus = "mutex__acquiring"
	MutexStatusLocked    MutexStatus = "mutex__locked"
	MutexStatusReleasing MutexStatus = "mutex__releasing"
	MutexStatusReleased  MutexStatus = "mutex__released"
	MutexStatusTimeout   MutexStatus = "mutex__timeout"
)

const (
	WorkflowQueryState queues.Signal = "query__mutex__state"
)

// set_query_state sets a query handler for the mutex workflow.
func (s *MutexState) set_query_state(ctx workflow.Context) error {
	return workflow.SetQueryHandler(ctx, WorkflowQueryState.String(), func() (*MutexState, error) {
		return s, nil
	})
}

// start logs the start of the workflow.
func (s *MutexState) start(ctx workflow.Context) {
	s.logger.info(s.Handler.WorkflowExecutionID(), "start", "mutex workflow started")
}

// wait_acquire transitions to acquiring state and logs.
func (s *MutexState) wait_acquire(ctx workflow.Context) {
	s.to_acquiring(ctx)
	s.logger.info(s.Handler.WorkflowExecutionID(), "main", "waiting for lock request ...")
}

// on_aquire returns a callback for handling the acquire signal.
func (s *MutexState) on_aquire(ctx workflow.Context) func(workflow.ReceiveChannel, bool) {
	return func(c workflow.ReceiveChannel, _ bool) {
		rx := &Handler{}
		c.Receive(ctx, rx)

		s.acquired(ctx, rx)
	}
}

// acquired handles the lock acquisition and signals the client.
func (s *MutexState) acquired(ctx workflow.Context, rx *Handler) {
	s.Handler = rx
	s.Timeout = rx.Timeout
	s.to_locked(ctx)

	_ = workflow.
		SignalExternalWorkflow(ctx, s.Handler.WorkflowExecutionID(), s.Handler.WorkflowRunID(), WorkflowSignalLocked.String(), true).
		Get(ctx, nil)

	s.logger.info(s.Handler.WorkflowExecutionID(), "main", "lock acquired", "holder", rx.WorkflowExecutionID())
}

// on_idle returns a callback for handling the idle timeout.
func (s *MutexState) on_idle(ctx workflow.Context) func(workflow.ReceiveChannel, bool) {
	return func(_ workflow.ReceiveChannel, _ bool) {
		s.logger.info(s.Handler.WorkflowExecutionID(), "timeout", "shutting down due to inactivity")
		s.stop_persisting(ctx)
	}
}

// on_release returns a callback for handling the release signal.
func (s *MutexState) on_release(ctx workflow.Context) func(workflow.ReceiveChannel, bool) {
	return func(c workflow.ReceiveChannel, _ bool) {
		rx := &Handler{}
		c.Receive(ctx, rx)

		if rx.WorkflowExecutionID() == s.Handler.WorkflowExecutionID() {
			s.released(ctx)
		} else {
			s.ignore_release(ctx, rx.WorkflowExecutionID())
		}
	}
}

// released handles the lock release process.
func (s *MutexState) released(ctx workflow.Context) {
	s.to_releasing(ctx)

	_ = workflow.
		SignalExternalWorkflow(ctx, s.Handler.WorkflowExecutionID(), s.Handler.WorkflowRunID(), WorkflowSignalReleased.String(), true).
		Get(ctx, nil)

	s.to_released(ctx)
}

// ignore_release logs a warning for a release attempt from a non-holder.
func (s *MutexState) ignore_release(ctx workflow.Context, senderID string) {
	s.logger.warn(s.Handler.WorkflowExecutionID(), "release", "ignored release from non-holder", "sender", senderID)
}

// on_expired returns a callback for handling the lease expiration.
func (s *MutexState) on_expired(ctx workflow.Context) func(workflow.Future) {
	return func(_ workflow.Future) {
		s.to_timeout(ctx)
		s.logger.warn(s.Handler.WorkflowExecutionID(), "lease", "lock lease expired", "holder", s.Handler.WorkflowExecutionID())
	}
}

// shutdown logs the workflow completion.
func (s *MutexState) shutdown(ctx workflow.Context) {
	s.logger.info(s.Handler.WorkflowExecutionID(), "shutdown", "workflow completed")
}

// to_locked transitions the state to Locked.
func (s *MutexState) to_locked(ctx workflow.Context) {
	s.Status = MutexStatusLocked
	s.logger.info(s.Handler.WorkflowExecutionID(), "transition", "to locked")
}

// to_releasing transitions the state to Releasing.
func (s *MutexState) to_releasing(ctx workflow.Context) {
	s.Status = MutexStatusReleasing
	s.logger.info(s.Handler.WorkflowExecutionID(), "transition", "to releasing")
}

// to_released transitions the state to Released.
func (s *MutexState) to_released(ctx workflow.Context) {
	s.Status = MutexStatusReleased
	s.logger.info(s.Handler.WorkflowExecutionID(), "transition", "to released")
}

// to_timeout transitions the state to Timeout.
func (s *MutexState) to_timeout(ctx workflow.Context) {
	s.Status = MutexStatusTimeout
	s.logger.info(s.Handler.WorkflowExecutionID(), "transition", "to timeout")
}

// to_acquiring transitions the state to Acquiring.
func (s *MutexState) to_acquiring(ctx workflow.Context) {
	s.Status = MutexStatusAcquiring
	s.logger.info(s.Handler.WorkflowExecutionID(), "transition", "to acquiring")
}

// stop_persisting sets the persist flag to false.
func (s *MutexState) stop_persisting(ctx workflow.Context) {
	s.Persist = false
	s.logger.info(s.Handler.WorkflowExecutionID(), "persist", "stopped")
}

// restore reinitializes the logger.
func (s *MutexState) restore(ctx workflow.Context) {
	s.logger = NewMutexControllerLogger(ctx, s.Handler.ResourceID)
}

// NewMutexState creates a new MutexState instance.
func NewMutexState(ctx workflow.Context, handler *Handler) *MutexState {
	state := &MutexState{
		Status:  MutexStatusAcquiring,
		Handler: handler,
		Timeout: 0,
		Persist: true,
		logger:  NewMutexControllerLogger(ctx, handler.ResourceID),
	}

	return state
}
