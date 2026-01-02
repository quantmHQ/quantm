package mutex_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/durable/mutex"
)

type (
	MutexTestSuite struct {
		suite.Suite
		testsuite.WorkflowTestSuite
		env *testsuite.TestWorkflowEnvironment
	}
)

func TestMutexTestSuite(t *testing.T) {
	suite.Run(t, new(MutexTestSuite))
}

func (s *MutexTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *MutexTestSuite) TearDownTest() {
	s.env.AssertExpectations(s.T())
}

// -----------------------------------------------------------------------------
// Part 1: MutexWorkflow Tests (The Engine)
// -----------------------------------------------------------------------------

func (s *MutexTestSuite) TestMutexWorkflow_IdleShutdown() {
	// Scenario: Workflow starts, no signals received.
	// Expected: Should shut down after IdleTimeout (10m).
	resourceID := "test-resource-idle"
	handler := &mutex.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "test-caller", RunID: "test-run"},
		},
		Timeout: time.Minute,
	}
	// Manual initialization to avoid calling workflow functions (NewMutex) outside of a workflow.
	state := &mutex.MutexState{
		Status:  mutex.MutexStatusAcquiring,
		Handler: handler,
		Timeout: handler.Timeout,
		Persist: true,
	}

	// We need to register the workflow to test it
	s.env.RegisterWorkflow(mutex.MutexWorkflow)

	// Step 1: Verify it runs before timeout
	s.env.RegisterDelayedCallback(func() {
		s.False(s.env.IsWorkflowCompleted(), "Workflow should be running before idle timeout")
	}, mutex.IdleTimeout-1*time.Second)

	// Step 2: Execute
	s.env.ExecuteWorkflow(mutex.MutexWorkflow, state)

	// Step 3: Verify completion
	s.True(s.env.IsWorkflowCompleted(), "Workflow should shut down after idle timeout")
	s.NoError(s.env.GetWorkflowError())
}

func (s *MutexTestSuite) TestMutexWorkflow_ActivityResetIdle() {
	// Scenario: Workflow starts, receives Acquire signals.
	// Expected: Idle timer should be reset, preventing shutdown at T=10m.
	resourceID := "test-resource-active"
	handler := &mutex.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "caller-workflow-id", RunID: "caller-run-id"},
		},
		Timeout: 10 * time.Minute,
	}
	// Manual initialization to avoid calling workflow functions (NewMutex) outside of a workflow.
	state := &mutex.MutexState{
		Status:  mutex.MutexStatusAcquiring,
		Handler: handler,
		Timeout: handler.Timeout,
		Persist: true,
	}

	s.env.RegisterWorkflow(mutex.MutexWorkflow)
	s.env.RegisterWorkflowWithOptions(func(ctx workflow.Context) error { return nil }, workflow.RegisterOptions{Name: "caller-workflow-id"})

	// Intercept SignalExternalWorkflow calls from Mutex -> Caller
	s.env.OnSignalExternalWorkflow(mock.Anything, "caller-workflow-id", "caller-run-id", mutex.WorkflowSignalLocked.String(), mock.Anything).Return(nil)   // nolint
	s.env.OnSignalExternalWorkflow(mock.Anything, "caller-workflow-id", "caller-run-id", mutex.WorkflowSignalReleased.String(), mock.Anything).Return(nil) // nolint

	// 1. Advance time half-way to idle, then signal ACQUIRE
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalAcquire.String(), handler)
	}, mutex.IdleTimeout/2)

	// 2. Advance time past the *original* idle timeout (T = Idle/2 + Idle/2 + 1s).
	// Since we signaled at T=Idle/2, the new timeout should be at T=Idle/2 + Idle = 1.5*Idle.
	// So at T=Idle + 1s, it should still be running.
	s.env.RegisterDelayedCallback(func() {
		s.False(s.env.IsWorkflowCompleted(), "Workflow should still be running after activity reset idle timer")
	}, mutex.IdleTimeout+1*time.Second)

	// 3. Release the lock at T = Idle + 2s
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalRelease.String(), handler)
	}, mutex.IdleTimeout+2*time.Second)

	// 4. Verification happens after ExecuteWorkflow returns (implicit wait for idle)
	s.env.ExecuteWorkflow(mutex.MutexWorkflow, state)

	s.True(s.env.IsWorkflowCompleted(), "Workflow should finally shut down after inactivity")
}

