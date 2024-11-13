package githubwfs

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	PushWorkflowState struct {
		log log.Logger
	}
)

func Push(ctx workflow.Context, payload *githubdefs.Push) error {
	acts := &githubacts.Push{}

	var eventory *githubdefs.Eventory[commonv1.RepoHook, eventsv1.Push]
	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &eventory); err != nil {
		return err
	}

	// TODO - need to confirm the signature
	return workflow.ExecuteActivity(
		ctx, githubacts.SignalCoreRepo, eventory.Repo, githubdefs.SignalWebhookPush, eventory.Event).
		Get(ctx, nil)
}
