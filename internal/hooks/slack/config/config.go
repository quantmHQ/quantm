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
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/slack-go/slack"
)

var (
	_once sync.Once
	_c    *Config
)

// Config holds the configuration for the Slack client.
type (
	Config struct {
		ClientID     string `koanf:"CLIENT_ID" validate:"required"`
		ClientSecret string `koanf:"CLIENT_SECRET" validate:"required"`
		RedirectURL  string `koanf:"REDIRECT_URL" validate:"required"`
		Debug        bool   `koanf:"DEBUG"`
	}

	ConfigOption func(*Config)
)

func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// GetSlackClient creates a new Slack client using the token.
func GetSlackClient(token string) (*slack.Client, error) {
	lgr := &logger{slog.Default().WithGroup("slack")}
	client := slack.New(
		token,
		slack.OptionDebug(_c.Debug),
		slack.OptionLog(lgr),
	)

	return client, nil
}

func ClientID() string {
	return _c.ClientID
}

func ClientSecret() string {
	return _c.ClientSecret
}

func ClientRedirectURL() string {
	return _c.RedirectURL
}

func WithConfig(cfg *Config) ConfigOption {
	return func(config *Config) {
		config.ClientID = cfg.ClientID
		config.ClientSecret = cfg.ClientSecret
		config.RedirectURL = cfg.RedirectURL
	}
}

func Instance(opts ...ConfigOption) *Config {
	_once.Do(func() {
		_c = &Config{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	return _c
}
