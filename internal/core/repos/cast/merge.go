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

package cast

import (
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// MergeConflictEvent converts a Push event to a diff event.
func RebaseEventToMergeConflictEvent(
	rebase *events.Event[eventsv1.RepoHook, eventsv1.Rebase],
	hook int32,
	payload *eventsv1.Merge,
) *events.Event[eventsv1.ChatHook, eventsv1.Merge] {
	return events.NextWithHook[eventsv1.RepoHook, eventsv1.ChatHook, eventsv1.Rebase, eventsv1.Merge](
		rebase,
		eventsv1.ChatHook(hook),
		events.ScopeMerge,
		events.ActionFailure,
	).SetPayload(payload)
}
