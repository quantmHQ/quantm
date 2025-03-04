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

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ResolveRevision resolves a revision string to its corresponding commit hash.
//
// Unsupported reference types (trees, annotated tags) are rejected.
//
// Supported revisions:
//   - HEAD, branches, tags,
//   - remote-tracking branches, HEAD~n, HEAD^,
//   - refspec selectors (e.g., HEAD^{/fix bug}),
//   - hash prefixes/full hashes.
func (r *Repository) ResolveRevision(ctx context.Context, revision string) (*plumbing.Hash, error) {
	if r.cloned == nil {
		if err := r.Open(); err != nil {
			return nil, err
		}
	}

	return r.cloned.ResolveRevision(plumbing.Revision(revision))
}

// ResolveCommit resolves a revision string to its corresponding commit object.
//
// The same rules as ResolveRevision apply here.
func (r *Repository) ResolveCommit(ctx context.Context, revision string) (*object.Commit, error) {
	hash, err := r.ResolveRevision(ctx, revision)
	if err != nil {
		return nil, err
	}

	commit, err := r.cloned.CommitObject(*hash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
