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

type (
	WebhookInstall struct {
		Action       string              `json:"action"`
		Installation Installation        `json:"installation"`
		Repositories []PartialRepository `json:"repositories"`
		Sender       User                `json:"sender"`
	}

	WebhookInstallRepos struct {
		Action              string              `json:"action"`
		Installation        Installation        `json:"installation"`
		RepositorySelection string              `json:"repository_selection"`
		RepositoriesAdded   []PartialRepository `json:"repositories_added"`
		RepositoriesRemoved []PartialRepository `json:"repositories_removed"`
	}

	WebhookRef struct {
		Ref          string         `json:"ref"`
		RefType      string         `json:"ref_type"`
		MasterBranch *string        `json:"master_branch"` // NOTE: This is only present in the create event.
		Description  *string        `json:"description"`   // NOTE: This is only present in the create event.
		PusherType   string         `json:"pusher_type"`
		Repository   Repository     `json:"repository"`
		Organization Organization   `json:"organization"`
		Sender       User           `json:"sender"`
		Installation InstallationID `json:"installation"`
		IsCreated    bool           `json:"is_created"`
	}
)

func (wr *WebhookRef) GetRef() string {
	return wr.Ref
}

func (wr *WebhookRef) GetRefType() string {
	return wr.RefType
}
