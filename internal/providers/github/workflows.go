// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
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

package github

import (
	"log/slog"
	"time"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared/queues"
)

var (
	activities *Activities
)

type (
	// Workflows is the entry point for all workflows for GitHub.
	Workflows struct{}

	// InstallationWorkflowStatus handles the status of the workflow Workflows.OnInstallationEvent.
	InstallationWorkflowStatus struct {
		WebhookDone bool
		RequestDone bool
	}

	// PullRequestWorkflowStatus handles the status of the workflow Workflows.OnPullRequestEvent.
	PullRequestWorkflowStatus struct {
		Complete bool
	}
)

// OnInstallationEvent workflow is executed when we initiate the installation of GitHub defs.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong, and we can receive the complete installation
// request before the push event. To handle this, we use temporal.io's signal API to provide two possible entry points
// for the system. See the README.md for a detailed explanation on how this workflow works.
//
// NOTE: This workflow is only meant to be started with SignalWithStartWorkflow.
// TODO: Refactor this workflow to reduce complexity.
func (w *Workflows) OnInstallationEvent(ctx workflow.Context) (*Installation, error) { // nolint:funlen
	// prelude
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)
	installation := &Installation{}
	webhook := &InstallationEvent{}
	request := &CompleteInstallationSignal{}
	status := &InstallationWorkflowStatus{WebhookDone: false, RequestDone: false}
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, activityOpts)

	// setting up channels to receive signals
	webhookChannel := workflow.GetSignalChannel(ctx, WorkflowSignalInstallationEvent.String())
	requestChannel := workflow.GetSignalChannel(ctx, WorkflowSignalCompleteInstallation.String())

	// setting up callbacks for the channels
	selector.AddReceive(webhookChannel, on_install_webhook_signal(ctx, webhook, status))
	selector.AddReceive(requestChannel, on_install_request_signal(ctx, request, status))

	logger.Info("github/installation: waiting for webhook and complete installation request signals ...")

	// keep listening for signals until we have received both the installation id and the team id
	for !(status.WebhookDone && status.RequestDone) {
		selector.Select(ctx)
	}

	logger.Info("github/installation: required signals processed ...")

	switch webhook.Action {
	// NOTE - Since a GitHub organization can only have one active installation at a time, when a new installation is created, it's
	// considered the first app installation for the organization, and we assume no teams have been created yet within the organization.
	//
	// TODO - we need to handle the case when an the app uninstallation and reinstallation case.
	//
	// - when delete event is received, we need to add a db field to mark the installation as deleted.
	// - on the subsequent installation, we need to check if the installation is deleted and update the installation status.
	case "created":
		user := &auth.User{}
		team := &auth.Team{}

		if err := workflow.ExecuteActivity(_ctx, activities.GetUserByID, request.UserID.String()).Get(ctx, user); err != nil {
			return nil, err
		}

		if user.TeamID.String() == db.NullUUID {
			logger.Info("github/installation: no team associated, creating a new team ...")

			team.Name = webhook.Installation.Account.Login

			_ = workflow.ExecuteActivity(_ctx, activities.CreateTeam, team).Get(ctx, team)

			logger.Info("github/installation: team created, assigning to user ...")

			user.TeamID = team.ID
			_ = workflow.ExecuteActivity(_ctx, activities.SaveUser, user).Get(ctx, user)
		} else {
			logger.Warn("github/installation: team already associated, fetching ...")

			_ = workflow.ExecuteActivity(_ctx, activities.GetTeamByID, user.TeamID.String()).Get(ctx, team)
		}

		// Finalizing the installation
		installation.TeamID = team.ID
		installation.InstallationID = webhook.Installation.ID
		installation.InstallationLogin = webhook.Installation.Account.Login // Github Organization name
		installation.InstallationLoginID = webhook.Installation.Account.ID  // Github organization ID
		installation.InstallationType = webhook.Installation.Account.Type
		installation.SenderID = webhook.Sender.ID
		installation.SenderLogin = webhook.Sender.Login
		installation.Status = webhook.Action

		logger.Info("github/installation: creating or updating installation ...")

		if err := workflow.ExecuteActivity(_ctx, activities.CreateOrUpdateInstallation, installation).Get(_ctx, installation); err != nil {
			logger.Error("github/installation: error saving installation ...", "error", err)
		}

		logger.Info("github/installation: updating user associations ...")

		membership := &CreateMembershipsPayload{
			UserID:        user.ID,
			TeamID:        team.ID,
			IsAdmin:       true,
			GithubOrgName: webhook.Installation.Account.Login,
			GithubOrgID:   webhook.Installation.Account.ID,
			GithubUserID:  webhook.Sender.ID,
		}

		if err := workflow.ExecuteActivity(_ctx, activities.CreateMemberships, membership).Get(_ctx, nil); err != nil {
			logger.Error("github/installation: error saving installation ...", "error", err)
		}

		logger.Info("github/installation: saving installation repos ...")

		for _, repo := range webhook.Repositories {
			logger.Info("github/installation: saving repository ...")
			logger.Debug("repository", "repository", repo)

			repo := &Repo{
				GithubID:        repo.ID,
				InstallationID:  installation.InstallationID,
				Name:            repo.Name,
				FullName:        repo.FullName,
				DefaultBranch:   "main",
				HasEarlyWarning: false,
				IsActive:        true,
				TeamID:          installation.TeamID,
			}

			future := workflow.ExecuteActivity(_ctx, activities.CreateOrUpdateGithubRepo, repo)

			// NOTE - ideally, we should use a new selector here, but since there will be no new signals comings in, we know that
			// selector.Select will only be waiting for the futures to complete.
			selector.AddFuture(future, on_repo_saved_future(ctx, repo))
		}

		logger.Info("github/installation: waiting for repositories to be saved ...")

		for range webhook.Repositories {
			selector.Select(ctx)
		}

		logger.Info("github/installation: installation repositories saved ...")
	case "deleted", "suspend", "unsuspend":
		logger.Warn("github/installation: installation removed, unhandled case ...")
	default:
		logger.Warn("github/installation: unhandled action during installation ...", slog.String("action", webhook.Action))
	}

	logger.Info("github/installation: complete", slog.Any("installation", installation))

	return installation, nil
}

