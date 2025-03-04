// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2025.
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
	"log/slog"
	"os"

	"go.breu.io/graceful"

	"go.breu.io/quantm/cmd/quantm/workers"
	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/nomad"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

type (
	ModeFn func(*graceful.Graceful) error
)

const (
	ServiceGithub     = "github"
	ServiceSlack      = "slack"
	ServiceKernel     = "kernel"
	ServiceDB         = "db"
	ServicePulse      = "pulse"
	ServiceDurable    = "durable"
	ServiceWebhook    = "webhook"
	ServiceNomad      = "nomad"
	ServiceCoreQueue  = "core_queue"
	ServiceHooksQueue = "hooks_queue"
)

// Setup configures the application based on the provided config.
func (c *Config) Setup(app *graceful.Graceful) error {
	modes := map[Mode]ModeFn{
		ModeMigrate: c.migrate,
		ModeWebhook: c.webhook,
		ModeGRPC:    c.grpc,
		ModeWorkers: c.workers,
		ModeDefault: c.all,
	}

	if fn, ok := modes[c.Mode]; ok {
		return fn(app)
	}

	return nil
}

func (c *Config) SetupLogger() {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	if c.Debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = false
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// SetupServices configures common services.
func (c *Config) SetupServices(app *graceful.Graceful) error {
	c.SetupLogger()
	auth.SetSecret(c.Secret)

	if err := c.Github.Validate(); err != nil {
		return err
	}

	github.Configure(github.WithConfig(c.Github))

	if err := c.Slack.Validate(); err != nil {
		return err
	}

	slack.Configure(slack.WithConfig(c.Slack))

	kernel.Configure(
		kernel.WithRepoHook(eventsv1.RepoHook_REPO_HOOK_GITHUB, &github.KernelImpl{}),
		kernel.WithChatHook(eventsv1.ChatHook_CHAT_HOOK_SLACK, &slack.KernelImpl{}),
	)

	if err := c.SetupDB(); err != nil {
		return err
	}

	if err := c.SetupDurable(); err != nil {
		return err
	}

	if err := c.SetupPulse(); err != nil {
		return err
	}

	app.Add(ServiceGithub, github.Get())
	// app.Add(ServicesSlack, slack.Get())
	app.Add(ServiceKernel, kernel.Get(), ServiceGithub)
	app.Add(ServiceDB, db.Get())
	app.Add(ServicePulse, pulse.Get())
	app.Add(ServiceDurable, durable.Get())

	return nil
}

// SetupDB configures the database.
func (c *Config) SetupDB() error {
	if err := c.DB.Validate(); err != nil {
		return err
	}

	db.Get(db.WithConfig(c.DB))

	return nil
}

// SetupDurable configures the durable service.
func (c *Config) SetupDurable() error {
	if err := c.Durable.Validate(); err != nil {
		return err
	}

	durable.Get(durable.WithConfig(c.Durable))

	return nil
}

// SetupPulse configures the pulse service.
func (c *Config) SetupPulse() error {
	if err := c.Pulse.Validate(); err != nil {
		return err
	}

	pulse.Get(pulse.WithConfig(c.Pulse))

	return nil
}

// migrate configures the application for database migrations.
func (c *Config) migrate(app *graceful.Graceful) error {
	c.SetupLogger()

	if err := c.SetupDB(); err != nil {
		return err
	}

	return nil
}

// webhook configures the application for webhook mode.
func (c *Config) webhook(app *graceful.Graceful) error {
	if err := c.SetupServices(app); err != nil {
		return err
	}

	app.Add(ServiceWebhook, NewWebhookServer(), ServiceDurable)

	return nil
}

// grpc configures the application for gRPC mode.
func (c *Config) grpc(app *graceful.Graceful) error {
	if err := c.SetupServices(app); err != nil {
		return err
	}

	app.Add(ServiceNomad, nomad.New(nomad.WithConfig(c.Nomad)), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)

	return nil
}

// workers configures the application for worker mode.
func (c *Config) workers(app *graceful.Graceful) error {
	if err := c.SetupServices(app); err != nil {
		return err
	}

	workers.Core()
	workers.Hooks()

	app.Add(ServiceCoreQueue, durable.OnCore(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	app.Add(ServiceHooksQueue, durable.OnHooks(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)

	return nil
}

// all configures the application with all services and modes.
func (c *Config) all(app *graceful.Graceful) error {
	if err := c.SetupServices(app); err != nil {
		return err
	}

	workers.Core()
	workers.Hooks()

	app.Add(ServiceWebhook, NewWebhookServer(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	app.Add(ServiceNomad, nomad.New(nomad.WithConfig(c.Nomad)), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	app.Add(ServiceCoreQueue, durable.OnCore(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	app.Add(ServiceHooksQueue, durable.OnHooks(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)

	return nil
}
