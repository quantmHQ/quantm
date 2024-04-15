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

package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/shared"
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

// OnInstallationEvent workflow is executed when we initiate the installation of GitHub core.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong, and we can receive the complete installation
// request before the push event. To handle this, we use temporal.io's signal API to provide two possible entry points
// for the system. See the README.md for a detailed explanation on how this workflow works.
//
// NOTE: This workflow is only meant to be started with SignalWithStartWorkflow.
func (w *Workflows) OnInstallationEvent(ctx workflow.Context) error {
	// prelude
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)
	webhook := &InstallationEvent{}
	request := &CompleteInstallationSignal{}
	status := &InstallationWorkflowStatus{WebhookDone: false, RequestDone: false}

	// setting up channels to receive signals
	webhookChannel := workflow.GetSignalChannel(ctx, WorkflowSignalInstallationEvent.String())
	requestChannel := workflow.GetSignalChannel(ctx, WorkflowSignalCompleteInstallation.String())

	// setting up callbacks for the channels
	selector.AddReceive(webhookChannel, onInstallationWebhookSignal(ctx, webhook, status))
	selector.AddReceive(requestChannel, onRequestSignal(ctx, request, status))

	// keep listening for signals until we have received both the installation id and the team id
	for !(status.WebhookDone && status.RequestDone) {
		logger.Info("waiting for signals ....")
		selector.Select(ctx)
	}

	logger.Info("all signals received, processing ...")

	// Finalizing the installation
	installation := &Installation{
		TeamID:            request.TeamID,
		InstallationID:    webhook.Installation.ID,
		InstallationLogin: webhook.Installation.Account.Login,
		InstallationType:  webhook.Installation.Account.Type,
		SenderID:          webhook.Sender.ID,
		SenderLogin:       webhook.Sender.Login,
		Status:            webhook.Action,
	}

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.
		ExecuteActivity(actx, activities.CreateOrUpdateInstallation, installation).
		Get(actx, installation)

	if err != nil {
		logger.Error("error saving installation", "error", err)
		return err
	}

	// If webhook.Action == "created", save the repository information to the database.
	if webhook.Action == "created" {
		logger.Info("saving associated repositories ...")

		// asynchronously save the repos
		for _, repository := range webhook.Repositories {
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
			selector.AddFuture(future, onCreateOrUpdateRepoActivityFuture(ctx, repo))
		}

		// wait for all repositories to be saved.
		for range webhook.Repositories {
			selector.Select(ctx)
		}
	}

	logger.Info("installation complete")
	logger.Debug("installation", "installation", installation)

	return nil
}

// OnPushEvent checks if the push event is associated with an open pull request.If so, it will get the idempotent key for
// the immutable rollout. Depending upon the target branch, it will either queue the rollout or update the existing
// rollout.
func (w *Workflows) OnPushEvent(ctx workflow.Context, payload *PushEvent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("received push event ...")

	branchName := payload.Ref
	if parts := strings.Split(payload.Ref, "/"); len(parts) == 3 {
		branchName = parts[len(parts)-1]
	}

	if branchName == payload.Repository.DefaultBranch || strings.Contains(branchName, "-tempcopy-for-target-") {
		shared.Logger().Debug("OnPushEvent", "push on default branch or temp branch", branchName)
		return nil
	}

	signalPayload := &shared.PushEventSignal{
		RefBranch:      branchName,
		RepoID:         payload.Repository.ID,
		RepoName:       payload.Repository.Name,
		RepoOwner:      payload.Repository.Owner.Login,
		DefaultBranch:  payload.Repository.DefaultBranch,
		InstallationID: payload.Installation.ID,
		RepoProvider:   "github",
	}

	cw := &core.Workflows{}
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("repo"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowElement("branch"),
			shared.WithWorkflowElementID(branchName),
			shared.WithWorkflowProp("type", "early_detection"),
		)

	_, err := shared.Temporal().
		Client().SignalWithStartWorkflow(
		context.Background(),
		opts.ID,
		shared.WorkflowPushEvent.String(),
		signalPayload,
		opts,
		cw.BranchController,
	)
	if err != nil {
		shared.Logger().Error("OnPushEvent", "Error signaling workflow", err)
		return err
	}

	return nil
}

