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

package states_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/states"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	SequencerTestSuite struct {
		suite.Suite
		testsuite.WorkflowTestSuite

		env *testsuite.TestWorkflowEnvironment
	}
)

const (
	PushSignal    queues.Signal = "push"
	PopSignal     queues.Signal = "pop"
	PromoteSignal queues.Signal = "promote"
	DemoteSignal  queues.Signal = "demote"
	DoneSignal    queues.Signal = "done"
)

const (
	Peek          queues.Query = "peek"
	PositionQuery queues.Query = "position"
	LengthQuery   queues.Query = "length"
	AllQuery      queues.Query = "all"
)

func (s *SequencerTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *SequencerTestSuite) Test_001_Push() {
	pr := &eventsv1.PullRequest{Number: 1}

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(PushSignal.String(), pr)
	}, time.Millisecond*50)

	s.env.ExecuteWorkflow(SequencerTestWorkflow)

	{
		result := &eventsv1.PullRequest{}

		ptr, err := s.env.QueryWorkflow(Peek.String())
		if s.NoError(err) {
			_ = ptr.Get(&result)
			s.Equal(pr.Number, result.Number)
		}
	}

	s.env.AssertExpectations(s.T())

	{
		result := 0

		ptr, err := s.env.QueryWorkflow(LengthQuery.String())
		if s.NoError(err) {
			_ = ptr.Get(&result)
			s.Equal(1, result)
		}
	}

	s.env.AssertExpectations(s.T())
}

func SequencerTestWorkflow(ctx workflow.Context) error {
	done := false
	seq := states.NewSequencer[int64, eventsv1.PullRequest]()
	seq.Init(ctx)
	selector := workflow.NewSelector(ctx)

	// Push Signal Handler
	{
		ch := workflow.GetSignalChannel(ctx, PushSignal.String())
		selector.AddReceive(ch, func(rx workflow.ReceiveChannel, more bool) {
			var item eventsv1.PullRequest

			rx.Receive(ctx, &item)
			seq.Push(ctx, item.Number, &item)
		})
	}

	// Pop Signal Handler
	{
		ch := workflow.GetSignalChannel(ctx, PopSignal.String())
		selector.AddReceive(ch, func(rx workflow.ReceiveChannel, more bool) {
			seq.Pop(ctx)
		})
	}

	// Promote Signal Handler
	{
		ch := workflow.GetSignalChannel(ctx, PromoteSignal.String())
		selector.AddReceive(ch, func(rx workflow.ReceiveChannel, more bool) {
			var key int64

			rx.Receive(ctx, &key)
			seq.Promote(ctx, key)
		})
	}

	// Demote Signal Handler
	{
		ch := workflow.GetSignalChannel(ctx, DemoteSignal.String())
		selector.AddReceive(ch, func(rx workflow.ReceiveChannel, more bool) {
			var key int64

			rx.Receive(ctx, &key)
			seq.Demote(ctx, key)
		})
	}

	// Done Signal Handler (optional - for workflow completion)
	{
		ch := workflow.GetSignalChannel(ctx, DoneSignal.String())
		selector.AddReceive(ch, func(rx workflow.ReceiveChannel, more bool) {
			rx.Receive(ctx, &done)
		})
	}

	// Peek Query Handler
	_ = workflow.SetQueryHandler(ctx, Peek.String(), func() (*eventsv1.PullRequest, error) {
		return seq.Peek(ctx), nil
	})

	// Position Query Handler
	_ = workflow.SetQueryHandler(ctx, PositionQuery.String(), func(key int64) (int, error) {
		return seq.Position(ctx, key), nil
	})

	// Length Query Handler
	_ = workflow.SetQueryHandler(ctx, LengthQuery.String(), func() (int, error) {
		return seq.Length(ctx), nil
	})

	// All Query Handler
	_ = workflow.SetQueryHandler(ctx, AllQuery.String(), func() ([]*eventsv1.PullRequest, error) {
		return seq.All(ctx), nil
	})

	for !done {
		selector.Select(ctx)
	}

	return nil
}

func TestSequenceSuite(t *testing.T) {
	suite.Run(t, new(SequencerTestSuite))
}
