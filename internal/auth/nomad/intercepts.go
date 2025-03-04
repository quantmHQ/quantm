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

package nomad

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/auth/config"
)

func AuthInterceptor() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			header := req.Header().Get("Authorization")

			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")

				cliams, err := config.DecodeJWE(config.Secret(), token)
				if err != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, err)
				}

				if cliams != nil {
					ctx = context.WithValue(ctx, AuthContextUser, cliams.UserID)
					ctx = context.WithValue(ctx, AuthContextOrg, cliams.OrgID)
				}
			} else {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing bearer token"))
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
