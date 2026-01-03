# Distributed Mutex Refactor Specification

## 1. Executive Summary
The current `internal/durable/mutex` implementation requires a verbose and error-prone four-step manual lifecycle (`Prepare`, `Acquire`, `Release`, `Cleanup`). This specification details the refactoring of the package to introduce a safer, functional API that manages the lifecycle automatically, preventing deadlocks and zombie workflows.

**Strategy:** To ensure stability and zero disruption, we will implement this new architecture in a parallel package `internal/durable/mutex2`. Only upon full verification and satisfaction will the new implementation replace the existing one.

**Primary Goal:** Reduce cognitive load and bug surface area by automating lock lifecycle and cleanup, utilizing a non-destructive, side-by-side implementation strategy.

## 2. Problem Statement
*   **Verbosity:** Consumers must manually orchestrate `Prepare`, `Acquire`, `Release`, and `Cleanup`.
*   **Brittleness:** Missing `defer Release()` causes deadlocks.
*   **Zombie Workflows:** Missing `Cleanup()` leaves the underlying `MutexWorkflow` running indefinitely, consuming resources.
*   **Error Mixing:** Infrastructure errors (lock acquisition) are currently mixed with business logic errors.

## 3. Objectives
1.  **Atomic Lifecycle:** Replace the manual sequence with a single functional entry point `OnAcquire`.
2.  **Safe Initialization:** Ensure the workflow exists and is ready upon instantiation.
3.  **Automatic Resource Management:** The `MutexWorkflow` must detect idleness and shut itself down, removing the need for client-side `Cleanup`.
4.  **Separation of Concerns:** Distinct handling for locking mechanics versus business logic.
5.  **Non-Destructive Evolution:** Implement the new logic in a fresh namespace (`mutex2`) to allow for iterative improvement and rigorous testing without destabilizing existing consumers.

## 4. Technical Specification

### 4.1. Package Structure
The new implementation will reside in `internal/durable/mutex2`.
*   It will mirror the necessary internal structure of the original package but with the new API and workflow logic.
*   It will **not** share state with the legacy `mutex` package to ensure isolation.

### 4.2. API Changes (Public Interface)

#### 4.2.1. Initialization (`New`)
The `New` constructor shall be promoted from a simple factory to an initializing constructor.

*   **Signature:** `func New(ctx workflow.Context, opts ...Option) (Mutex, error)`
*   **Behavior:**
    1.  Validates configuration (Resource ID, etc.).
    2.  Executes the `Prepare` logic (sending `SignalWithStart` to the underlying workflow) immediately.
    3.  Returns the `Mutex` interface only upon successful initialization of the backend workflow.
    4.  Returns an error if the workflow cannot be reached or started.

#### 4.2.2. Execution (`OnAcquire`)
A new method `OnAcquire` shall replace the explicit `Acquire` and `Release` methods in the public interface.

*   **Signature:** `OnAcquire(ctx workflow.Context, fn func(workflow.Context))`
*   **Parameters:**
    *   `ctx`: The parent workflow context.
    *   `fn`: A closure containing the critical section code. It accepts a `workflow.Context` which is valid *only* for the duration of the lock.
*   **Behavior:**
    1.  **Acquisition:** Blocks until the distributed lock is acquired or the context is cancelled/timed out.
    2.  **Execution:** Invokes the provided `fn` closure.
    3.  **Release:** Automatically releases the lock when `fn` returns, regardless of success or panic (via `defer`).
    4.  **Error Handling:** Returns an error *only* if the locking mechanism fails (e.g., timeout, workflow unreachable). Logic errors within `fn` are handled internally by the closure or through side effects; they do not propagate through `OnAcquire`'s return value.

### 4.3. Workflow Logic Updates (`MutexWorkflow`)
The new workflow in `mutex2` will incorporate automatic lifecycle management.

#### 4.3.1. Automatic Idle Shutdown
The `MutexWorkflow` event loop shall be modified to support self-termination.

*   **Idle Detection:** The main event loop must track the number of active locks and pending requests (the "Queue").
*   **Timer Strategy:**
    *   The workflow must utilize the `internal/durable/periodic` package to manage the idle timeout.
    *   Create a timer using `periodic.New(ctx, IdleTimeout)`.
    *   When the queue is empty, start the timer loop.
    *   If a new `Acquire` signal is received, the timer should be reset or adjusted using `periodic.Interval` methods (e.g., `Restart` or `Stop`).
    *   If the timer fires (tick completes) before any new activity, the workflow gracefully terminates.
*   **Removal of Explicit Cleanup:** The `WorkflowSignalCleanup` and associated logic can be deprecated or repurposed as an immediate "force kill" rather than a required lifecycle step.

### 4.4. Interface Definition
The `Mutex` interface shall be updated for the new package:

```go
type Mutex interface {
    // OnAcquire blocks until the lock is acquired, executes fn, and then releases the lock.
    // It returns an error if the lock cannot be acquired or the context is cancelled.
    OnAcquire(ctx workflow.Context, fn func(workflow.Context)) error
}
```

(Note: `Prepare`, `Acquire`, `Release`, and `Cleanup` are removed from the public interface or made internal/private).

## 5. Implementation Plan (Iterative)

The implementation will proceed in distinct, verifiable phases within the `internal/durable/mutex2` namespace.

### Phase 1: Scaffolding & Setup
1.  Create `internal/durable/mutex2`.
2.  Copy core data structures (`MutexState`, `Pool`, `Handler`) from `mutex` to `mutex2` as a starting point.
3.  Rename/Namespacing: Ensure workflow names/IDs used by `mutex2` are distinct (e.g., prefix with `v2-` or similar) to prevent collision with running v1 workflows during testing.

### Phase 2: Workflow Logic (The Engine)
1.  Implement `MutexWorkflow` in `mutex2`.
2.  **Crucial:** Implement the "Idle Timer" loop logic.
    *   Logic: `Loop { Select { Case SignalAcquire: Handle; Case Timer(Idle): Return } }`.
3.  Verify the workflow terminates automatically after the timeout when no signals are received.

### Phase 3: Client API (The Interface)
1.  Implement `mutex2.New(...)`.
    *   Must perform `SignalWithStart` to ensure the V2 workflow is running.
2.  Implement `mutex2.OnAcquire(...)`.
    *   Logic:
        *   `defer Signal(Release)`
        *   `Signal(Acquire).Get(Future)` // Block
        *   `fn(childCtx)`
3.  Ensure `OnAcquire` correctly handles panics in `fn` (releasing the lock).

### Phase 4: Verification & Cutover
1.  Write comprehensive tests for `mutex2` covering:
    *   Basic acquisition/release.
    *   Concurrent contention (queueing).
    *   Timeout/Cancellation.
    *   **Idle Shutdown verification.**
2.  Once satisfied:
    *   Delete `internal/durable/mutex`.
    *   Rename `internal/durable/mutex2` to `internal/durable/mutex`.
    *   Update documentation `doc.go` with the new usage patterns.

### Usage Example (Target)

```go
// Intended usage after migration
lock, err := mutex.New(ctx, mutex.WithResourceID("..."))
if err != nil { return err }

err = lock.OnAcquire(ctx, func(lockCtx workflow.Context) {
    // Critical section
    // Lock is held here.
})
// Lock is guaranteed released here.
```