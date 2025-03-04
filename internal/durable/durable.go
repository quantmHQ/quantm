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

package durable

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/durable/config"
)

// -- Internal --

var (
	// configured is the instantiated configuration.
	configured *Config
	configonce sync.Once

	// coreq is the core queue.
	coreq     queues.Queue
	coreqonce sync.Once

	// hooksq is the hooks queue.
	hooksq     queues.Queue
	hooksqonce sync.Once
)

// -- Types --

type (
	// Config represents the configuration for the durable layer.
	Config = config.Config

	// ConfigOption is an option for configuring the durable layer.
	ConfigOption = config.Option
)

// -- Configuration --

var (
	// DefaultConfig is the default configuration for the durable layer.
	DefaultConfig = config.Default

	WithConfig = config.WithConfig
)

// Get is a singleton that holds the temporal client.
//
// Please note that the actual client is not created until the first call Get().Client().
func Get(opts ...ConfigOption) *Config {
	configonce.Do(func() {
		configured = config.New(opts...)
	})

	return configured
}

// -- Queues --

// OnCore returns the core queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.core.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnCore() queues.Queue {
	coreqonce.Do(func() {
		client, err := Get().Client()
		if err != nil {
			slog.Error("durable: unable to connect to durable server ...", "error", err.Error())
			panic(err)
		}

		coreq = queues.New(queues.WithName("core"), queues.WithClient(client))
	})

	return coreq
}

// OnHooks returns the hooks queue.
//
// All workflows on this queue will have the ID prefix of
//
//	io.ctrlpane.hooks.{block}.{block_id}.{element}.{element_id}.{modifier}.{modifier_id}....
func OnHooks() queues.Queue {
	hooksqonce.Do(func() {
		client, err := Get().Client()
		if err != nil {
			slog.Error("durable: unable to connect to durable server ...", "error", err.Error())
			panic(err)
		}

		hooksq = queues.New(queues.WithName("hooks"), queues.WithClient(client))
	})

	return hooksq
}
