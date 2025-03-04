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

package pulse

import (
	"log/slog"
	"sync"

	"go.breu.io/quantm/internal/pulse/config"
)

type (
	Config = config.Config // Config represents the configuration for the Pulse package.
	Option = config.Option // Option is a functional option to configure Pulse.
)

var (
	DefaultConfig = config.DefaultConfig // DefaultConfig holds the default configuration values.

	_c   *Config   // _c stores the configured instance of Pulse.
	once sync.Once // once ensures the initialization happens only once.
)

// WithConfig allows customizing the Pulse configuration using functional options.
//
// It takes a Config pointer as input and returns an Option. This Option can then be passed to the Instance function.
func WithConfig(cfg *Config) Option {
	return config.WithConfig(cfg)
}

// Get returns the singleton instance of the Pulse configuration.
//
// It initializes Pulse with the provided options if it hasn't been initialized yet. The initialization is thread-safe,
// guaranteed by the sync.Once usage.  Get returns a pointer to the initialized Config instance.
func Get(opts ...Option) *Config {
	once.Do(func() {
		slog.Info("pulse: configuring ...")

		_c = config.New(opts...)
	})

	return _c
}
