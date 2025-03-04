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

package git

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/cast"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Repository struct {
		Entity *entities.Repo
		Branch string
		Path   string

		cloned *gogit.Repository
	}
)

func (r *Repository) Clone(ctx context.Context) error {
	if r.cloned != nil {
		return NewRepositoryError(r, OpClone)
	}

	hook := cast.HookToProto(r.Entity.Hook)
	ref := plumbing.NewBranchReferenceName(r.Branch)

	if err := ref.Validate(); err != nil {
		return NewRepositoryError(r, OpClone).Wrap(err)
	}

	url, err := kernel.Get().RepoHook(hook).TokenizedCloneUrl(ctx, r.Entity)
	if err != nil {
		return NewRepositoryError(r, OpClone).Wrap(err)
	}

	cloned, err := gogit.PlainCloneContext(ctx, r.Path, false, &gogit.CloneOptions{
		URL:           url,
		ReferenceName: ref,
		SingleBranch:  false,
	})
	if err != nil {
		return NewRepositoryError(r, OpClone).Wrap(err)
	}

	r.cloned = cloned

	return nil
}

func (r *Repository) Open() error {
	if r.cloned != nil {
		return nil
	}

	cloned, err := gogit.PlainOpen(r.Path)
	if err != nil {
		return NewRepositoryError(r, OpOpen).Wrap(err)
	}

	r.cloned = cloned

	return nil
}

func NewRepository(entity *entities.Repo, branch, path string) *Repository {
	return &Repository{
		Entity: entity,
		Branch: branch,
		Path:   path,
	}
}
