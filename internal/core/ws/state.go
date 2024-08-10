package ws

import (
	"fmt"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	// Connections represents the state of websocket connections,
	// managing the relationship between users and queues.
	Connections struct {
		UserQueue  map[string]string              `json:"user_queue"`
		QueueUsers map[string]map[string]struct{} `json:"queue_users"`
		mu         workflow.Mutex
		logger     log.Logger
	}

	AddUserSignal struct {
		UserID string
		Queue  string
	}

	RemoveUserSignal struct {
		UserID string
	}

	FlushQueueSignal struct {
		Queue string
	}

	WorkerAddedSignal struct {
		Queue string
	}
)

// GetQueueForUser returns the queue name for a given user ID.
//
// Example:
//
//	queue, exists := connections.GetQueueForUser(ctx, "user123")
//	if exists {
//	    fmt.Printf("User is in queue: %s\n", queue)
//	}
func (con *Connections) GetQueueForUser(ctx workflow.Context, user_id string) (string, bool) {
	if err := con.mu.Lock(ctx); err != nil {
		return "", false
	}
	defer con.mu.Unlock()

	queue, exists := con.UserQueue[user_id]

	return queue, exists
}

// AddUserToQueue adds a user to a specified queue.
// If the user is already in a queue, they are removed from the old queue first.
//
// Example:
//
//	err := connections.AddUserToQueue(ctx, "user123", "queue1")
//	if err != nil {
//	    log.Printf("Failed to add user to queue: %v", err)
//	}
func (con *Connections) AddUserToQueue(ctx workflow.Context, user_id, queue string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if oldQueue, exists := con.UserQueue[user_id]; exists {
		delete(con.QueueUsers[oldQueue], user_id)

		if len(con.QueueUsers[oldQueue]) == 0 {
			delete(con.QueueUsers, oldQueue)
		}
	}

	con.UserQueue[user_id] = queue
	if _, exists := con.QueueUsers[queue]; !exists {
		con.QueueUsers[queue] = make(map[string]struct{})
	}

	con.QueueUsers[queue][user_id] = struct{}{}

	return nil
}

// RemoveUserFromQueue removes a user from their current queue.
//
// Example:
//
//	err := connections.RemoveUserFromQueue(ctx, "user123")
//	if err != nil {
//	    log.Printf("Failed to remove user from queue: %v", err)
//	}
func (con *Connections) RemoveUserFromQueue(ctx workflow.Context, user_id string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if queueName, exists := con.UserQueue[user_id]; exists {
		delete(con.UserQueue, user_id)
		delete(con.QueueUsers[queueName], user_id)

		if len(con.QueueUsers[queueName]) == 0 {
			delete(con.QueueUsers, queueName)
		}
	}

	return nil
}

// GetUsersInQueue returns a list of user IDs in a specified queue.
//
// Example:
//
//	users, err := connections.GetUsersInQueue(ctx, "queue1")
//	if err != nil {
//	    log.Printf("Failed to get users in queue: %v", err)
//	} else {
//	    fmt.Printf("Users in queue: %v\n", users)
//	}
func (con *Connections) GetUsersInQueue(ctx workflow.Context, queue string) ([]string, error) {
	if err := con.mu.Lock(ctx); err != nil {
		return nil, err
	}
	defer con.mu.Unlock()

	users := make([]string, 0, len(con.QueueUsers[queue]))
	for user_id := range con.QueueUsers[queue] {
		users = append(users, user_id)
	}

	return users, nil
}

// ClearQueue removes all users from a specified queue.
//
// Example:
//
//	err := connections.ClearQueue(ctx, "queue1")
//	if err != nil {
//	    log.Printf("Failed to clear queue: %v", err)
//	}
func (con *Connections) ClearQueue(ctx workflow.Context, queue string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if users, exists := con.QueueUsers[queue]; exists {
		for user_id := range users {
			delete(con.UserQueue, user_id)
		}

		delete(con.QueueUsers, queue)
	}

	return nil
}

// Restore reinitializes the mutex.
// This should be called when deserializing the Connections struct.
//
// Example:
//
//	var connections Connections
//	// ... deserialize connections ...
//	connections.Restore(ctx)
func (con *Connections) Restore(ctx workflow.Context) {
	con.mu = workflow.NewMutex(ctx)
	con.logger = workflow.GetLogger(ctx)
}

func (con *Connections) prefixed(msg string) string {
	return fmt.Sprintf("ws: %s", msg)
}

func (con *Connections) info(msg string, keyvals ...any) {
	con.logger.Info(con.prefixed(msg), keyvals...)
}

func (con *Connections) debug(msg string, keyvals ...any) {
	con.logger.Debug(con.prefixed(msg), keyvals...)
}

func (con *Connections) warn(msg string, keyvals ...any) {
	con.logger.Warn(con.prefixed(msg), keyvals...)
}

func (con *Connections) error(msg string, keyvals ...any) {
	con.logger.Error(con.prefixed(msg), keyvals...)
}

func (con *Connections) on_add(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal AddUserSignal

		rx.Receive(ctx, &signal)

		if err := con.AddUserToQueue(ctx, signal.UserID, signal.Queue); err != nil {
			con.error("Failed to add user to queue", "user_id", signal.UserID, "queue", signal.Queue, "error", err)
		} else {
			con.info("Added user to queue", "user_id", signal.UserID, "queue", signal.Queue)
		}
	}
}

func (con *Connections) on_remove(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal RemoveUserSignal

		rx.Receive(ctx, &signal)

		if err := con.RemoveUserFromQueue(ctx, signal.UserID); err != nil {
			con.error("Failed to remove user from queue", "user_id", signal.UserID, "error", err)
		} else {
			con.info("Removed user from queue", "user_id", signal.UserID)
		}
	}
}

func (con *Connections) on_flush(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal FlushQueueSignal

		rx.Receive(ctx, &signal)

		if err := con.ClearQueue(ctx, signal.Queue); err != nil {
			con.error("Failed to flush queue", "queue", signal.Queue, "error", err)
		} else {
			con.info("Flushed queue", "queue", signal.Queue)
		}
	}
}

func (con *Connections) on_worker_added(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal WorkerAddedSignal

		rx.Receive(ctx, &signal)

		con.info("Worker added", "queue", signal.Queue)
	}
}

// NewConnections creates a new Connections instance.
//
// Example:
//
//	connections := NewConnections(ctx)
//	// Use the connections object to manage websocket connections
func NewConnections(ctx workflow.Context) *Connections {
	return &Connections{
		UserQueue:  make(map[string]string),
		QueueUsers: make(map[string]map[string]struct{}),
		mu:         workflow.NewMutex(ctx),
		logger:     workflow.GetLogger(ctx),
	}
}
