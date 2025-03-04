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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (r *Repository) Ancestor(a, b plumbing.Hash) (*object.Commit, error) {
	onto, err := r.cloned.CommitObject(a)
	if err != nil {
		return nil, err
	}

	upstream, err := r.cloned.CommitObject(b)
	if err != nil {
		return nil, err
	}

	ancestors, err := onto.MergeBase(upstream)
	if err != nil {
		return nil, err
	}

	if len(ancestors) == 0 {
		return nil, NewCompareError(r, OpAncestor, a.String(), b.String())
	}

	return ancestors[0], nil
}
