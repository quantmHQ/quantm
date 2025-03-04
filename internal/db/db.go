// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2025.
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

package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Config = config.Config
)

var (
	DefaultConfig = config.Default
)

func WithConfig(conf *Config) config.ConfigOption {
	return config.WithConfig(conf)
}

// Get is a wrapper around the dbcfg.Instance singleton.
func Get(opts ...config.ConfigOption) *config.Config {
	return config.Instance(opts...)
}

// Queries is a wrapper around the dbcfg.Queries singleton.
func Queries() *entities.Queries {
	return config.Queries()
}

// Transaction begins the transaction and wraps the queries in a transaction.
//
// Example:
//
//	tx, qtx, err := db.Transaction(ctx)
//	if err != nil { return err }
//
//	defer func() { _ = tx.Rollback(ctx) }()
//
//	// Do something with qtx. Any time you return on error, the transaction will be rolled back.
//	...
//
//	// Commit the transaction.
//	err = tx.Commit(ctx)
//	if err != nil { return err }
//
//	return nil
func Transaction(ctx context.Context) (pgx.Tx, *entities.Queries, error) {
	tx, err := Get().Pool().Begin(ctx)
	if err != nil {
		slog.Error("db: error creating transaction ...", "error", err.Error())

		return nil, nil, err
	}

	return tx, Queries().WithTx(tx), nil
}
