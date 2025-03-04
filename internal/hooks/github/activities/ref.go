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

package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Ref groups all the activities required for the Github Ref.
	Ref struct{}
)

// HydrateGithubRefEvent hydrates the branch event with the given parameters.
func (b *Ref) HydrateGithubRefEvent(ctx context.Context, params *defs.HydratedRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

// SignalRepoWithGithubRef signals the repository with the hydrated branch event.
func (b *Ref) SignalRepoWithGithubRef(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.GitRef]) error {
	return SignalRepo(ctx, hydrated)
}
