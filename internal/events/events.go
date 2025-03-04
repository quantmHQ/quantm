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

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Event represents an event.  Events are created using the New function and must be persisted immediately.
	Event[H Hook, P Payload] struct {
		Version   EventVersion `json:"version"`   // Version is the version of the event.
		ID        uuid.UUID    `json:"id"`        // ID is the ID of the event.
		Timestamp time.Time    `json:"timestamp"` // Timestamp is the Event occurrence time.
		Context   Context[H]   `json:"context"`   // Context is the context of the event.
		Subject   Subject      `json:"subject"`   // Subject is the subject of the event.
		Payload   *P           `json:"payload"`   // Payload is the payload of the event.
	}
)

// SetParents sets the adds the given id to parents.
func (e *Event[H, P]) SetParents(id ...uuid.UUID) *Event[H, P] {
	e.Context.Parents = append(e.Context.Parents, id...)
	return e
}

// SetHook sets the hook of the event.
func (e *Event[H, P]) SetHook(hook H) *Event[H, P] {
	e.Context.Hook = hook
	return e
}

// SetScope sets the scope of the event.
func (e *Event[H, P]) SetScope(scope Scope) *Event[H, P] {
	e.Context.Scope = scope
	return e
}

// SetAction sets the action of the event.
func (e *Event[H, P]) SetAction(action Action) *Event[H, P] {
	e.Context.Action = action
	return e
}

// SetSource sets the source of the event.
func (e *Event[H, P]) SetSource(source string) *Event[H, P] {
	e.Context.Source = source
	return e
}

// SetSubjectID sets the subject ID of the event.
func (e *Event[H, P]) SetSubjectID(id uuid.UUID) *Event[H, P] {
	e.Subject.ID = id
	return e
}

// SetSubjectName sets the subject name of the event.
func (e *Event[H, P]) SetSubjectName(name string) *Event[H, P] {
	e.Subject.Name = name
	return e
}

// SetOrg sets the organization ID of the event.
func (e *Event[H, P]) SetOrg(id uuid.UUID) *Event[H, P] {
	e.Subject.OrgID = id
	return e
}

// SetTeam sets the team ID of the event.
func (e *Event[H, P]) SetTeam(id uuid.UUID) *Event[H, P] {
	e.Subject.TeamID = id
	return e
}

// SetUser sets the user ID of the event.
func (e *Event[H, P]) SetUser(id uuid.UUID) *Event[H, P] {
	e.Subject.UserID = id
	return e
}

func (e *Event[H, P]) SetContext(ctx Context[H]) *Event[H, P] {
	e.Context = ctx
	return e
}

func (e *Event[H, P]) SetSubject(subject Subject) *Event[H, P] {
	e.Subject = subject
	return e
}

// SetPayload sets the payload of the event.
func (e *Event[H, P]) SetPayload(payload *P) *Event[H, P] {
	e.Payload = payload
	return e
}

// Flatten flattens the event into a simpler structure.
func (e *Event[H, P]) Flatten() *Flat[H] {
	return &Flat[H]{
		Version:     e.Version,
		ID:          e.ID,
		Timestamp:   e.Timestamp,
		Parents:     e.Context.Parents,
		Hook:        e.Context.Hook,
		Scope:       e.Context.Scope,
		Action:      e.Context.Action,
		Source:      e.Context.Source,
		SubjectID:   e.Subject.ID,
		SubjectName: e.Subject.Name,
		OrgID:       e.Subject.OrgID,
		TeamID:      e.Subject.TeamID,
		UserID:      e.Subject.UserID,
	}
}

// Next creates a new event based on the provided event, scope, and action.
func Next[H Hook, F Payload, T Payload](event *Event[H, F], scope Scope, action Action) *Event[H, T] {
	return NextWithHook[H, H, F, T](event, event.Context.Hook, scope, action)
}

// NextWithHook transforms the event from one hook type to another while keeping the payload.
func NextWithHook[H1 Hook, H2 Hook, F Payload, T Payload](event *Event[H1, F], hook H2, scope Scope, action Action) *Event[H2, T] {
	ctx := Context[H2]{
		Parents: append(event.Context.Parents, event.ID),
		Hook:    hook,
		Scope:   scope,
		Action:  action,
		Source:  event.Context.Source,
	}

	return New[H2, T]().SetContext(ctx).SetSubject(event.Subject)
}

// New creates a new event with default values.
func New[H Hook, P Payload]() *Event[H, P] {
	event := &Event[H, P]{
		Version:   EventVersionDefault,
		ID:        MustUUID(),
		Timestamp: time.Now(),
		Context:   Context[H]{Parents: make([]uuid.UUID, 0)},
		Subject:   Subject{},
	}

	return event
}
