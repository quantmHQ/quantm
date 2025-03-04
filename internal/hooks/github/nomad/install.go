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
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/hooks/github/defs"
	"go.breu.io/quantm/internal/hooks/github/workflows"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/proto/hooks/github/v1/githubv1connect"
)

type (
	GithubService struct {
		githubv1connect.UnimplementedGithubServiceHandler
	}
)

func (s *GithubService) Install(
	ctx context.Context, req *connect.Request[githubv1.InstallRequest],
) (*connect.Response[emptypb.Empty], error) {
	if req.Msg.Action != githubv1.SetupAction_INSTALL {
		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	opts := defs.NewInstallWorkflowOptions(req.Msg.InstallationId, req.Msg.Action)
	args := defs.RequestInstall{
		InstallationID: req.Msg.InstallationId,
		SetupAction:    req.Msg.Action,
		OrgID:          uuid.MustParse(req.Msg.OrgId),
	}

	_, err := durable.OnHooks().SignalWithStartWorkflow(ctx, opts, defs.SignalRequestInstall, args, workflows.Install)
	if err != nil {
		return nil, erratic.NewSystemError(erratic.HooksGithubModule).WithReason("unable to signal hook")
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewGithubServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return githubv1connect.NewGithubServiceHandler(&GithubService{}, opts...)
}
