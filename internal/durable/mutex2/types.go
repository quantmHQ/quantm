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
	"time"

	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/workflow"
)

const (
	DefaultTimeout = 0 * time.Minute // DefaultTimeout is the default timeout for the mutex.
)

const (
	WorkflowSignalPrepare        queues.Signal = "mutex__prepare"
	WorkflowSignalAcquire        queues.Signal = "mutex__acquire"
	WorkflowSignalLocked         queues.Signal = "mutex__locked"
	WorkflowSignalRelease        queues.Signal = "mutex__release"
	WorkflowSignalReleased       queues.Signal = "mutex__released"
	WorkflowSignalCleanup        queues.Signal = "mutex__cleanup"
	WorkflowSignalCleanupDone    queues.Signal = "mutex__cleanup_done"
	WorkflowSignalCleanupDoneAck queues.Signal = "mutex__cleanup_done_ack"
	WorkflowSignalShutDown       queues.Signal = "mutex__shutdown"
)

type (
	Option func(*Handler)

	// Handler is the Mutex handler.
	Handler struct {
		ResourceID string              `json:"resource_id"` // ResourceID identifies the resource being locked.
		Info       *workflow.Info      `json:"info"`        // Info holds the workflow info that requests the mutex.
		Execution  *workflow.Execution `json:"execution"`   // Info holds the workflow info that holds the mutex.
		Timeout    time.Duration       `json:"timeout"`     // Timeout sets the timeout, after which the lock is automatically released.
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
