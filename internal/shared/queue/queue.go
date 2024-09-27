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

package queue

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"go.temporal.io/sdk/client"
	sdk "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type (
	// queue defines the basic queue.
	queue struct {
		name                Name          // The name of the temporal queue.
		prefix              string        // The prefix for the Workflow ID.
		workflowMaxAttempts int32         // The maximum number of attempts for a workflow.
		worker              worker.Worker // worker singleton.
		workeronce          sync.Once     // worker singleton lock.
	}
)

func (q Name) String() string {
	return string(q)
}

func (q *queue) Name() string {
	return q.name.String()
}

func (q *queue) Prefix() string {
	return q.prefix
}

func (q *queue) WorkflowID(options ...WorkflowOptionProvider) string {
	prefix := ""
	opts := NewWorkflowOptions(options...)

	if opts.IsChild() {
		prefix = opts.ParentWorkflowID()
	} else {
		prefix = q.prefix
	}

	return fmt.Sprintf("%s.%s", prefix, opts.Suffix())
}

func (q *queue) WorkflowOptions(options ...WorkflowOptionProvider) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:          q.WorkflowID(options...),
		TaskQueue:   q.Name(),
		RetryPolicy: &sdk.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

func (q *queue) ChildWorkflowOptions(options ...WorkflowOptionProvider) workflow.ChildWorkflowOptions {
	return workflow.ChildWorkflowOptions{
		WorkflowID:  q.WorkflowID(options...),
		RetryPolicy: &sdk.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

func (q *queue) Worker(c client.Client) worker.Worker {
	q.workeronce.Do(func() {
		slog.Info("queue: creating worker ...", "queue", q.name, "id_prefix", q.Prefix())

		options := worker.Options{OnFatalError: func(err error) { panic(err) }, EnableSessionWorker: true}

		q.worker = worker.New(c, q.Name(), options)
	})

	return q.worker
}

func (q *queue) Listen(interrupt <-chan any) error {
	if q.worker == nil {
		return ErrNilWorker
	}

	return q.worker.Run(interrupt)
}

func (q *queue) Stop(ctx context.Context) error {
	if q.worker == nil {
		return ErrNilWorker
	}

	q.worker.Stop()

	return nil
}

// WithName sets the queue name and the prefix for the workflow ID.
func WithName(name Name) QueueOption {
	return func(q Queue) {
		q.(*queue).name = name
		q.(*queue).prefix = DefaultPrefix + name.String()
	}
}

// WithWorkflowMaxAttempts sets the maximum number of attempts for a workflow.
func WithWorkflowMaxAttempts(attempts int32) QueueOption {
	return func(q Queue) {
		q.(*queue).workflowMaxAttempts = attempts
	}
}

// NewQueue creates a new queue with the given options.
func NewQueue(opts ...QueueOption) Queue {
	q := &queue{workflowMaxAttempts: DefaultWorkflowMaxAttempts}
	for _, opt := range opts {
		opt(q)
	}

	return q
}