// PostInstall refresh the default branch for all repositories associated with the given teamID and gets orgs users.
// NOTE - this workflow runs complete for the first time but when reinstall the github app and configure the same repos.
// it will give the, It will give the access_token error: could not refresh installation id XXXXXXX's token error. TODO
// - handle when the github app is reinstall and confgure the same repos, and also need to test when configure the same
// repo or new repos.
func (w *Workflows) PostInstall(ctx workflow.Context, payload *Installation) error {
	logger := workflow.GetLogger(ctx)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	logger.Info(
		"github/installation/post: starting ...",
		slog.String("installation_id", payload.InstallationID.String()),
		slog.String("installation_login", payload.InstallationLogin),
	)

	// TODO: move this inside a workflow.Go statement
	logger.Info("github/installation/post: syncing repos ...", "installation_id", payload.InstallationID.String())

	sync := &SyncReposFromGithubPayload{
		InstallationID: payload.InstallationID,
		Owner:          payload.InstallationLogin,
		TeamID:         payload.TeamID,
	}
	if err := workflow.ExecuteActivity(_ctx, activities.SyncReposFromGithub, sync).Get(_ctx, nil); err != nil {
		logger.Error("github/installation/post: error syncing repos ...", "error", err)
	}

	logger.Info("github/installation/post: syncing github org users ...", "installation_id", payload.InstallationID.String())

	// TODO: sync users
	orgsync := &SyncOrgUsersFromGithubPayload{
		InstallationID: payload.InstallationID,
		GithubOrgName:  payload.InstallationLogin,
		GithubOrgID:    payload.InstallationLoginID,
	}
	if err := workflow.ExecuteActivity(_ctx, activities.SyncOrgUsersFromGithub, orgsync).Get(_ctx, nil); err != nil {
		logger.Error("github/installation/post: error syncing org users ...", "error", err)
	}

	return nil
}

// OnPushEvent is run when ever a repo event is received. Repo Event can be push event or a create event.
func (w *Workflows) OnCreateOrDeleteEvent(ctx workflow.Context, payload *CreateOrDeleteEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/create-delete: fetching metadata ...")

	prepared := PrepareRepoEventPayload(payload)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/create-delete: preparing event ...")

	event := payload.normalize(state.CoreRepo)
	if state.User != nil {
		event.SetUserID(state.User.ID)
	} else {
		logger.Warn("github/create-delete: unable to set user id ...")
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	logger.Info("github/create-delete: dispatching event ...")

	return workflow.
		ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalCreateOrDelete, event).
		Get(ctx, nil)
}

