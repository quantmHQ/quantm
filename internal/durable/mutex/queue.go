package mutex

import (
	"sync"

	"go.breu.io/durex/queues"
)

var (
	_q     queues.Queue
	_qonce sync.Once
)

// Queue returns the singleton instance of the queues.Queue.
//
// When used in Temporal workflows, Queue must be instantiated
// during application startup, prior to any workflow execution logic accessing it.
// Lazy instantiation within a workflow violates Temporal's deterministic bounds,
// leading to non-deterministic errors.
//
// Call mutex.Queue() once in your main function or service initialization.
func Queue(opts ...queues.QueueOption) queues.Queue {
	_qonce.Do(func() {
		_q = queues.New(opts...)
	})

	return _q
}
