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

	"go.breu.io/durex/dispatch"
	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/workflow"
)

const (
	DefaultTimeout = 10 * time.Minute // DefaultTimeout is the default lease timeout.
	SignalTimeout  = 1 * time.Minute  // SignalTimeout is the timeout for signal handshakes.
	MaxAcquireWait = 24 * time.Hour   // MaxAcquireWait is the maximum time to wait for the lock.
)

const (
	WorkflowSignalAcquire  queues.Signal = "mutex__acquire"
	WorkflowSignalLocked   queues.Signal = "mutex__locked"
	WorkflowSignalRelease  queues.Signal = "mutex__release"
	WorkflowSignalReleased queues.Signal = "mutex__released"
)

type (
	Option func(*Handler)

	// Mutex defines the signature for the workflow mutex.
	Mutex interface {
		// OnAcquire blocks until the lock is acquired, executes fn, and then releases the lock.
		// It returns an error if the lock cannot be acquired or the context is cancelled.
		OnAcquire(ctx workflow.Context, fn func(workflow.Context)) error
	}

	// Handler is the Mutex handler.
	Handler struct {
		ResourceID string              `json:"resource_id"` // ResourceID identifies the resource being locked.
		Info       *workflow.Info      `json:"info"`        // Info holds the workflow info that requests the mutex.
		Execution  *workflow.Execution `json:"execution"`   // Execution holds the mutex workflow execution details.
		Timeout    time.Duration       `json:"timeout"`     // Timeout sets the lease timeout.
		logger     *MutexLogger
	}
)

// WithResourceID sets the resource ID for the mutex workflow.
func WithResourceID(id string) Option {
	return func(m *Handler) {
		m.ResourceID = id
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) Option {
	return func(m *Handler) {
		m.Timeout = timeout
	}
}

// New returns a new Mutex.
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

	h.logger.info(h.Info.WorkflowExecution.ID, "create", "mutex handler initialized")

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
			h.logger.error(h.Info.WorkflowExecution.ID, "release", "failed to release lock", err)
		}
	}()

	// 3. Execute Critical Section
	fn(ctx)

	return nil
}

// Internal helper to acquire the lock.
func (h *Handler) acquire(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "requesting lock")

	c := dispatch.WithDefaultActivityContext(ctx)

	exe := &workflow.Execution{}
	if err := workflow.ExecuteActivity(c, AcquireMutexActivity, h).Get(c, exe); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "acquire", "unable to request lock", err)
		return NewAcquireLockError(h.ResourceID)
	}

	h.Execution = exe
	h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "waiting for lock")

	locked := false
	timeout := false
	waiter := workflow.NewSelector(ctx)

	waiter.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalLocked.String()), func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, &locked)
	})

	waiter.AddFuture(workflow.NewTimer(ctx, MaxAcquireWait), func(_ workflow.Future) {
		timeout = true
	})

	waiter.Select(ctx)

	if timeout {
		h.logger.warn(h.Info.WorkflowExecution.ID, "acquire", "timeout waiting for lock")
		return NewAcquireLockError(h.ResourceID)
	}

	if locked {
		h.logger.info(h.Info.WorkflowExecution.ID, "acquire", "lock acquired")
		return nil
	}

	return NewAcquireLockError(h.ResourceID)
}

// Internal helper to release the lock.
func (h *Handler) release(ctx workflow.Context) error {
	h.logger.info(h.Info.WorkflowExecution.ID, "release", "requesting release")

	if err := workflow.
		SignalExternalWorkflow(ctx, h.Execution.ID, h.Execution.RunID, WorkflowSignalRelease.String(), h).
		Get(ctx, nil); err != nil {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "unable to request release", err)
		return NewReleaseLockError(h.ResourceID)
	}

	h.logger.info(h.Info.WorkflowExecution.ID, "release", "waiting for release confirmation")

	released := false
	timeout := false
	waiter := workflow.NewSelector(ctx)

	waiter.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalReleased.String()), func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, &released)
	})

	waiter.AddFuture(workflow.NewTimer(ctx, SignalTimeout), func(_ workflow.Future) {
		timeout = true
	})

	waiter.Select(ctx)

	if timeout {
		h.logger.warn(h.Info.WorkflowExecution.ID, "release", "timeout waiting for release confirmation")
		return NewReleaseLockError(h.ResourceID)
	}

	if released {
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
