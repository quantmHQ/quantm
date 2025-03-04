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

package nomad

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/auth/cast"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	AccountService struct {
		authv1connect.UnimplementedAccountServiceHandler
	}
)

func (s *AccountService) GetAccountByProviderAccountID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountByProviderAccountIDRequest],
) (*connect.Response[authv1.GetAccountByProviderAccountIDResponse], error) {
	params := cast.ProtoToGetAccountByProviderAccountIDParams(req.Msg)

	account, err := db.Queries().GetOAuthAccountByProviderAccountID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, erratic.NewNotFoundError(erratic.AuthModule).
				AddHint("account_id", params.ProviderAccountID).
				AddHint("provider", params.Provider)
		}

		return nil, erratic.NewDatabaseError(erratic.AuthModule).
			AddHint("account_id", params.ProviderAccountID).
			AddHint("provider", params.Provider).Wrap(err)
	}

	proto := &authv1.GetAccountByProviderAccountIDResponse{Account: cast.AccountToProto(&account)}

	return connect.NewResponse(proto), nil
}

func (s *AccountService) GetAccountsByUserID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountsByUserIDRequest],
) (*connect.Response[authv1.GetAccountsByUserIDResponse], error) {
	id, err := uuid.Parse(req.Msg.GetUserId())
	if err != nil {
		return nil, erratic.NewBadRequestError(erratic.AuthModule).
			AddHint("user_id", req.Msg.GetUserId()).Wrap(err)
	}

	accounts, err := db.Queries().GetOAuthAccountsByUserID(ctx, id)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.AuthModule).
			AddHint("user_id", req.Msg.GetUserId()).Wrap(err)
	}

	proto := make([]*authv1.Account, len(accounts))
	for i, account := range accounts {
		proto[i] = cast.AccountToProto(&account)
	}

	return connect.NewResponse(&authv1.GetAccountsByUserIDResponse{Accounts: proto}), nil
}

func (s *AccountService) CreateAccount(
	ctx context.Context,
	req *connect.Request[authv1.CreateAccountRequest],
) (*connect.Response[authv1.CreateAccountResponse], error) {
	params := cast.ProtoToCreateAccountParams(req.Msg)

	account, err := db.Queries().CreateOAuthAccount(ctx, params)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.AuthModule).Wrap(err)
	}

	return connect.NewResponse(&authv1.CreateAccountResponse{Account: cast.AccountToProto(&account)}), nil
}

func (s *AccountService) GetAccountByID(
	ctx context.Context,
	req *connect.Request[authv1.GetAccountByIDRequest],
) (*connect.Response[authv1.GetAccountByIDResponse], error) {
	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, erratic.NewBadRequestError(erratic.AuthModule).
			AddHint("id", req.Msg.GetId()).Wrap(err)
	}

	account, err := db.Queries().GetOAuthAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, erratic.NewNotFoundError(erratic.AuthModule, "account").AddHint("id", req.Msg.GetId())
		}

		return nil, erratic.NewDatabaseError(erratic.AuthModule).AddHint("id", req.Msg.GetId()).Wrap(err)
	}

	return connect.NewResponse(&authv1.GetAccountByIDResponse{Account: cast.AccountToProto(&account)}), nil
}

func NewAccountSericeServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewAccountServiceHandler(
		&AccountService{},
		opts...,
	)
}
