package mutex2_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/durable/mutex2"
)

type MutexTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

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
	handler := &mutex2.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "test-caller", RunID: "test-run"},
		},
		Timeout: time.Minute,
	}
	// Manual initialization to avoid calling workflow functions (NewMutex) outside of a workflow.
	state := &mutex2.MutexState{
		Status:  mutex2.MutexStatusAcquiring,
		Handler: handler,
		Timeout: handler.Timeout,
		Persist: true,
	}

	// We need to register the workflow to test it
	s.env.RegisterWorkflow(mutex2.MutexWorkflow)

	// Step 1: Verify it runs before timeout
	s.env.RegisterDelayedCallback(func() {
		s.False(s.env.IsWorkflowCompleted(), "Workflow should be running before idle timeout")
	}, mutex2.IdleTimeout-1*time.Second)

	// Step 2: Execute
	s.env.ExecuteWorkflow(mutex2.MutexWorkflow, state)

	// Step 3: Verify completion
	s.True(s.env.IsWorkflowCompleted(), "Workflow should shut down after idle timeout")
	s.NoError(s.env.GetWorkflowError())
}

func (s *MutexTestSuite) TestMutexWorkflow_ActivityResetIdle() {
	// Scenario: Workflow starts, receives Prepare/Acquire signals.
	// Expected: Idle timer should be reset, preventing shutdown at T=10m.
	resourceID := "test-resource-active"
	handler := &mutex2.Handler{
		ResourceID: resourceID,
		Info: &workflow.Info{
			WorkflowExecution: workflow.Execution{ID: "caller-workflow-id", RunID: "caller-run-id"},
		},
		Timeout: 10 * time.Minute,
	}
	// Manual initialization to avoid calling workflow functions (NewMutex) outside of a workflow.
	state := &mutex2.MutexState{
		Status:  mutex2.MutexStatusAcquiring,
		Handler: handler,
		Timeout: handler.Timeout,
		Persist: true,
		Pool: &mutex2.Pool{
			Data: map[string]time.Duration{
				"caller-workflow-id": handler.Timeout,
			},
		},
	}

	s.env.RegisterWorkflow(mutex2.MutexWorkflow)
	s.env.RegisterWorkflowWithOptions(func(ctx workflow.Context) error { return nil }, workflow.RegisterOptions{Name: "caller-workflow-id"})

	// Intercept SignalExternalWorkflow calls from Mutex -> Caller
	s.env.OnSignalExternalWorkflow(mock.Anything, "caller-workflow-id", "", mutex2.WorkflowSignalLocked.String(), mock.Anything).Return(nil)
	s.env.OnSignalExternalWorkflow(mock.Anything, "caller-workflow-id", "", mutex2.WorkflowSignalReleased.String(), mock.Anything).Return(nil)

	// 1. Advance time half-way to idle, then signal ACQUIRE
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex2.WorkflowSignalAcquire.String(), handler)
	}, mutex2.IdleTimeout/2)

	// 2. Advance time past the *original* idle timeout (T = Idle/2 + Idle/2 + 1s).
	// Since we signaled at T=Idle/2, the new timeout should be at T=Idle/2 + Idle = 1.5*Idle.
	// So at T=Idle + 1s, it should still be running.
	s.env.RegisterDelayedCallback(func() {
		s.False(s.env.IsWorkflowCompleted(), "Workflow should still be running after activity reset idle timer")
	}, mutex2.IdleTimeout+1*time.Second)

	// 3. Release the lock at T = Idle + 2s
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(mutex2.WorkflowSignalRelease.String(), handler)
	}, mutex2.IdleTimeout+2*time.Second)

	// 4. Verification happens after ExecuteWorkflow returns (implicit wait for idle)
	s.env.ExecuteWorkflow(mutex2.MutexWorkflow, state)

	s.True(s.env.IsWorkflowCompleted(), "Workflow should finally shut down after inactivity")
}

// -----------------------------------------------------------------------------
// Part 2: Client API Tests (New / OnAcquire)
// -----------------------------------------------------------------------------

// ConsumerWorkflowForTest is a fake workflow that uses the mutex2 library.
// We test the library by running this workflow.
func ConsumerWorkflowForTest(ctx workflow.Context, resourceID string) error {
	// 1. Initialize
	m, err := mutex2.New(ctx, mutex2.WithResourceID(resourceID))
	if err != nil {
		return err
	}

	// 2. Use OnAcquire
	return m.OnAcquire(ctx, func(lockCtx workflow.Context) {
		// Critical section simulation
		workflow.Sleep(lockCtx, 5*time.Second)
	})
}

func (s *MutexTestSuite) TestClient_OnAcquire_Success() {
	resourceID := "test-resource-client"
	mutexWorkflowID := "ai.ctrlplane.mutex.resource-v2." + resourceID

	// Mock the PrepareMutexActivity
	s.env.OnActivity(mutex2.PrepareMutexActivity, mock.Anything, mock.Anything).Return(
		&workflow.Execution{ID: mutexWorkflowID, RunID: "mutex-run-id"},
		nil,
	)

	// Mock signals sent FROM Client TO MutexWorkflow
	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "", mutex2.WorkflowSignalAcquire.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			s.env.SignalWorkflow(mutex2.WorkflowSignalLocked.String(), true)
		}).
		Return(nil)

	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "", mutex2.WorkflowSignalRelease.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			s.env.SignalWorkflow(mutex2.WorkflowSignalReleased.String(), true) // Not orphan
		}).
		Return(nil)

	s.env.ExecuteWorkflow(ConsumerWorkflowForTest, resourceID)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *MutexTestSuite) TestClient_OnAcquire_PanicSafety() {
	resourceID := "test-resource-panic"
	mutexWorkflowID := "ai.ctrlplane.mutex.resource-v2." + resourceID

	s.env.OnActivity(mutex2.PrepareMutexActivity, mock.Anything, mock.Anything).Return(
		&workflow.Execution{ID: mutexWorkflowID, RunID: "mutex-run-id"},
		nil,
	)

	// Expect Acquire signal
	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "", mutex2.WorkflowSignalAcquire.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			s.env.SignalWorkflow(mutex2.WorkflowSignalLocked.String(), true)
		}).
		Return(nil)

	// Expect Release signal DESPITE panic
	s.env.OnSignalExternalWorkflow(mock.Anything, mutexWorkflowID, "", mutex2.WorkflowSignalRelease.String(), mock.Anything).
		Run(func(args mock.Arguments) {
			s.env.SignalWorkflow(mutex2.WorkflowSignalReleased.String(), true)
		}).
		Return(nil)

	// Define a workflow that panics inside the lock
	panicWorkflow := func(ctx workflow.Context) error {
		m, _ := mutex2.New(ctx, mutex2.WithResourceID(resourceID))

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
