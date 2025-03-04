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

package events

type (
	Scope string // EventScope is the scope of the event.
)

// String returns the string representation of the EventScope.
func (es Scope) String() string { return string(es) }

const (
	ScopeBranch     Scope = "branch"      // ScopeBranch scopes branch event.
	ScopeTag        Scope = "tag"         // ScopeTag scopes tag event.
	ScopePush       Scope = "push"        // ScopePush scopes push event.
	ScopeRebase     Scope = "rebase"      // ScopeRebase scopes rebase event.
	ScopeDiff       Scope = "diff"        // ScopeDiff scopes diff event.
	ScopePr         Scope = "pr"          // ScopePr scopes pull request event.
	ScopePrLabel    Scope = "pr_label"    // ScopePrLabel scopes pull request label event.
	ScopeMerge      Scope = "merge"       // ScopeMerge scopes merge event.
	ScopeMergeQueue Scope = "merge_queue" // ScopeMergeQueue scopes merge queue event.
)
