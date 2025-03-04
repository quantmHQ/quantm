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
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

// // SyncRepos synchronizes repositories for a given installation.
//
// This workflow handles the addition and removal of repositories
// from a GitHub installation. It retrieves the installation details,
// then iterates through the added and removed repositories, executing
// activities to handle the synchronization process for each repository.
//
// The workflow uses a selector to manage concurrent execution of
// activities for added and removed repositories. It waits for all
// activities to complete before returning.
func SyncRepos(ctx workflow.Context, payload *defs.WebhookInstallRepos) error {
	selector := workflow.NewSelector(ctx)
	acts := &activities.InstallRepos{}
	total := make([]string, len(payload.RepositoriesAdded)+len(payload.RepositoriesRemoved))
	install := &entities.GithubInstallation{}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.
		ExecuteActivity(ctx, acts.GetInstallationForSync, payload.Installation.ID).
		Get(ctx, install); err != nil {
		return err
	}

	for _, repo := range payload.RepositoriesAdded {
		payload := &defs.SyncRepoPayload{InstallationID: install.ID, Repo: repo, OrgID: install.OrgID}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoAdded, payload), func(f workflow.Future) {})
	}

	for _, repo := range payload.RepositoriesRemoved {
		payload := &defs.SyncRepoPayload{InstallationID: install.ID, Repo: repo}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoRemoved, payload), func(f workflow.Future) {})
	}

	for range total {
		selector.Select(ctx)
	}

	return nil
}
