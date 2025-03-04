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
	"errors"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	Install struct{}
)

// GetOrCreateInstallation retrieves a Github installation from the database by installation ID.
// If the installation does not exist, it creates a new one.
func (a *Install) GetOrCreateInstallation(
	ctx context.Context, install *entities.GithubInstallation,
) (*entities.GithubInstallation, error) {
	response, err := db.Queries().GetGithubInstallationByInstallationID(ctx, install.InstallationID)
	if err == nil {
		return &response, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		create := entities.CreateGithubInstallationParams{
			OrgID:             install.OrgID,
			InstallationID:    install.InstallationID,
			InstallationLogin: install.InstallationLogin,
			InstallationType:  install.InstallationType,
			SenderID:          install.SenderID,
			SenderLogin:       install.SenderLogin,
		}

		response, err = db.Queries().CreateGithubInstallation(ctx, create)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, err
}

func (a *Install) AddRepoToInstall(ctx context.Context, payload *defs.SyncRepoPayload) error {
	return AddRepo(ctx, payload)
}
