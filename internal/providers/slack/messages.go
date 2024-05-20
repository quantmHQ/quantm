// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package slack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/core"
)

const (
	footer = "Powered by quantm"
)

func formatLineThresholdExceededAttachment(payload *core.LinesExceedSlackMessageProviderPayload) slack.Attachment {
	return slack.Attachment{
		Color: "warning",
		Pretext: "The number of lines in this pull request exceeds the allowed threshold. " +
			"Please review and adjust accordingly.", // TODO: need to finalize
		Title:     "PR Lines Exceed",
		TitleLink: payload.DetectChanges.CompareUrl,
		Fields: []slack.AttachmentField{
			createRepositoryField(payload.RepoName, payload.DetectChanges.RepoUrl),
			createBranchField(payload.BranchName, payload.DetectChanges.CompareUrl),
			{
				Title: "*Threshold*",
				Value: fmt.Sprintf("%d", payload.Threshold),
				Short: true,
			},
			{
				Title: "*Total Lines Count*",
				Value: fmt.Sprintf("%d", payload.DetectChanges.Delta),
				Short: true,
			},
			{
				Title: "*Lines Added*",
				Value: fmt.Sprintf("%d", payload.DetectChanges.Added),
				Short: true,
			},
			{
				Title: "*Lines Deleted*",
				Value: fmt.Sprintf("%d", payload.DetectChanges.Removed),
				Short: true,
			},
			{
				Title: "*Details*",
				Value: fmt.Sprintf("*Number of Files Changed:* %d\n*Files Changed:*\n%s",
					len(payload.DetectChanges.Modified), formatFilesList(payload.DetectChanges.Modified)),
				Short: false,
			},
		},
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func formatMergeConflictAttachment(merge *core.LatestCommit) slack.Attachment {
	return slack.Attachment{
		Color:     "warning",
		Pretext:   "Merge conflict detected. Please resolve the conflict.", // TODO: need to finalize
		Title:     "Merge Conflict",
		TitleLink: merge.CommitUrl,
		Fields: []slack.AttachmentField{
			createRepositoryField(merge.RepoName, merge.RepoUrl),
			createBranchField(merge.Branch, merge.CommitUrl),
		},
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func formatStaleBranchAttachment(staleBranch *core.LatestCommit) slack.Attachment {
	return slack.Attachment{
		Color:     "warning",
		Pretext:   "Stale branch is detected. Please review and take necessary action.", // TODO: need to finalize
		Title:     "Stale Branch",
		TitleLink: staleBranch.CommitUrl,
		Fields: []slack.AttachmentField{
			createRepositoryField(staleBranch.RepoName, staleBranch.RepoUrl),
			createBranchField(staleBranch.Branch, staleBranch.CommitUrl),
		},
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func createRepositoryField(repoName, repoURL string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Repository*",
		Value: fmt.Sprintf("<%s|%s>", repoURL, repoName),
		Short: true,
	}
}

func createBranchField(branchName, compareUrl string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s|%s>", compareUrl, branchName),
		Short: true,
	}
}

func formatFilesList(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}
