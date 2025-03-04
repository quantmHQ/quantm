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
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/repos/cast"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/core/v1/corev1connect"
)

type (
	RepoService struct {
		corev1connect.UnimplementedRepoServiceHandler
	}
)

func (s *RepoService) ListRepos(
	ctx context.Context, req *connect.Request[emptypb.Empty],
) (*connect.Response[corev1.ListReposResponse], error) {
	_, org_id := auth.NomadAuthContext(ctx)

	rows, err := db.Queries().ListRepos(ctx, org_id)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.CoreModule).Wrap(err)
	}

	protos := cast.RepoExtendedListToProto(rows)

	return connect.NewResponse(&corev1.ListReposResponse{Repos: protos}), nil
}

func NewRepoServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return corev1connect.NewRepoServiceHandler(&RepoService{}, opts...)
}
