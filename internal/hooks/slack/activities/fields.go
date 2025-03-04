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

package activities

import (
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/slack/attach"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func fields_lines_exceeded(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		attach.Repo(event),
		attach.Branch(event),
		attach.Threshold(),
		attach.TotalLinesCount(event),
		attach.LinesAdded(event),
		attach.LinesDeleted(event),
		attach.AddedFiles(event),
		attach.DeletedFiles(event),
		attach.ModifiedFiles(event),
		attach.RenameFiles(event),
	}

	return fields
}

func fields_merge_conflict(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		attach.Repo(event),
		attach.BranchMerge(event),
		attach.CurrentHead(event),
		attach.ConflictHead(),
		attach.AffectedFiles(),
	}

	return fields
}
