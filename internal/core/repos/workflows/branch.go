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

package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

func Branch(ctx workflow.Context, state *states.Branch) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	// - activity monitors -

	state.PullRequestMonitor(ctx)
	state.StaleMonitor(ctx)

	// - signal handlers -

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	rebase := workflow.GetSignalChannel(ctx, defs.SignalRebase.String())
	selector.AddReceive(rebase, state.OnRebase(ctx))

	label := workflow.GetSignalChannel(ctx, defs.SignalPullRequestLabel.String())
	selector.AddReceive(label, state.OnLabel(ctx))

	prr := workflow.GetSignalChannel(ctx, defs.SignalPullRequestReview.String())
	selector.AddReceive(prr, state.OnPrReview(ctx))

	prrc := workflow.GetSignalChannel(ctx, defs.SignalPullRequestReviewComment.String())
	selector.AddReceive(prrc, state.OnPRReviewComment(ctx))

	// - event loop -

	for !state.ExitLoop(ctx) {
		selector.Select(ctx)
	}

	// - exit or continue -

	if state.RestartRecommended(ctx) {
		return workflow.NewContinueAsNewError(ctx, Branch, state)
	}

	return nil
}
