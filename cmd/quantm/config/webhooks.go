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
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/hooks/github"
)

type (
	// WebhookService is a webserver to manage webhooks for all the hooks. It conforms to the graceful.Service
	// interface, allowing for graceful start and shutdown. It wraps echo.Echo to provide this functionality.
	WebhookService struct {
		*echo.Echo
	}
)

func (w *WebhookService) Start(ctx context.Context) error {
	slog.Info("webhook: starting ...", "port", 8000)

	err := w.Echo.Start(":8000")
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (w *WebhookService) Stop(ctx context.Context) error {
	return w.Echo.Shutdown(ctx)
}

func NewWebhookServer() *WebhookService {
	webhook := echo.New()
	webhook.HideBanner = true
	webhook.HidePort = true

	github := &github.Webhook{}

	webhook.POST("/webhooks/github", github.Handler)

	return &WebhookService{webhook}
}
