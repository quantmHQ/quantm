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

package config

import (
	"context"
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/db/entities"
)

var (
	_qry       *entities.Queries // Global database queries instance.
	_queryonce sync.Once         // Ensures queries initialization occurs only once.
)

// Queries returns a singleton instance of SQLC-generated queries, initialized with the Connection singleton's database connection.
//
// If no connection exists, Queries establishes one using the default environment-based configuration.  For more predictable behavior,
// manually initialize the connection singleton using Instance() followed by Start() in the main function.
func Queries() *entities.Queries {
	_queryonce.Do(func() {
		slog.Info("db: initializing queries ...")

		if _c == nil {
			slog.Warn("db: no connection, attempting to create connection using environment variables ...")

			conn := Instance(WithConfigFromEnvironment())
			_ = conn.Start(context.Background())
		}

		_qry = entities.New(Instance().pool)
	})

	return _qry
}
