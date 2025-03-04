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
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (r *Repository) CherryPick(ctx context.Context, branch, hash string) (*object.Commit, error) {
	if r.cloned == nil {
		if err := r.Open(); err != nil {
			return nil, NewRepositoryError(r, OpOpen).Wrap(err)
		}
	}

	pick, err := r.ResolveCommit(ctx, hash)
	if err != nil {
		return nil, NewResolveError(r, OpResolveCommit, hash).Wrap(err)
	}

	worktree, err := r.cloned.Worktree()
	if err != nil {
		return nil, NewCherryPickError(r, "worktree", hash).Wrap(err)
	}

	err = worktree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Create: false,
	})

	if err != nil {
		return nil, NewCherryPickError(r, "checkout", hash).Wrap(err)
	}

	commit, err := worktree.Commit(pick.Message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  pick.Author.Name,
			Email: pick.Author.Email,
			When:  pick.Author.When,
		},
	})

	if err != nil {
		return nil, NewCherryPickError(r, "commit", hash).Wrap(err)
	}

	cp, err := r.cloned.CommitObject(commit)
	if err != nil {
		return nil, NewCherryPickError(r, "commit_object", hash).Wrap(err)
	}

	if err := worktree.Checkout(&gogit.CheckoutOptions{
		Hash:   cp.Hash,
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	}); err != nil {
		return nil, NewCherryPickError(r, "checkout_post_cherrypick", hash).Wrap(err)
	}

	return cp, nil
}
