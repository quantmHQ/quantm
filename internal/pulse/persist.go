// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024, 2025.
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

package pulse

import (
	"context"
	"fmt"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

const (
	statement__events__persist = `
INSERT INTO %s (
	version,
	id,
	parents,
	hook,
	scope,
	action,
	source,
	subject_id,
	subject_name,
	user_id,
	team_id,
	org_id,
	timestamp
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
)

// Persist persists an event to clickhouse, routing it to the appropriate activity handler based on the
// event's associated hook.  It's a workflow-scoped function, mandating execution immediately post-event creation.
func Persist[H events.Hook, P events.Payload](ctx workflow.Context, event *events.Event[H, P]) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	flat := event.Flatten()

	var future workflow.Future

	switch any(flat.Hook).(type) {
	case eventsv1.RepoHook:
		future = workflow.ExecuteActivity(ctx, PersistRepoEvent, flat)
	case eventsv1.ChatHook:
		future = workflow.ExecuteActivity(ctx, PersistChatEvent, flat)
	}

	return future.Get(ctx, nil)
}

// PersistRepoEvent persists a repo event to the database.
func PersistRepoEvent(ctx context.Context, flat events.Flat[eventsv1.RepoHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, flat.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Get().
		Connection().
		Exec(
			ctx,
			stmt,
			flat.Version,
			flat.ID,
			flat.Parents,
			flat.Hook.Number(),
			flat.Scope,
			flat.Action,
			flat.Source,
			flat.SubjectID,
			flat.SubjectName,
			flat.UserID,
			flat.TeamID,
			flat.OrgID,
			flat.Timestamp,
		)
}

// PersistChatEvent persists a chat event to the database.
func PersistChatEvent(ctx context.Context, flat events.Flat[eventsv1.ChatHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, flat.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Get().
		Connection().
		Exec(
			ctx,
			stmt,
			flat.Version,
			flat.ID,
			flat.Parents,
			flat.Hook.Number(),
			flat.Scope,
			flat.Action,
			flat.Source,
			flat.SubjectID,
			flat.SubjectName,
			flat.UserID,
			flat.TeamID,
			flat.OrgID,
			flat.Timestamp,
		)
}
