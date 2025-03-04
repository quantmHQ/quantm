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

// ProtoToTeam converts a protobuf team to its Team entity representation.
func ProtoToTeam(proto *authv1.Team) *entities.Team {
	return &entities.Team{
		ID:        uuid.MustParse(proto.GetId()),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
		Name:      proto.GetName(),
		Slug:      proto.GetSlug(),
	}
}

// TeamToProto converts a Team entity to its protobuf representation.
func TeamToProto(team *entities.Team) *authv1.Team {
	return &authv1.Team{
		Id:        team.ID.String(),
		CreatedAt: timestamppb.New(team.CreatedAt),
		UpdatedAt: timestamppb.New(team.UpdatedAt),
		Name:      team.Name,
		Slug:      team.Slug,
	}
}

// ProtoToCreateTeamParams converts a protobuf CreateTeamRequest to a CreateTeamParams.
func ProtoToCreateTeamParams(proto *authv1.CreateTeamRequest) entities.CreateTeamParams {
	return entities.CreateTeamParams{
		OrgID: uuid.MustParse(proto.GetOrgId()),
		Name:  proto.GetName(),
	}
}

// GetTeamBySlugRowToProto converts a GetTeamBySlugRow entity to its protobuf representation.
func GetTeamBySlugRowToProto(team entities.GetTeamBySlugRow) *authv1.Team {
	return &authv1.Team{
		Id:   team.ID.String(),
		Name: team.Name,
	}
}
