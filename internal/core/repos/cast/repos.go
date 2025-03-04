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

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// HookToProto converts an int32 representation of a RepoHook to a RepoHook proto.
func HookToProto(hook int32) eventsv1.RepoHook {
	v, ok := eventsv1.RepoHook_name[hook]
	if !ok {
		return eventsv1.RepoHook_REPO_HOOK_UNSPECIFIED
	}

	return eventsv1.RepoHook(eventsv1.RepoHook_value[v])
}

// RepoToProto converts a Repo entity to a Repo proto.
func RepoToProto(repo *entities.Repo) *corev1.Repo {
	return &corev1.Repo{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		OrgId:         repo.OrgID.String(),
		Name:          repo.Name,
		Hook:          HookToProto(repo.Hook),
		HookId:        repo.HookID.String(),
		DefaultBranch: repo.DefaultBranch,
		IsMonorepo:    repo.IsMonorepo,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToProto(repo.StaleDuration),
		Url:           repo.Url,
		IsActive:      repo.IsActive,
	}
}

// RepoToProto converts a Repo entity to a Repo proto.
func RepoExtendedRowToProto(repo *entities.ListReposRow) *corev1.RepoExtended {
	return &corev1.RepoExtended{
		Id:            repo.ID.String(),
		CreatedAt:     timestamppb.New(repo.CreatedAt),
		UpdatedAt:     timestamppb.New(repo.UpdatedAt),
		OrgId:         repo.OrgID.String(),
		Name:          repo.Name,
		Hook:          HookToProto(repo.Hook),
		HookId:        repo.HookID.String(),
		DefaultBranch: repo.DefaultBranch,
		IsMonorepo:    repo.IsMonorepo,
		Threshold:     repo.Threshold,
		StaleDuration: db.IntervalToProto(repo.StaleDuration),
		Url:           repo.Url,
		IsActive:      repo.IsActive,
		ChatEnabled:   repo.HasChat,
		ChannelName:   repo.ChannelName,
	}
}

// ReposToProto converts a slice of Repo entities to a slice of Repo protos.
func RepoExtendedListToProto(repos []entities.ListReposRow) []*corev1.RepoExtended {
	protos := make([]*corev1.RepoExtended, 0)
	for _, repo := range repos {
		protos = append(protos, RepoExtendedRowToProto(&repo))
	}

	return protos
}

// PushEventToRebaseEvent converts a Push event to a Rebase event.
func PushEventToRebaseEvent(
	push *events.Event[eventsv1.RepoHook, eventsv1.Push], parent uuid.UUID, base string,
) *events.Event[eventsv1.RepoHook, eventsv1.Rebase] {
	payload := &eventsv1.Rebase{
		Base:       base,
		Head:       push.Payload.After,
		Repository: push.Payload.Repository,
	}

	return events.
		Next[eventsv1.RepoHook, eventsv1.Push, eventsv1.Rebase](push, events.ScopeRebase, events.ActionRequested).
		SetPayload(payload)
}
