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

package workers

import (
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/pulse"
)

// Hooks registers the activites and workflows for the hooks queue.
func Hooks() {
	q := durable.OnHooks()

	q.CreateWorker()

	if q != nil {
		// Register pulse activities
		q.RegisterActivity(pulse.PersistRepoEvent)
		q.RegisterActivity(pulse.PersistChatEvent)

		// Register github install workflow and activity
		q.RegisterWorkflow(github.InstallWorkflow)
		q.RegisterActivity(&github.InstallActivity{})

		// Register github sync repos workflow and activity
		q.RegisterWorkflow(github.SyncReposWorkflow)
		q.RegisterActivity(&github.InstallReposActivity{})

		// Register github push workflow and activity
		q.RegisterWorkflow(github.PushWorkflow)
		q.RegisterActivity(&github.PushActivity{})

		// Register github ref workflow and activity
		q.RegisterWorkflow(github.RefWorkflow)
		q.RegisterActivity(&github.RefActivity{})
	}
}
