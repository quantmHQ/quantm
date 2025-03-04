// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2025.
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

package activities

import (
	"context"
	"log/slog"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Notify sends chat notifications.
	Notify struct{}
)

// LinesExceeded notifies a chat service of exceeded lines. It uses the context and event to dispatch a
// notification via a chat hook. Returns error if notification fails, logging a warning.
func (n *Notify) LinesExceeded(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyLinesExceed(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

// MergeConflict notifies a chat service of a merge conflict. It uses the context and event to dispatch a
// notification via a chat hook. Returns error if notification fails, logging a warning.
func (n *Notify) MergeConflict(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyMergeConflict(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}
