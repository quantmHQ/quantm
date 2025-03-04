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
	"github.com/google/uuid"
)

type (
	// Subject represents the entity within the system that is the subject of an event.
	//
	// It encapsulates:
	//   - ID: The unique identifier of the entity i.e. the primary key within its respective database table.
	//   - Name: The name of the entity's corresponding database table. This provides context for the event's subject.
	//     For instance, an event related to a branch would have "repos" as the subject name, as branches are associated
	//     with repositories.
	//   - TeamID: The unique identifier of the team to which this entity belongs. This allows for team-based filtering
	//     and organization
	//     of events.
	Subject struct {
		Name   string    `json:"name"`    // Name of the database table.
		ID     uuid.UUID `json:"id"`      // ID is the ID of the subject.
		OrgID  uuid.UUID `json:"org_id"`  // OrgID is the ID of the organization that the subject belongs to.
		TeamID uuid.UUID `json:"team_id"` // Team ID of the subject's team in the organization. It can be null uuid.
		UserID uuid.UUID `json:"user_id"` // UserID is the ID of the user that the subject belongs to. It can be null uuid.
	}
)

const (
	SubjectNameRepos = "repos"
	SubjectNameChat  = "chat"
)
