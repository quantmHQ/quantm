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

package db

import (
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	dbcfg "go.breu.io/quantm/internal/db/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	//go:embed migrations/postgres/*.sql
	sql embed.FS
)

// WithPostgresMigrations configures PostgreSQL database migrations.
// TODO - move to function return.
func WithPostgresMigrations() {
	// TODO: read from .env
	c := &dbcfg.Default

	dir, err := iofs.New(sql, "migrations/postgres")
	if err != nil {
		slog.Error("db: failed to initialize migrations", "error", err.Error())
		return
	}

	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		dir,
		c.ConnectionURI(),
	)
	if err != nil {
		slog.Error("db: failed to create migrations instance", "error", err.Error())
		return
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Warn("db: failed to run migrations", "error", err.Error())
		return
	}

	if err == migrate.ErrNoChange {
		slog.Info("db: no new migrations to run")
	}

	slog.Info("db: migrations done successfully")
}
