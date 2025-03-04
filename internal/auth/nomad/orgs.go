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
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	OrgService struct {
		authv1connect.UnimplementedOrgServiceHandler
	}
)

func (s *OrgService) SetOrgHooks(
	ctx context.Context, req *connect.Request[authv1.SetOrgHooksRequest],
) (*connect.Response[emptypb.Empty], error) {
	hooks, err := json.Marshal(req.Msg.Hooks)
	if err != nil {
		return nil, erratic.NewBadRequestError(erratic.AuthModule).WithReason("unable to detect hook").Wrap(err)
	}

	params := entities.SetOrgHooksParams{ID: uuid.MustParse(req.Msg.GetOrgId()), Hooks: hooks}

	err = db.Queries().SetOrgHooks(ctx, params)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.AuthModule).WithReason("unable to set org hooks").Wrap(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewOrgServiceServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewOrgServiceHandler(
		&OrgService{},
		opts...,
	)
}
