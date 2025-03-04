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

package defs

import (
	"go.breu.io/durex/workflows"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
)

// RepoWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repos.{id}.name.{name}"
func RepoWorkflowOptions(repo *entities.Repo) workflows.Options {
	opts := durable.NewWorkflowOptions(
		durable.WithOrg(repo.OrgID.String()),
		durable.WithSubject("repos"),
		durable.WithSubjectID(repo.ID.String()),
		durable.WithMeta("name", repo.Name),
	)

	return opts
}

// BranchWorkflowOptions returns workflow options for BranchCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repos.{id}.name.{name}.branch.{branch}"
func BranchWorkflowOptions(repo *entities.Repo, branch string) workflows.Options {
	opts := durable.NewWorkflowOptions(
		durable.WithOrg(repo.OrgID.String()),
		durable.WithSubject("repos"),
		durable.WithSubjectID(repo.ID.String()),
		durable.WithMeta("name", repo.Name),
		durable.WithMeta("branch", branch),
	)

	return opts
}

// TrunkWorkflowOptions returns workflow options for TrunkCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repos.{id}.name.{name}.branch.trunk"
func TrunkWorkflowOptions(repo *entities.Repo) workflows.Options {
	opts := durable.NewWorkflowOptions(
		durable.WithOrg(repo.OrgID.String()),
		durable.WithSubject("repos"),
		durable.WithSubjectID(repo.ID.String()),
		durable.WithMeta("name", repo.Name),
		durable.WithMeta("branch", "trunk"),
	)

	return opts
}
