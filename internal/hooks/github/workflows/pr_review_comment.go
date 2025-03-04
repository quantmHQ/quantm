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
	"github.com/google/uuid"
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/cast"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

func PullRequestReviewComment(ctx workflow.Context, prrc *defs.PrReviewComment) error {
	acts := &activities.PullRequestReviewComment{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PrReviewCommentToProto(prrc)
	hydrated := &defs.HydratedRepoEvent{}

	email := ""
	if prrc.GetSenderEmail() != nil {
		email = *prrc.GetSenderEmail()
	}

	{
		payload := &defs.HydratedRepoEventPayload{
			RepoID:         prrc.GetRepositoryID(),
			InstallationID: prrc.GetInstallationID(),
			Email:          email,
			Branch:         repos.BranchNameFromRef(prrc.GetHeadBranch()),
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPREvent, payload).Get(ctx, hydrated); err != nil {
			return err
		}
	}

	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequestReviewComment]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePr).
		SetSource(hydrated.GetRepoUrl()).
		SetOrg(hydrated.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(hydrated.GetRepoID()).
		SetPayload(&proto)

	switch prrc.GetAction() {
	case "created":
		event.SetActionCreated()
	case "edited":
		event.SetActionUpdated()
	case "deleted":
		event.SetActionDismissed()
	default:
		return nil
	}

	if hydrated.GetParentID() != uuid.Nil {
		event.SetParents(hydrated.GetParentID())
	}

	if hydrated.GetTeam() != nil {
		event.SetTeam(hydrated.GetTeamID())
	}

	if hydrated.GetUser() != nil {
		event.SetUser(hydrated.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequestReviewComment]{
		Event: event, Meta: hydrated, Signal: repos.SignalPullRequestReviewComment,
	}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}
