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
	// Flat is the flat structure of an event for time series databases. It does not
	// contain payload data, only metadata.
	Flat[H Hook] struct {
		Version     EventVersion `json:"version"`      // Version is the version of the event.
		ID          uuid.UUID    `json:"id"`           // ID is the ID of the event.
		Parents     []uuid.UUID  `json:"parents"`      // ParentID is the ID of the parent event.
		Hook        H            `json:"provider"`     // Provider is the provider of the event.
		Scope       Scope        `json:"scope"`        // Scope is the scope of the event.
		Action      Action       `json:"action"`       // Action is the action of the event.
		Source      string       `json:"source"`       // Source is the source of the event. For every hook it will be in different format.
		SubjectID   uuid.UUID    `json:"subject_id"`   // SubjectID is the ID of the subject.
		SubjectName string       `json:"subject_name"` // SubjectName is the name of the subject.
		UserID      uuid.UUID    `json:"user_id"`      // UserID is the ID of the user that the subject belongs to. Can be empty.
		TeamID      uuid.UUID    `json:"team_id"`      // TeamID is the ID of the team that the subject belongs to. Can be empty.
		OrgID       uuid.UUID    `json:"org_id"`       // OrgID is the ID of the organization that the subject belongs to.
		Timestamp   time.Time    `json:"timestamp"`    // Timestamp is the timestamp of the event.
	}
)
