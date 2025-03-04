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

package intercepts

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/erratic"
)

// RequestLogger returns a unary interceptor that logs request and response information.
//
// It logs the peer address, protocol, HTTP method, and latency.
// If the request returns an error, it logs the error details, including the error ID, code, message, hints, and
// internal error (if any).
// It converts any error to a connect.Error using erratic.QuantmError.ToConnectError.
// The procedure name is used as the log message. Errors are logged at ERROR level, and successes are logged at
// INFO level.
func RequestLogger() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			peer := req.Peer()
			procedure := req.Spec().Procedure
			fields := []any{
				"peer_address", peer.Addr,
				"protocol", peer.Protocol,
				"method", req.HTTPMethod(),
			}

			resp, err := next(ctx, req)

			elapsed := time.Since(start)
			fields = append(fields, "latency", elapsed)

			if err != nil {
				qerr, ok := err.(*erratic.QuantmError)
				if !ok {
					qerr = erratic.NewUnknownError(erratic.CommonModule).Wrap(err)
				}

				fields = append(fields, "error_id", qerr.ID)
				fields = append(fields, "error_code", qerr.Code)
				fields = append(fields, "error", qerr.Error())

				for k, v := range qerr.Hints {
					fields = append(fields, k, v)
				}

				if qerr.Unwrap() != nil {
					fields = append(fields, "internal_error", qerr.Unwrap().Error())
				}

				slog.Warn(procedure, fields...)

				err = qerr.ToConnectError()
			} else {
				slog.Info(procedure, fields...)
			}

			return resp, err
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
