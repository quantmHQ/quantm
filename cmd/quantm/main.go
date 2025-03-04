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

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/graceful"

	"go.breu.io/quantm/cmd/quantm/config"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
)

func main() {
	conf := config.New()
	ctx := context.Background()

	conf.Load()
	conf.Parse()

	// - run migrations and exit if mode is migrate
	if conf.Mode == config.ModeMigrate {
		if err := migrations.Run(ctx, db.Get(db.WithConfig(conf.DB))); err != nil {
			slog.Error("unable to run migrations", "error", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	// - run the app based on mode, exit 1 on error else wait for signal

	quit := make(chan os.Signal, 1)
	app := graceful.New()

	if err := conf.Setup(app); err != nil {
		slog.Error("unable to setup ...", "error", err.Error())
		os.Exit(1)
	}

	if err := app.Start(ctx); err != nil {
		slog.Error("unable to start ...", "error", err.Error())

		os.Exit(1)
	}

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	<-quit

	// - gracefully stop the services

	if err := app.Stop(ctx); err != nil {
		slog.Error("unable to stop service", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("service stopped, exiting...")

	os.Exit(0)
}