func (s *MutexTestSuite) TestMutexWorkflow_Contention() {
	// Scenario:
	// 1. Client A acquires.
	// 2. Client B tries to acquire (should block).
	// 3. Client A releases.
	// 4. Client B should get the lock.
	// 5. Client B releases.
	resourceID := "test-resource-contention"
	handlerA := &mutex.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "client-A", RunID: "run-A"},
		},
		Timeout: 10 * time.Minute,
	}
	handlerB := &mutex.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "client-B", RunID: "run-B"},
		},
		Timeout: 10 * time.Minute,
	}

	// Manual initialization
	state := &mutex.MutexState{
		Status:  mutex.MutexStatusAcquiring,
		Handler: handlerA, // Initial handler for logging context
		Persist: true,
	}

	s.env.RegisterWorkflow(mutex.MutexWorkflow)

	// Track event order
	var eventLog []string

	// Mocks for Client A
	s.env.OnSignalExternalWorkflow(mock.Anything, "client-A", "run-A", mutex.WorkflowSignalLocked.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			eventLog = append(eventLog, "A-Locked")
		}).Return(nil)
	s.env.OnSignalExternalWorkflow(mock.Anything, "client-A", "run-A", mutex.WorkflowSignalReleased.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			eventLog = append(eventLog, "A-Released")
		}).Return(nil)

	// Mocks for Client B
	s.env.OnSignalExternalWorkflow(mock.Anything, "client-B", "run-B", mutex.WorkflowSignalLocked.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			eventLog = append(eventLog, "B-Locked")
		}).Return(nil)
	s.env.OnSignalExternalWorkflow(mock.Anything, "client-B", "run-B", mutex.WorkflowSignalReleased.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			eventLog = append(eventLog, "B-Released")
		}).Return(nil)

	// Sequence

	// T+1s: Acquire A
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalAcquire.String(), handlerA)
	}, 1*time.Second)

	// T+2s: Acquire B (Should be buffered)
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalAcquire.String(), handlerB)
	}, 2*time.Second)

	// T+4s: Release A
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalRelease.String(), handlerA)
	}, 4*time.Second)

	// T+6s: Release B
	// Note: We need enough delay to ensure B processes its Locked signal
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalRelease.String(), handlerB)
	}, 6*time.Second)

	// T+15m: Ensure workflow eventually shuts down to finish test
	// (Idle timeout is 10m, so it should shut down after last activity)

	s.env.ExecuteWorkflow(mutex.MutexWorkflow, state)

	// Verify Order
	expected := []string{"A-Locked", "A-Released", "B-Locked", "B-Released"}
	s.Equal(expected, eventLog)
}

// -----------------------------------------------------------------------------
// Part 2: Client API Tests (New / OnAcquire)
// -----------------------------------------------------------------------------

// ConsumerWorkflowForTest is a fake workflow that uses the mutex library.
// We test the library by running this workflow.
func ConsumerWorkflowForTest(ctx workflow.Context, resourceID string) error {
	// 1. Initialize
	m, err := mutex.New(ctx, mutex.WithResourceID(resourceID))
	if err != nil {
		return err
	}

	// 2. Use OnAcquire
	return m.OnAcquire(ctx, func(lockCtx workflow.Context) {
		_ = workflow.Sleep(lockCtx, 5*time.Second)
	})
}

func (s *MutexTestSuite) TestClient_OnAcquire_Success() {
	resourceID := "test-resource-client"
	mutexWorkflowID := "ai.ctrlplane.mutex.resource-v2." + resourceID

	// Mock the AcquireMutexActivity
	s.env.OnActivity(mutex.AcquireMutexActivity, mock.Anything, mock.Anything).Return(
		&workflow.Execution{ID: mutexWorkflowID, RunID: "mutex-run-id"},
		nil,
	)

	// Mock signals sent FROM Client TO MutexWorkflow (Acquire is sent via Activity, so only signals used in wait logic here)
	// Actually, Acquire is sent via SignalExternalWorkflow inside AcquireMutexActivity helper in real usage?
	// No, AcquireMutexActivity calls SignalWithStart.
	// The client workflow (ConsumerWorkflowForTest) -> m.OnAcquire -> calls AcquireMutexActivity.
	// Then m.OnAcquire waits for "mutex__locked".

	// We need to simulate the MutexWorkflow signaling BACK to the client.
	// Since we are running the ConsumerWorkflow, we need to mock the "environment" responding to it.
	// The Consumer calls Activity, then waits for Signal "mutex__locked".
	// We can use RegisterDelayedCallback to simulate the signal arrival.

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalLocked.String(), true)
	}, 1*time.Second)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalReleased.String(), true)
	}, 7*time.Second) // 1s (start) + 5s (sleep) + buffer

	// Also mock the Release signal from Client -> Mutex
	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "mutex-run-id", mutex.WorkflowSignalRelease.String(), mock.Anything).
		Return(nil)

	s.env.ExecuteWorkflow(ConsumerWorkflowForTest, resourceID)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *MutexTestSuite) TestClient_OnAcquire_PanicSafety() {
	resourceID := "test-resource-panic"
	mutexWorkflowID := "ai.ctrlplane.mutex.resource-v2." + resourceID

	s.env.OnActivity(mutex.AcquireMutexActivity, mock.Anything, mock.Anything).Return(
		&workflow.Execution{ID: mutexWorkflowID, RunID: "mutex-run-id"},
		nil,
	)

	// Simulate lock acquired signal
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalLocked.String(), true)
	}, 1*time.Second)

	// Simulate release confirmation signal
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex.WorkflowSignalReleased.String(), true)
	}, 2*time.Second)

	// Expect Release signal DESPITE panic
	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "mutex-run-id", mutex.WorkflowSignalRelease.String(), mock.Anything).
		Return(nil)

	// Define a workflow that panics inside the lock
	panicWorkflow := func(ctx workflow.Context) error {
		m, _ := mutex.New(ctx, mutex.WithResourceID(resourceID))

		return m.OnAcquire(ctx, func(c workflow.Context) {
			panic("business logic failure")
		})
	}

	s.env.ExecuteWorkflow(panicWorkflow)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	// Check for panic error string safely
	s.True(strings.Contains(err.Error(), "business logic failure"), "Error should contain panic message")
}
