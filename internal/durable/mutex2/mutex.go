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

package mutex2

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"
)

type (
	// Mutex defines the signature for the workflow mutex.
	Mutex interface {
		// OnAcquire blocks until the lock is acquired, executes fn, and then releases the lock.
		// It returns an error if the lock cannot be acquired or the context is cancelled.
		OnAcquire(ctx workflow.Context, fn func(workflow.Context)) error
	}
)

// New returns a new Mutex and initializes the underlying workflow.
// It performs the "Prepare" step immediately.
func New(ctx workflow.Context, opts ...Option) (Mutex, error) {
	h := &Handler{Timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(h)
	}

	h.Info = workflow.GetInfo(ctx)
	h.logger = NewMutexHandlerLogger(ctx, h.ResourceID)

	if err := h.validate(); err != nil {
		h.logger.error(h.Info.WorkflowExecution.ID, "create", "validate error", err)
		return nil, err
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "create", "initializing mutex")

	// Prepare/Initialize the workflow via Activity
	ctx = dispatch.WithDefaultActivityContext(ctx)

	exe := &workflow.Execution{}
	if err := workflow.ExecuteActivity(ctx, PrepareMutexActivity, h).Get(ctx, exe); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "create", "Unable to prepare mutex", err)
		return nil, NewPrepareMutexError(h.ResourceID)
	}

	h.Execution = exe
	h.logger.info(h.Info.WorkflowExecution.ID, "create", "mutex initialized", "id", h.Execution.ID)

	return h, nil
}

// OnAcquire blocks until acquired (or timeout), executes the closure, and releases the lock.
func (h *Handler) OnAcquire(ctx workflow.Context, fn func(workflow.Context)) error {
	// 1. Acquire
	if err := h.acquire(ctx); err != nil {
		return err
	}

	// 2. Ensure Release
	defer func() {
		if err := h.release(ctx); err != nil {
			h.logger.error(h.Info.WorkflowExecution.ID, "release", "failed to release lock during cleanup", err)
		}
	}()

	// 3. Execute Critical Section
	// We pass the same context for now. In the future, we could pass a cancellable context
	// linked to the lock lease.
	fn(ctx)

	return nil
}

// Internal helper to acquire the lock
func (h *Handler) acquire(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "requesting lock")

	ok := true

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalAcquire.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "acquire", "Unable to request lock", err)
		return NewAcquireLockError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "waiting for lock")
	workflow.GetSignalChannel(ctx, WorkflowSignalLocked.String()).Receive(ctx, &ok)
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "lock acquired")

	if ok {
		return nil
	}

	return NewAcquireLockError(h.ResourceID)
}

// Internal helper to release the lock
func (h *Handler) release(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "release", "requesting release")

	orphan := false

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, "", WorkflowSignalRelease.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "Unable to request release", err)
		return NewReleaseLockError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "release", "waiting for release")
	workflow.GetSignalChannel(ctx, WorkflowSignalReleased.String()).Receive(ctx, &orphan)

	if orphan {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "lock released, orphaned", nil)
	} else {
		h.logger.info(h.Info.WorkflowExecution.ID, "release", "lock released")
	}

	return nil
}

// validate checks if the mutex is properly configured.
func (h *Handler) validate() error {
	if h.ResourceID == "" {
		return ErrNoResourceID
	}

	if h.Info == nil {
		return ErrNilContext
	}

	return nil
}
