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

package migrations

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	dbcfg "go.breu.io/quantm/internal/db/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	//go:embed postgres/*.sql
	sql embed.FS
)

// Run runs the migrations for the PostgreSQL database.
func Run(ctx context.Context, connection *dbcfg.Config) error {
	slog.Info("migrations: running ...")

	if !connection.IsConnected() {
		_ = connection.Start(ctx)
		defer func() { _ = connection.Stop(ctx) }()
	}

	dir, err := iofs.New(sql, "postgres")
	if err != nil {
		slog.Error("migrations: unable to read ...", "error", err.Error())

		return err
	}

	migrations, err := migrate.NewWithSourceInstance(
		"iofs",
		dir,
		connection.ConnectionURI(),
	)

	if err != nil {
		slog.Error("migrations: unable to read data ...", "error", err.Error())
		return err
	}

	version, dirty, _ := migrations.Version()

	if dirty {
		return fmt.Errorf("migrations:  unapplied migrations at version, %d", version)
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		slog.Info("migrations: nothing new since ...", "version", version)
	}

	slog.Info("migrations: migrations done successfully")

	return nil
}
