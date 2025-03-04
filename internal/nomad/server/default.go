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

package server

import (
	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/nomad/intercepts"
)

// DefaultServer creates a new Nomad server instance with the provided options.
//
// FIXME: create an insecure handler for user registration and login. AccountServiceHandler?
func DefaultServer(opts ...Option) *Server {
	srv := New(opts...)

	// -- config/interceptors --

	interceptors := []connect.Interceptor{
		intercepts.RequestLogger(),
	}

	// -- config/handlers --
	options := []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}

	// - insecure handlers -
	// -- auth --
	srv.add(auth.NomadAccountServiceHandler(options...))
	srv.add(auth.NomadOrgServiceHandler(options...))
	srv.add(auth.NomadUserServiceHandler(options...))

	// - secure handlers -

	options = append(options, connect.WithInterceptors(auth.NomadInterceptor()))

	// -- core/repos --
	srv.add(repos.NomadHandler(options...))

	// -- hooks/github --
	srv.add(github.NomadHandler(options...))

	// -- hooks/slack --
	srv.add(slack.NomadHandler(options...))

	return srv
}