// OnPushEvent is run when ever a repo event is received. Repo Event can be push event or a create event.
func (w *Workflows) OnPushEvent(ctx workflow.Context, payload *PushEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/push: fetching metadata ...")

	prepared := PrepareRepoEventPayload(payload)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/push: preparing event ...")

	event := payload.normalize(state.CoreRepo)
	if state.User != nil {
		event.SetUserID(state.User.ID)
	} else {
		logger.Warn("github/push: unable to set user id ...")
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	logger.Info("github/push: dispatching event ...")

	return workflow.ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPush, event).Get(ctx, nil)
}

// OnPullRequestEvent normalize the pull request event and then signal the core repo.
func (w *Workflows) OnPullRequestEvent(ctx workflow.Context, payload *PullRequestEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/pull_request: fetching metadata ...")

	prepared := PrepareRepoEventPayload(payload)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/pull_request: preparing event ...")

	event := payload.normalize(state.CoreRepo)
	if state.User != nil {
		event.SetUserID(state.User.ID)
	} else {
		logger.Warn("github/pull_request: unable to set user id ...")
	}

	label := payload.as_label(event) // this will be nil if scope is label

	ctx = dispatch.WithDefaultActivityContext(ctx)

	logger.Info("github/pull_request: dispatching event ...")

	if label == nil {
		return workflow.
			ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequest, event).Get(ctx, nil)
	}

	return workflow.
		ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestLabel, label).
		Get(ctx, nil)
}

// OnPullRequestReviewEvent normalize the pull request review event and then signal the core repo.
func (w *Workflows) OnPullRequestReviewEvent(ctx workflow.Context, event *PullRequestReviewEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/pull_request_review: fetching metadata ...")

	prepared := PrepareRepoEventPayload(event)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/pull_request_review: preparing event ...")

	payload := event.normalize(state.CoreRepo)
	if state.User != nil {
		payload.SetUserID(state.User.ID)
	} else {
		logger.Warn("github/pull_request_review: unable to set user id ...")
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	logger.Info("github/pull_request_review: dispatching event ...")

	return workflow.
		ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestReview, payload).Get(ctx, nil)
}

// OnPullRequestReviewCommentEvent normalize the pull request review comment event and then signal the core repo.
func (w *Workflows) OnPullRequestReviewCommentEvent(ctx workflow.Context, event *PullRequestReviewCommentEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/pull_request_review_comment: fetching metadata ...")

	prepared := PrepareRepoEventPayload(event)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/pull_request_review_comment: preparing event ...")

	payload := event.normalize(state.CoreRepo)
	if state.User != nil {
		payload.SetUserID(state.User.ID)
	} else {
		logger.Warn("github/pull_request_review_comment: unable to set user id ...")
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	logger.Info("github/pull_request_review_comment: dispatching event ...")

	return workflow.
		ExecuteActivity(ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestComment, payload).
		Get(ctx, nil)
}

// OnInstallationRepositoriesEvent is responsible when a repository is added or removed from an installation.
func (w *Workflows) OnInstallationRepositoriesEvent(ctx workflow.Context, payload *InstallationRepositoriesEvent) error {
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)

	logger.Info("received installation repositories event ...")

	installation := &Installation{}
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	err := workflow.
		ExecuteActivity(actx, activities.GetInstallation, payload.Installation.ID).
		Get(actx, installation)
	if err != nil {
		logger.Error("error getting installation", "error", err)
		return err
	}

	for _, repository := range payload.RepositoriesAdded {
		logger.Info("saving repository ...")
		logger.Debug("repository", "repository", repository)

		repo := &Repo{
			GithubID:       repository.ID,
			InstallationID: installation.InstallationID,
			Name:           repository.Name,
			FullName:       repository.FullName,
			TeamID:         installation.TeamID,
		}

		future := workflow.ExecuteActivity(actx, activities.CreateOrUpdateGithubRepo, repo)
		selector.AddFuture(future, on_repo_saved_future(ctx, repo))
	}

	// wait for all the repositories to be saved.
	for range payload.RepositoriesAdded {
		selector.Select(ctx)
	}

	return nil
}