func (w *Workflows) OnGithubActionResult(ctx workflow.Context, payload *WorkflowRun) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("OnGithubActionResult", "entry", "workflow started")

	// wait for github action to return success status
	ch := workflow.GetSignalChannel(ctx, WorkflowSignalActionResult.String())
	gh_result := &GithubActionResult{}
	ch.Receive(ctx, gh_result)

	logger.Info("OnGithubActionResult", "action recvd", gh_result)

	// acquiring lock here
	lock, err := LockInstance(ctx, fmt.Sprint(payload.Repository.ID))
	if err != nil {
		logger.Error("Error in getting lock instance", "Error", err)
		return err
	}

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	var mergeCommit string
	err = workflow.ExecuteActivity(actx, activities.RebaseAndMerge, payload.Repository.Owner.Login, payload.Repository.Name,
		payload.WR.HeadBranch, payload.Installation.ID).Get(ctx, &mergeCommit)

	if err != nil {
		logger.Error("error getting installation", "error", err)
		return err
	}

	err = lock.Release(ctx)
	if err != nil {
		logger.Error("error releasing lock", "error", err)
		return err
	}

	// Signal stack workflow about changeset update
	// info to be sent: repo, commit
	// get workflowID for the stack attached to this repo

	// get core repo
	repo := &Repo{GithubID: payload.Repository.ID}
	coreRepo := &core.Repo{}

	err = workflow.ExecuteActivity(actx, activities.GetCoreRepo, repo).Get(ctx, coreRepo)
	if err != nil {
		logger.Error("error getting core repo", "error", err)
		return err
	}

	// get core workflow ID for this stack
	coreWorkflowID := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowID(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	// signal core stack workflow
	logger.Info("core workflow id", "ID", coreWorkflowID)

	signalPayload := &shared.CreateChangesetSignal{
		RepoTableID: coreRepo.ID,
		RepoID:      fmt.Sprint(payload.Repository.ID),
		CommitID:    mergeCommit,
	}

	options := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	cw := &core.Workflows{}
	_, err = shared.Temporal().Client().SignalWithStartWorkflow(
		context.Background(),
		coreWorkflowID,
		shared.WorkflowSignalCreateChangeset.String(),
		signalPayload,
		options,
		cw.StackController,
		coreRepo.StackID.String(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (w *Workflows) OnLabelEvent(ctx workflow.Context, payload *PullRequestEvent) error {
	shared.Logger().Info("OnLabelEvent", "entry", "workflow started")

	logger := workflow.GetLogger(ctx)

	logger.Info("received PR label event ...")

	installationID := payload.Installation.ID
	repoOwner := payload.Repository.Owner.Login
	repoName := payload.Repository.Name
	pullRequestID := payload.Number
	label := payload.Label.Name
	branch := payload.PullRequest.Head.Ref

	if label == fmt.Sprintf("quantm ready") {
		logger.Debug("quantm ready label applied")

		cw := &core.Workflows{}
		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			WorkflowOptions(
				shared.WithWorkflowBlock("repo"),
				shared.WithWorkflowBlockID(strconv.FormatInt(payload.Repository.ID, 10)),
				shared.WithWorkflowElement("PR"),
				shared.WithWorkflowElementID(fmt.Sprint(pullRequestID)),
				shared.WithWorkflowProp("type", "merge_queue"),
			)

		payload2 := &shared.MergeQueueSignal{
			PullRequestID:  pullRequestID,
			InstallationID: installationID,
			RepoOwner:      repoOwner,
			RepoName:       repoName,
			Branch:         branch,
		}

		_, err := shared.Temporal().
			Client().SignalWithStartWorkflow(
			context.Background(),
			opts.ID,
			shared.MergeQueueStarted.String(),
			payload2,
			opts,
			cw.PollMergeQueue,
		)
		if err != nil {
			shared.Logger().Error("OnLabelEvent", "Error signaling workflow", err)
			return err
		}

		shared.Logger().Info("PR sent to MergeQueue")
	}

	return nil
}

// OnPullRequestEvent workflow is responsible to get or create the idempotency key for the changeset controller workflow.
// Regardless of the action on PR, the algorithm needs to arrive at the same idempotency key! One possible way is
// to calculate the checksum  of different components. The trick would be to handle "synchronize" event as this relates
// to a new commit on the PR.
//
//   - One possible way to handle "synchronize" would be to only listen to label events on the PR.
//   - The other possible way to create an idempotency key would be to take the state, create a version set and then tag
//     the git commit with the version set. We can also take a look at aviator.co to see how they are creating version-sets.
//
// After the creation of the idempotency key, we pass the idempotency key as a signal to the Aperture Workflow.
func (w *Workflows) OnPullRequestEvent(ctx workflow.Context, payload *PullRequestEvent) error {
	shared.Logger().Info("OnPullRequestEvent", "entry", "workflow started")

	logger := workflow.GetLogger(ctx)
	// status := &PullRequestWorkflowStatus{Complete: false}

	// wait for artifact to generate and push to registery
	ch := workflow.GetSignalChannel(ctx, WorkflowSignalArtifactReady.String())
	artifact := &ArtifactReadySignal{}
	ch.Receive(ctx, artifact)

	// setting activity options
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	// get core repo
	repo := &Repo{GithubID: payload.Repository.ID}
	coreRepo := &core.Repo{}

	err := workflow.ExecuteActivity(actx, activities.GetCoreRepo, repo).Get(ctx, coreRepo)
	if err != nil {
		logger.Error("error getting core repo", "error", err)
		return err
	}

	// get core workflow ID for this stack
	corePRWfID := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowID(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	// payload for core stack workflow
	signalPayload := &shared.PullRequestSignal{
		RepoID:           coreRepo.ID,
		SenderWorkflowID: workflow.GetInfo(ctx).WorkflowExecution.ID,
		TriggerID:        payload.PullRequest.ID,
		Image:            artifact.Image,
		Digest:           artifact.Digest,
		ImageRegistry:    artifact.Registry,
	}

	// signal core stack workflow
	logger.Info("core workflow id", "ID", corePRWfID)

	options := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	cw := &core.Workflows{}
	_, _ = shared.Temporal().Client().SignalWithStartWorkflow(
		context.Background(),
		corePRWfID,
		shared.WorkflowSignalDeploymentStarted.String(),
		signalPayload,
		options,
		cw.StackController,
		coreRepo.StackID.String(),
	)
	// workflow.SignalExternalWorkflow(ctx, corePRWfID, "", shared.WorkflowSignalPullRequest.String(), signalPayload).Get(ctx, nil)
	logger.Debug("Signaled workflow", "ID", signalPayload.SenderWorkflowID, " core repo ID: ", signalPayload.RepoID.String())

	// workflow.GetSignalChannel(ctx, WorkflowSignalPullRequestProcessed.String()).Receive(ctx, &status)

	// signal processor
	// selector.AddReceive(prChannel, onPRSignal(ctx, pr, status))

	// logger.Info("PR created: scheduling new aperture at the application level.")

	// // keep listening to signals until complete = true
	// for !status.Complete {
	// 	selector.Select(ctx)
	// }

	return nil
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
		Get(ctx, installation)
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
		selector.AddFuture(future, onCreateOrUpdateRepoActivityFuture(ctx, repo))
	}

	// wait for all the repositories to be saved.
	for range payload.RepositoriesAdded {
		selector.Select(ctx)
	}

	return nil
}

// onCreateOrUpdateRepoActivityFuture handles post-processing after a repository is saved against an installation.
func onCreateOrUpdateRepoActivityFuture(ctx workflow.Context, payload *Repo) shared.FutureHandler {
	logger := workflow.GetLogger(ctx)
	return func(f workflow.Future) { logger.Info("repository saved ...", "repo", payload.GithubID) }
}

// onInstallationWebhookSignal handles webhook events for installation that is in progress.
func onInstallationWebhookSignal(
	ctx workflow.Context, installation *InstallationEvent, status *InstallationWorkflowStatus,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("received webhook installation event ...", "action", installation.Action)
		channel.Receive(ctx, installation)

		status.WebhookDone = true

		switch installation.Action {
		case "deleted", "suspend", "unsuspend":
			logger.Info("installation removed, skipping complete installation request ...")

			status.RequestDone = true
		default:
			logger.Info("installation created, waiting for complete installation request ...")
		}
	}
}

// onRequestSignal handles new http requests on an installation in progress.
func onRequestSignal(
	ctx workflow.Context, installation *CompleteInstallationSignal, status *InstallationWorkflowStatus,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("received complete installation request ...")
		channel.Receive(ctx, installation)

		status.RequestDone = true
	}
}

// // onPRSignal handles incoming signals on open PR.
// func onPRSignal(ctx workflow.Context, pr *PullRequestEvent, status *PullRequestWorkflowStatus) shared.ChannelHandler {
// 	logger := workflow.GetLogger(ctx)

// 	return func(channel workflow.ReceiveChannel, more bool) {
// 		channel.Receive(ctx, pr)

// 		switch pr.Action {
// 		case "closed":
// 			logger.Info("PR closed: scheduling aperture to be abandoned.", "action", pr.Action)

// 			if pr.PullRequest.Merged {
// 				logger.Info("PR merged: scheduling aperture to finish with conclusion.")

// 				// TODO: send the signal to the aperture workflow to finish with conclusion.
// 				status.Complete = true
// 			} else {
// 				logger.Info("PR closed: abort aperture.")

// 				status.Complete = true
// 			}
// 		case "synchronize":
// 			logger.Info("PR updated: checking the status of the environment ...", "action", pr.Action)
// 			// TODO: here we need to check the app associated with repo & get the `release` branch. If the PR branch is not
// 			// the default branch, then we update in place, otherwise, we queue a new rollout.
// 		default:
// 			logger.Info("PR: no action required, skipping ...", "action", pr.Action)
// 		}
// 	}
// }
