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

package activities

import (
	"context"
	"fmt"
	"net/http"

	ghi "github.com/bradleyfalzon/ghinstallation/v2"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/config"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Kernel is the implmentation of kernel.Repo interface.
	//
	// Please note that this must never be called from the workflows.
	Kernel struct{}
)

func (k *Kernel) TokenizedCloneUrl(ctx context.Context, repo *entities.Repo) (string, error) {
	ghrepo, err := db.Queries().GetGithubRepoByID(ctx, repo.HookID)
	if err != nil {
		return "", err
	}

	install, err := db.Queries().GetGithubInstallation(ctx, ghrepo.InstallationID)
	if err != nil {
		return "", err
	}

	client, err := ghi.New(http.DefaultTransport, config.Instance().AppID, install.InstallationID, []byte(config.Instance().PrivateKey))
	if err != nil {
		return "", err
	}

	token, err := client.Token(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://git:%s@github.com/%s.git", token, ghrepo.FullName), nil
}

func (k *Kernel) DetectChanges(ctx context.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Push]) error {
	return nil
}
