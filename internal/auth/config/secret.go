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

package config

import (
	"log/slog"
	"sync/atomic"
)

const (
	_default string = "set me"
)

var (
	secret atomic.Value
)

func init() { secret.Store(_default) }

// Secret returns the configured secret value. It will log a warning if the secret is not set.
func Secret() string {
	if !IsValid() {
		slog.Warn("auth: secret is not set, configure it using the environment variable 'SECRET'")
	}

	return secret.Load().(string)
}

// SetSecret sets the secret value.
func SetSecret(val string) {
	secret.Store(val)
}

// IsValid returns true if the secret is valid, false otherwise.
func IsValid() bool {
	return secret.Load().(string) != _default
}
