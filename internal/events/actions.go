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

package events

type (
	// Action is the action of the event.
	Action string
)

const (
	ActionCreated      Action = "created"   // ActionCreated indicates the initial creation.
	ActionDeleted      Action = "deleted"   // ActionDeleted indicates the removal.
	ActionUpdated      Action = "updated"   // ActionUpdated indicates that something has been modified.
	ActionForced       Action = "forced"    // ActionForced indicates an action was applied regardless of normal constraints.
	ActionReopened     Action = "reopened"  // ActionReopened indicates a previously closed item has been reopened.
	ActionClosed       Action = "closed"    // ActionClosed indicates an item or process has reached a terminal/inactive state.
	ActionStarted      Action = "started"   // ActionStarted indicates the start of a process or task.
	ActionCompleted    Action = "completed" // ActionCompleted indicates a process, task, or work item was successfully finished.
	ActionDismissed    Action = "dismissed" // ActionDismissed indicates a user dismissed an item or alert.
	ActionFailure      Action = "failure"   // ActionAbandoned indicates a failure or abandonment of a process or task.
	EventActionAdded   Action = "added"     // EventActionAdded indicates something was added to something else.
	EventActionRemoved Action = "removed"   // EventActionRemoved indicates something was removed from something else.
	ActionRequested    Action = "requested" // ActionRequested indicates a request for an action, approval, or resource was initiated.
)

// String returns the string representation of the EventAction.
func (a Action) String() string { return string(a) }

// SetActionCreated sets the action of the Event to ActionCreated.
func (e *Event[T, P]) SetActionCreated() {
	e.Context.Action = ActionCreated
}

// SetActionDeleted sets the action of the Event to ActionDeleted.
func (e *Event[T, P]) SetActionDeleted() {
	e.Context.Action = ActionDeleted
}

// SetActionUpdated sets the action of the Event to ActionUpdated.
func (e *Event[T, P]) SetActionUpdated() {
	e.Context.Action = ActionUpdated
}

// SetActionDismissed sets the action of the Event to ActionDismissed.
func (e *Event[T, P]) SetActionDismissed() {
	e.Context.Action = ActionDismissed
}