func (w *Workflows) OnWorkflowRunEvent(ctx workflow.Context, pl *GithubWorkflowRunEvent) error {
	logger := workflow.GetLogger(ctx)
	state := &RepoEventState{}

	logger.Info("github/workflow_run: fetching metadata ...")

	prepared := PrepareRepoEventPayload(pl)

	future, err := queues.
		Providers().
		ExecuteChildWorkflow(ctx, PrepareRepoEventChildWorkflowOptions(ctx), PrepareRepoEventWorkflow, prepared)
	if err != nil {
		return err
	}

	if err := future.Get(ctx, state); err != nil {
		return err
	}

	logger.Info("github/workflow_run: preparing event ...")

	payload := &defs.RepoIOSignalWorkflowRunPayload{
		RepoName:       pl.Repository.Name,
		RepoOwner:      pl.Repository.Owner.Login,
		InstallationID: pl.Installation.ID,
	}

	p := &defs.RepoIOWorkflowActionPayload{
		RepoName:       payload.RepoName,
		RepoOwner:      payload.RepoOwner,
		InstallationID: payload.InstallationID,
	}

	winfo := &defs.RepoIOWorkflowInfo{}

	return workflow.ExecuteActivity(ctx, activities.GithubWorkflowInfo, p).Get(ctx, winfo)
}

// on_repo_saved_future handles post-processing after a repository is saved against an installation.
func on_repo_saved_future(ctx workflow.Context, payload *Repo) defs.FutureHandler {
	logger := workflow.GetLogger(ctx)
	return func(f workflow.Future) { logger.Info("repository saved ...", "repo", payload.GithubID) }
}

// on_install_webhook_signal handles webhook events for installation that is in progress.
func on_install_webhook_signal(
	ctx workflow.Context, installation *InstallationEvent, status *InstallationWorkflowStatus,
) defs.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("github/installation: webhook received ...", "action", installation.Action)
		channel.Receive(ctx, installation)

		status.WebhookDone = true

		switch installation.Action {
		case "deleted", "suspend", "unsuspend":
			logger.Info("github/installation: installation removed ....", "action", installation.Action)

			status.RequestDone = true
		case "request":
			logger.Info("github/installation: installation request ...", "action", installation.Action)

			status.RequestDone = true
		default:
			logger.Info("github/installation: create action ...", "action", installation.Action)
		}
	}
}

// on_install_request_signal handles new http requests on an installation in progress.
func on_install_request_signal(
	ctx workflow.Context, installation *CompleteInstallationSignal, status *InstallationWorkflowStatus,
) defs.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("github/installation: received complete installation request ...")
		channel.Receive(ctx, installation)

		status.RequestDone = true
	}
}

func PrepareRepoEventWorkflow(ctx workflow.Context, event *RepoEventPayload) (*RepoEventState, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("github/repo_event: preparing ...")

	acts := &Activities{}
	state := &RepoEventState{}
	ctx = dispatch.WithDefaultActivityContext(ctx)
	repos := make([]Repo, 0)

	/**
	 * Get github repo for the installation.
	 */

	if err := workflow.ExecuteActivity(ctx, acts.GetReposForInstallation, event.InstallationID.String(), event.RepoID.String()).
		Get(ctx, &repos); err != nil {
		return state, err
	}

	if len(repos) == 0 {
		return state, NewRepoNotFoundRepoEventError(event.InstallationID, event.RepoID, event.RepoName)
	}

	if len(repos) > 1 {
		return state, NewMultipleReposFoundRepoEventError(event.InstallationID, event.RepoID, event.RepoName)
	}

	repo := &repos[0]

	if !repo.IsActive {
		return state, NewInactiveRepoRepoEventError(event.InstallationID, event.RepoID, event.RepoName)
	}

	if !repo.HasEarlyWarning {
		return state, NewHasNoEarlyWarningRepoEventError(event.InstallationID, event.RepoID, event.RepoName)
	}

	state.Repo = repo

	/**
	 * Get the linked core repository.
	 */

	if err := workflow.ExecuteActivity(ctx, acts.GetCoreRepo, repo.ID.String()).Get(ctx, &state.CoreRepo); err != nil {
		return state, err
	}

	/**
	 * Get the user for the installation.
	 */

	if err := workflow.ExecuteActivity(ctx, acts.GetTeamUserByLoginID, event.InstallationID.String()).Get(ctx, &state.User); err != nil {
		return state, err
	}

	return state, nil
}
