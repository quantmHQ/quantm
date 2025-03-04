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

package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Trunk struct {
		*Base      `json:"base"`
		MergeQueue *Sequencer[int64, eventsv1.MergeQueue] `json:"merge_queue"`

		done     bool                   // done flag
		channel  workflow.Channel       // for cross loop communication
		inflight []*eventsv1.MergeQueue // in-flight merges
	}
)

// - queue process -

// OnMergeQueue is the signal handler for the merge queue.
func (state Trunk) OnMergeQueue(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		mq := &events.Event[eventsv1.RepoHook, eventsv1.MergeQueue]{}
		state.rx(ctx, rx, mq)

		if mq.Context.Action == events.EventActionRemoved {
			state.MergeQueue.Remove(ctx, mq.Payload.GetNumber())

			return
		}

		if mq.Payload.IsPriority {
			state.MergeQueue.Priority(ctx, mq.Payload.GetNumber(), mq.Payload)

			return
		}

		state.MergeQueue.Push(ctx, mq.Payload.GetNumber(), mq.Payload)
	}
}

// StartQueue is the main queue processing loop.
func (state *Trunk) StartQueue(ctx workflow.Context) {
	log := workflow.GetLogger(ctx)

	for state.Continue() && state.MergeQueue.Peek(ctx) != nil {
		next := state.MergeQueue.Pop(ctx) // next item

		// ahead of line testing
		// we rebase the changes from the branches that are being tested, this way, we can run tests on each.
		//
		// TODO: implement ahead of line testing
		// we will gather all the branches that are being tested and rebase them on top of the current branch.
		// this will allow us to run tests on each branch and merge them in order.
		//
		// we also will create a shadow branch that will be used to merge the changes into the main branch.
		log.Info("merge_queue: attempting ahead of line merge ...", "next", next, "in_prgress", state.inflight)
	}
}

func (state *Trunk) Continue() bool {
	return !state.done
}

func (state *Trunk) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
	state.MergeQueue.Init(ctx)
	state.channel = workflow.NewChannel(ctx)
}

func NewTrunk(repo *entities.Repo, chat *entities.ChatLink) *Trunk {
	return &Trunk{
		&Base{Repo: repo, ChatLink: chat},
		NewSequencer[int64, eventsv1.MergeQueue](),
		false,
		nil,
		make([]*eventsv1.MergeQueue, 0),
	}
}
