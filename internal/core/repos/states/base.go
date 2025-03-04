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
	"fmt"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	// Base represents the base state for repository workflows.  It encapsulates
	// core data and provides logging utilities.
	Base struct {
		Repo     *entities.Repo     `json:"repo"`      // Repository entity.
		ChatLink *entities.ChatLink `json:"chat_link"` // ChatLink entity.

		logger log.Logger // Workflow logger.
	}
)

// - private

// rx wraps workflow.ReceiveChannel.Receive, adding logging.  It receives a message
// from the specified Temporal channel. The target parameter must be a pointer to the
// data structure expected to be received.
func (state *Base) rx(ctx workflow.Context, ch workflow.ReceiveChannel, target any) {
	state.logger.Info(fmt.Sprintf("rx: %s", ch.Name()))
	ch.Receive(ctx, target)
}

// run wraps workflow.ExecuteActivity with logging with the default activity context. If you need to
// provide a custom context, use run_ex.
func (state *Base) run(ctx workflow.Context, action string, activity, event, result any, keyvals ...any) error {
	state.logger.Info(fmt.Sprintf("dispatch(%s): init ...", action), keyvals...)

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.ExecuteActivity(ctx, activity, event).Get(ctx, result); err != nil {
		state.logger.Error(fmt.Sprintf("dispatch(%s): error", action), keyvals...)
		return err
	}

	state.logger.Info(fmt.Sprintf("dispatch(%s): success", action), keyvals...)

	return nil
}

// - public

// RestartRecommended checks if the workflow should be continued as new.
func (state *Base) RestartRecommended(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// Init initializes the base state with the provided context.
func (state *Base) Init(ctx workflow.Context) {
	state.logger = workflow.GetLogger(ctx)
}
