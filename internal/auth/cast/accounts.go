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

package cast

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/db/entities"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
)

// AccountToProto converts an OauthAccount entity to its protobuf representation.
func AccountToProto(account *entities.OauthAccount) *authv1.Account {
	return &authv1.Account{
		Id:                account.ID.String(),
		CreatedAt:         timestamppb.New(account.CreatedAt),
		UpdatedAt:         timestamppb.New(account.UpdatedAt),
		ExpiresAt:         timestamppb.New(account.ExpiresAt),
		UserId:            account.UserID.String(),
		Provider:          AuthProviderToProto(account.Provider),
		ProviderAccountId: account.ProviderAccountID,
		Kind:              account.Type,
	}
}

// ProtoToAccount converts a protobuf account to its OauthAccount entity representation.
func ProtoToAccount(proto *authv1.Account) *entities.OauthAccount {
	return &entities.OauthAccount{
		ID:                uuid.MustParse(proto.GetId()),
		CreatedAt:         proto.GetCreatedAt().AsTime(),
		UpdatedAt:         proto.GetUpdatedAt().AsTime(),
		UserID:            uuid.MustParse(proto.GetUserId()),
		Provider:          ProtoToAuthProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
		ExpiresAt:         proto.GetExpiresAt().AsTime(),
		Type:              proto.GetKind(),
	}
}

// ProtoToGetAccountsByUserIDParams converts a protobuf GetAccountsByUserIDRequest to a UUID.
func ProtoToGetAccountsByUserIDParams(proto *authv1.GetAccountsByUserIDRequest) uuid.UUID {
	return uuid.MustParse(proto.GetUserId())
}

// ProtoToCreateAccountParams converts a protobuf CreateAccountRequest to an entities.CreateOAuthAccountParams.
func ProtoToCreateAccountParams(proto *authv1.CreateAccountRequest) entities.CreateOAuthAccountParams {
	return entities.CreateOAuthAccountParams{
		UserID:            uuid.MustParse(proto.GetUserId()),
		Provider:          ProtoToAuthProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
		ExpiresAt:         proto.GetExpiresAt().AsTime(),
		Type:              proto.GetKind(),
	}
}

// ProtoToGetAccountByIDParams converts a protobuf GetAccountByIDRequest to a UUID.
func ProtoToGetAccountByIDParams(proto *authv1.GetAccountByIDRequest) uuid.UUID {
	return uuid.MustParse(proto.GetId())
}

// ProtoToGetAccountByProviderAccountIDParams converts a protobuf GetAccountByProviderAccountIDRequest to an
// entities.GetOAuthAccountByProviderAccountIDParams.
func ProtoToGetAccountByProviderAccountIDParams(
	proto *authv1.GetAccountByProviderAccountIDRequest,
) entities.GetOAuthAccountByProviderAccountIDParams {
	return entities.GetOAuthAccountByProviderAccountIDParams{
		Provider:          ProtoToAuthProvider(proto.GetProvider()),
		ProviderAccountID: proto.GetProviderAccountId(),
	}
}
