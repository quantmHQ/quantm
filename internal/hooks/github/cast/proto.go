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
	"slices"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func RefToProto(ref *defs.WebhookRef) eventsv1.GitRef {
	return eventsv1.GitRef{
		Ref:  ref.GetRef(),
		Kind: ref.GetRefType(),
	}
}

func PushToProto(push *defs.Push) eventsv1.Push {
	return eventsv1.Push{
		Ref:        push.GetRef(),
		Before:     push.GetBefore(),
		After:      push.GetAfter(),
		Repository: push.GetRepositoryName(),
		SenderId:   push.GetSenderID(),
		Timestamp:  timestamppb.New(time.Now()),
		Commits:    CommitsToProto(push.GetCommits()),
	}
}

func CommitsToProto(commits []defs.Commit) []*eventsv1.Commit {
	result := make([]*eventsv1.Commit, len(commits))
	for i, commit := range commits {
		result[i] = &eventsv1.Commit{
			Sha:       commit.GetID(),
			Message:   commit.GetMessage(),
			Url:       commit.GetURL(),
			Timestamp: timestamppb.New(commit.GetTimestamp()),
			Added:     commit.GetAdded(),
			Removed:   commit.GetRemoved(),
			Modified:  commit.GetModified(),
		}
	}

	return result
}

func PullRequestToProto(pr *defs.PR) eventsv1.PullRequest {
	return eventsv1.PullRequest{
		Number:     pr.GetNumber(),
		Title:      pr.GetTitle(),
		Body:       pr.GetBody(),
		Author:     pr.GetAuthor(),
		HeadBranch: pr.GetHeadBranch(),
		BaseBranch: pr.GetBaseBranch(),
		Timestamp:  timestamppb.New(pr.GetTimestamp()),
	}
}

func PullRequestLabelToProto(pr *defs.PR) *eventsv1.MergeQueue {
	valid := []string{repos.LabelMerge, repos.LabelPriority}

	if slices.Contains(valid, pr.GetLabelName()) {
		proto := &eventsv1.MergeQueue{
			Number:    pr.GetNumber(),
			Branch:    pr.GetHeadBranch(),
			Timestamp: timestamppb.New(pr.GetTimestamp()),
		}

		if pr.GetLabelName() == repos.LabelPriority {
			proto.IsPriority = true
		}

		return proto
	}

	return nil
}

func PrReviewToProto(prr *defs.PrReview) eventsv1.PullRequestReview {
	return eventsv1.PullRequestReview{
		Id:                prr.GetPrReviewID(),
		PullRequestNumber: prr.GetPrNumber(),
		Branch:            prr.GetHeadBranch(),
		State:             prr.GetState(),
		AuthorEmail:       *prr.GetSenderEmail(),
		SubmittedAt:       timestamppb.New(prr.GetSubmittedAt()),
	}
}

func PrReviewCommentToProto(prrc *defs.PrReviewComment) eventsv1.PullRequestReviewComment {
	return eventsv1.PullRequestReviewComment{
		Id:                prrc.GetCommentID(),
		PullRequestNumber: prrc.GetPrNumber(),
		Branch:            prrc.GetHeadBranch(),
		State:             prrc.GetState(),
		ReviewId:          prrc.GetReviewID(),
		CommitSha:         prrc.GetCommitSha(),
		Path:              prrc.GetPath(),
		Position:          prrc.GetPosition(),
		InReplyTo:         *prrc.GetInReplyTo(),
		AuthorEmail:       *prrc.GetSenderEmail(),
		SubmittedAt:       timestamppb.New(prrc.GetSubmittedAt()),
	}
}
