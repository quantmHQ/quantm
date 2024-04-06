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

package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/shared"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 // TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Workflows        struct{}
	GetAssetsPayload struct {
		StackID       string
		RepoID        gocql.UUID
		ChangeSetID   gocql.UUID
		Image         string
		ImageRegistry string
		Digest        string
	}
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// StackController runs indefinitely and controls and synchronizes all actions on stack.
// This workflow will start when createStack call is received. it will be the master workflow for all child stack workflows
// for tasks like creating infrastructure, doing deployment, apperture controller etc.
//
// The workflow waits for the signals from the git provider. It consumes events for PR created, updated, merged etc.
func (w *Workflows) StackController(ctx workflow.Context, stackID string) error {
	// wait for merge complete signal
	ch := workflow.GetSignalChannel(ctx, shared.WorkflowSignalCreateChangeset.String())
	payload := &shared.CreateChangesetSignal{}
	ch.Receive(ctx, payload)

	shared.Logger().Info("StackController", "signal payload", payload)

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	// get repos for stack
	repos := SlicedResult[Repo]{}
	if err := workflow.ExecuteActivity(actx, activities.GetRepos, stackID).Get(ctx, &repos); err != nil {
		shared.Logger().Error("GetRepos providers failed", "error", err)
		return err
	}

	shared.Logger().Debug("StackController : going to create repomarkers")

	providerActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
	}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get commits against the repos
	repoMarkers := make([]ChangeSetRepoMarker, len(repos.Data))
	for idx, repo := range repos.Data {
		marker := &repoMarkers[idx]
		// p := Instance().Provider(repo.Provider) // get the specific provider
		p := Instance().RepoProvider(repo.Provider) // get the specific provider
		commitID := ""

		if err := workflow.
			ExecuteActivity(pctx, p.GetLatestCommit, repo.ProviderID, repo.DefaultBranch).
			Get(ctx, &commitID); err != nil {
			shared.Logger().Error("Error in getting latest commit ID", "repo", repo.Name, "provider", repo.Provider)
			return fmt.Errorf("Error in getting latest commit ID repo:%s, provider:%s", repo.Name, repo.Provider.String())
		}

		marker.CommitID = commitID
		marker.Provider = repo.Provider.String()
		marker.RepoID = repo.ID.String()

		// update commit id for the recently changed repo
		if marker.RepoID == payload.RepoID {
			marker.CommitID = payload.CommitID
			shared.Logger().Debug("Commit ID updated for repo " + marker.RepoID)
			marker.HasChanged = true // the repo in which commit was made
		}

		shared.Logger().Debug("Repo", "Name", repo.Name, "Repo marker", marker)
	}

	shared.Logger().Debug("StackController", "repomarkers created", repoMarkers)

	// create changeset before deploying the updated changeset
	changesetID, _ := gocql.RandomUUID()
	stackUUID, _ := gocql.ParseUUID(stackID)
	changeset := &ChangeSet{
		RepoMarkers: repoMarkers,
		ID:          changesetID,
		StackID:     stackUUID,
	}

	if err := workflow.ExecuteActivity(actx, activities.CreateChangeset, changeset, changeset.ID).Get(ctx, nil); err != nil {
		shared.Logger().Error("Error in creating changeset")
	}

	shared.Logger().Info("StackController", "Changeset created", changeset.ID)

	for idx, repo := range repos.Data {
		p := Instance().RepoProvider(repo.Provider) // get the specific provider

		if err := workflow.ExecuteActivity(pctx, p.TagCommit, repo.ProviderID, repoMarkers[idx].CommitID, changesetID.String(),
			"Tagged by quantm"); err != nil {
			shared.Logger().Error("StackController", "error in tagging the repo", repo.ProviderID, "err", err)
		}

		if err := workflow.ExecuteActivity(pctx, p.DeployChangeset, repo.ProviderID, changeset.ID).Get(ctx, nil); err != nil {
			shared.Logger().Error("StackController", "error in deploying for repo", repo.ProviderID, "err", err)
		}
	}

	shared.Logger().Debug("deployment done........")

	// // deployment map is designed to be used in OnPullRequestWorkflow only
	// logger := workflow.GetLogger(ctx)
	// lockID := "stack." + stackID // stack.<stack id>
	// deployments := make(Deployments)

	// // the idea is to save active infra which will be serving all the traffic and use this active infra as reference for next deployment
	// // this is not being used that as active infra for cloud run is being fetched from the cloud which is not an efficient approach
	// activeInfra := make(Infra)

	// // create and initialize mutex, initializing mutex will start a mutex workflow
	// logger.Info("creating mutex for stack", "stack", stackID)

	// lock := mutex.New(
	// 	mutex.WithCallerContext(ctx),
	// 	mutex.WithID(lockID),
	// )

	// if err := lock.Start(ctx); err != nil {
	// 	logger.Debug("unable to start mutex workflow", "error", err)
	// }

	// triggerChannel := workflow.GetSignalChannel(ctx, shared.WorkflowSignalDeploymentStarted.String())
	// assetsChannel := workflow.GetSignalChannel(ctx, WorkflowSignalAssetsRetrieved.String())
	// infraChannel := workflow.GetSignalChannel(ctx, WorkflowSignalInfraProvisioned.String())
	// deploymentChannel := workflow.GetSignalChannel(ctx, WorkflowSignalDeploymentCompleted.String())
	// manualOverrideChannel := workflow.GetSignalChannel(ctx, WorkflowSignalManualOverride.String())

	// selector := workflow.NewSelector(ctx)
	// selector.AddReceive(triggerChannel, onDeploymentStartedSignal(ctx, stackID, deployments))
	// selector.AddReceive(assetsChannel, onAssetsRetrievedSignal(ctx, stackID, deployments))
	// selector.AddReceive(infraChannel, onInfraProvisionedSignal(ctx, stackID, lock, deployments, activeInfra))
	// selector.AddReceive(deploymentChannel, onDeploymentCompletedSignal(ctx, stackID, deployments))
	// selector.AddReceive(manualOverrideChannel, onManualOverrideSignal(ctx, stackID, deployments))

	// // var prSignalsCounter int = 0
	// // return continue as new if this workflow has processed signals upto a limit
	// // if prSignalsCounter >= OnPullRequestWorkflowPRSignalsLimit {
	// // 	return workflow.NewContinueAsNewError(ctx, w.OnPullRequestWorkflow, stackID)
	// // }
	// for {
	// 	logger.Info("waiting for signals ....")
	// 	selector.Select(ctx)
	// }

	return nil
}

func (w *Workflows) EarlyDetection(ctx workflow.Context) error {
	shared.Logger().Debug("EarlyDetection", "waiting for signal", shared.WorkflowEarlyDetection.String())

	// get push event data via workflow signal
	ch := workflow.GetSignalChannel(ctx, shared.WorkflowEarlyDetection.String())

	signalPayload := &shared.PushEventSignal{}

	// receive signal payload
	ch.Receive(ctx, signalPayload)

	branchName := signalPayload.RefBranch
	installationID := signalPayload.InstallationID
	repoID := signalPayload.RepoID
	repoName := signalPayload.RepoName
	repoOwner := signalPayload.RepoOwner
	defaultBranch := signalPayload.DefaultBranch
	repoProvider := signalPayload.RepoProvider

	shared.Logger().Info("EarlyDetection", "signal payload", signalPayload)

	repoProviderInst := Instance().RepoProvider(RepoProvider(repoProvider))

	// msgProvider := Instance().MessageProvider(MessageProvider("slack"))
	msgProviderInst := Instance().MessageProvider(MessageProviderSlack) // TODO - maybe not hardcode to slack and get from payload

	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	pctx := workflow.WithActivityOptions(ctx, providerActOpts)

	// detect 200+ changes
	shared.Logger().Debug("going to detect 200+ changes")

	changes := 0
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.CalculateChangesInBranch, installationID, repoName, repoOwner,
		defaultBranch, branchName).Get(ctx, &changes); err != nil {
		shared.Logger().Error("EarlyDetection", "error from CalculateChangesInBranch activity", err)
		return err
	}

	if changes > 200 {
		message := "200+ lines changed on branch " + branchName

		if err := workflow.ExecuteActivity(
			pctx,
			msgProviderInst.SendChannelMessage,
			message,
		).Get(ctx, nil); err != nil {
			shared.Logger().Error("Error notifying Slack", "error", err.Error())
			return err
		}

		return nil
	} else {
		shared.Logger().Debug("200+ changes NOT detected")
	}

	// check merge conflicts
	shared.Logger().Debug("going to detect merge conflicts")

	latestDefBranchCommitSHA := ""
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.GetLatestCommit, strconv.FormatInt(repoID, 10), defaultBranch).Get(ctx,
		&latestDefBranchCommitSHA); err != nil {
		shared.Logger().Error("EarlyDetection", "error from GetLatestCommit activity", err)
		return err
	}

	// create a temp branch/ref
	tempBranchName := defaultBranch + "-tempcopy-for-target-" + branchName

	// delete the branch if it is present already
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.DeleteBranch, installationID, repoName, repoOwner,
		tempBranchName).Get(ctx, nil); err != nil {
		shared.Logger().Error("EarlyDetection", "error from DeleteBranch activity", err)
		return err
	}

	// create new ref
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.CreateBranch, installationID, repoID, repoName, repoOwner,
		latestDefBranchCommitSHA, tempBranchName).Get(ctx, nil); err != nil {
		shared.Logger().Error("EarlyDetection", "error from DeleteBranch activity", err)
		return err
	}

	if err := workflow.ExecuteActivity(pctx, repoProviderInst.MergeBranch, installationID, repoName, repoOwner, tempBranchName,
		branchName).Get(ctx, nil); err != nil {
		// dont want to retry this workflow so not returning error, just log and return
		shared.Logger().Error("EarlyDetection", "Error merging branch", err)

		message := "Merge Conflicts are expected on branch " + branchName
		if err = workflow.ExecuteActivity(
			pctx,
			msgProviderInst.SendChannelMessage,
			message,
		).Get(ctx, nil); err != nil {
			shared.Logger().Error("Error notifying Slack", "error", err.Error())
			return err
		}

		return nil
	}

	shared.Logger().Debug("merge conflicts NOT detected")

	// execute child workflow for stale detection
	shared.Logger().Debug("going to detect stale branch")

	latestDefBranchCommitSHA = ""
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.GetLatestCommit, strconv.FormatInt(repoID, 10), branchName).Get(ctx,
		&latestDefBranchCommitSHA); err != nil {
		shared.Logger().Error("EarlyDetection", "error from GetLatestCommit activity", err)
		return err
	}

	wf := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("repo"),
			shared.WithWorkflowBlockID(strconv.FormatInt(repoID, 10)),
			shared.WithWorkflowElement("branch"),
			shared.WithWorkflowElementID(branchName),
			shared.WithWorkflowProp("type", "Stale-Detection"),
		)

	if _, err := shared.Temporal().Client().ExecuteWorkflow(
		context.Background(),
		opts,
		wf.StaleBranchDetection,
		signalPayload, branchName, latestDefBranchCommitSHA); err != nil {
		// dont want to retry this workflow so not returning error, just log and return
		shared.Logger().Error("EarlyDetection", "error executing child workflow", err)
		return nil
	}

	return nil
}

func (w *Workflows) StaleBranchDetection(ctx workflow.Context, event *shared.PushEventSignal, branchName string,
	lastBranchCommit string) error {
	// Sleep for 5 days before raising stale detection
	_ = workflow.Sleep(ctx, 5*24*time.Hour)
	// _ = workflow.Sleep(ctx, 30*time.Second)

	shared.Logger().Debug("StaleBranchDetection", "woke up from sleep", "checking for stale branch")

	repoProviderInst := Instance().RepoProvider(RepoProvider(event.RepoProvider))

	// msgProvider := Instance().MessageProvider(MessageProvider("slack"))
	msgProviderInst := Instance().MessageProvider(MessageProviderSlack) // TODO - maybe not hardcode to slack and get from payload

	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	pctx := workflow.WithActivityOptions(ctx, providerActOpts)

	latestCommitSHA := ""
	if err := workflow.ExecuteActivity(pctx, repoProviderInst.GetLatestCommit, strconv.FormatInt(event.RepoID, 10), branchName).Get(ctx,
		&latestCommitSHA); err != nil {
		shared.Logger().Error("EarlyDetection", "error from GetLatestCommit activity", err)
		return err
	}

	// check if the branchName branch has the lastBranchCommit as the latest commit
	if lastBranchCommit == latestCommitSHA {
		message := "Stale branch " + branchName
		if err := workflow.ExecuteActivity(
			pctx,
			msgProviderInst.SendChannelMessage,
			message,
		).Get(ctx, nil); err != nil {
			shared.Logger().Error("Error notifying Slack", "error", err.Error())
			return err
		}

		return nil
	}

	// at this point, the branch is not stale so just return
	shared.Logger().Debug("stale branch NOT detected")

	return nil
}

// DeProvisionInfra de-provisions the infrastructure created for stack deployment.
func (w *Workflows) DeProvisionInfra(ctx workflow.Context, stackID string, resourceData *ResourceConfig) error {
	return nil
}

// GetAssets gets assets for stack including resources, workloads and blueprint.
func (w *Workflows) GetAssets(ctx workflow.Context, payload *GetAssetsPayload) error {
	var (
		future workflow.Future
		err    error = nil
	)

	logger := workflow.GetLogger(ctx)
	assets := NewAssets()
	workloads := SlicedResult[Workload]{}
	resources := SlicedResult[Resource]{}
	repos := SlicedResult[Repo]{}
	blueprint := Blueprint{}

	shared.Logger().Info("Get assets workflow")

	selector := workflow.NewSelector(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	providerActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
	}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get resources for stack
	future = workflow.ExecuteActivity(actx, activities.GetResources, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &resources); err != nil {
			logger.Error("GetResources providers failed", "error", err)
			return
		}
	})

	// get workloads for stack
	future = workflow.ExecuteActivity(actx, activities.GetWorkloads, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &workloads); err != nil {
			logger.Error("GetWorkloads providers failed", "error", err)
			return
		}
	})

	// get repos for stack
	future = workflow.ExecuteActivity(actx, activities.GetRepos, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &repos); err != nil {
			logger.Error("GetRepos providers failed", "error", err)
			return
		}
	})

	// get blueprint for stack
	future = workflow.ExecuteActivity(actx, activities.GetBluePrint, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &blueprint); err != nil {
			logger.Error("GetBluePrint providers failed", "error", err)
			return
		}
	})

	// TODO: come up with a better logic for this
	for i := 0; i < 4; i++ {
		selector.Select(ctx)
		// return if providers failed. TODO: handle race conditions as the 'err' variable is shared among all providers
		if err != nil {
			logger.Error("Exiting due to providers failure")
			return err
		}
	}

	// Tag the build image with changeset ID
	// TODO: update it to tag one or multiple images against every workload
	switch payload.ImageRegistry {
	case "GCPArtifactRegistry":
		err := workflow.ExecuteActivity(actx, activities.TagGcpImage, payload.Image, payload.Digest, payload.ChangeSetID).Get(ctx, nil)
		if err != nil {
			return err
		}
	default:
		shared.Logger().Error("This image registry is not supported in quantm yet", "registry", payload.ImageRegistry)
	}

	// get commits against the repos
	repoMarker := make([]ChangeSetRepoMarker, len(repos.Data))

	for idx, repo := range repos.Data {
		marker := &repoMarker[idx]
		// p := Instance().Provider(repo.Provider) // get the specific provider
		p := Instance().RepoProvider(repo.Provider) // get the specific provider
		commitID := ""

		if err := workflow.
			ExecuteActivity(pctx, p.GetLatestCommit, repo.ProviderID, repo.DefaultBranch).
			Get(ctx, &commitID); err != nil {
			logger.Error("Error in getting latest commit ID", "repo", repo.Name, "provider", repo.Provider)
			return fmt.Errorf("Error in getting latest commit ID repo:%s, provider:%s", repo.Name, repo.Provider.String())
		}

		marker.CommitID = commitID
		marker.HasChanged = repo.ID == payload.RepoID
		marker.Provider = repo.Provider.String()
		marker.RepoID = repo.ID.String()
		logger.Debug("Repo", "Name", repo.Name, "Repo marker", marker)
	}

	// save changeset
	stackID, _ := gocql.ParseUUID(payload.StackID)
	changeset := &ChangeSet{
		RepoMarkers: repoMarker,
		ID:          payload.ChangeSetID,
		StackID:     stackID,
	}

	err = workflow.ExecuteActivity(actx, activities.CreateChangeset, changeset, payload.ChangeSetID).Get(ctx, nil)
	if err != nil {
		logger.Error("Error in creating changeset")
	}

	assets.ChangesetID = payload.ChangeSetID
	assets.Blueprint = blueprint
	assets.Repos = append(assets.Repos, repos.Data...)
	assets.Resources = append(assets.Resources, resources.Data...)
	assets.Workloads = append(assets.Workloads, workloads.Data...)
	logger.Debug("Assets retrieved", "Assets", assets)

	// signal parent workflow
	parent := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.
		SignalExternalWorkflow(ctx, parent, "", WorkflowSignalAssetsRetrieved.String(), assets).
		Get(ctx, nil)

	return nil
}

// ProvisionInfra provisions the infrastructure required for stack deployment.
func (w *Workflows) ProvisionInfra(ctx workflow.Context, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	futures := make([]workflow.Future, 0)

	shared.Logger().Debug("provision infra", "assets", assets)

	for _, rsc := range assets.Resources {
		logger.Info("Creating resource", "Name", rsc.Name)

		// get the resource constructor specific to the driver e.g gke, cloudrun for GCP, sns, fargate for AWS
		resconstr := Instance().ResourceConstructor(rsc.Provider, rsc.Driver)

		if rsc.IsImmutable {
			// assuming a single region for now
			region := getRegion(rsc.Provider, &assets.Blueprint)
			providerConfig := getProviderConfig(rsc.Provider, &assets.Blueprint)
			r, err := resconstr.Create(rsc.Name, region, rsc.Config, providerConfig)

			if err != nil {
				logger.Error("could not create resource object", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}

			// resource is an interface and cannot be sent as a parameter in workflow because workflow cannot unmarshal an interface.
			// So we need to send the marshalled value in this workflow and then unmarshal and resconstruct the resource again in the
			// Deploy workflow
			ser, err := r.Marshal()
			if err != nil {
				logger.Error("could not marshal resource", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}

			assets.Infra[rsc.ID] = ser

			// TODO: initiate all resource provisions in parallel in child workflows and wait for all child workflows
			// completion before sending infra provisioned signal
			if f, err := r.Provision(ctx); err != nil {
				logger.Error("could not start resource provisioning", "ID", rsc.ID, "name", rsc.Name, "Error", err)

				return err
			} else if f != nil {
				futures = append(futures, f)
			}
		}
	}

	// not tested yet as provision workflow for cloudrun doesn't return any future.
	// the idea is to run all provision workflows asynchronously and in parallel and wait here for their completion before moving forward
	// also need to test if we can call child.Get with parent workflow's context
	for _, child := range futures {
		_ = child.Get(ctx, nil)
	}

	shared.Logger().Info("Signaling infra provisioned")

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalInfraProvisioned.String(), assets).Get(ctx, nil)

	return nil
}

// Deploy deploys the stack.
func (w *Workflows) Deploy(ctx workflow.Context, stackID string, lock *mutex.Lock, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Deployment initiated", "changeset", assets.ChangesetID, "infra", assets.Infra)
	infra := make(Infra)

	// Acquire lock
	err := lock.Acquire(ctx)
	if err != nil {
		logger.Error("Error in acquiring lock", "Error", err)
		return err
	}

	// create deployable, map of one or more workloads against each resource
	deployables := make(map[gocql.UUID][]Workload) // map of resource id and workloads
	for _, w := range assets.Workloads {
		_, ok := deployables[w.ResourceID]
		if !ok {
			deployables[w.ResourceID] = make([]Workload, 0)
		}

		deployables[w.ResourceID] = append(deployables[w.ResourceID], w)
	}

	// create the resource object again from marshaled data and deploy workload on resource
	for _, rsc := range assets.Resources {
		resconstr := Instance().ResourceConstructor(rsc.Provider, rsc.Driver)
		inf := assets.Infra[rsc.ID] // get marshaled resource from ID
		r := resconstr.CreateFromJson(inf)
		infra[rsc.ID] = r
		_ = r.Deploy(ctx, deployables[rsc.ID], assets.ChangesetID)
	}

	// update traffic on resource from 50 to 100
	var i int32
	for i = 50; i <= 100; i += 25 {
		for id, r := range infra {
			logger.Info("updating traffic", id, r)
			_ = r.UpdateTraffic(ctx, i)
		}
	}

	_ = lock.Release(ctx)

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalDeploymentCompleted.String(), assets)

	return nil
}

func onManualOverrideSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	triggerID := int64(0)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &triggerID)
		logger.Info("manual override for", "Trigger ID", triggerID)
	}
}

// onDeploymentStartedSignal is the channel handler for trigger channel
// It will execute GetAssets and update PR deployment state to "GettingAssets".
func onDeploymentStartedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		// Receive signal data
		payload := &shared.PullRequestSignal{}
		channel.Receive(ctx, payload)
		logger.Info("received deployment request", "Trigger ID", payload.TriggerID)

		// We want to filter workflows with changeset ID, so create changeset ID here and use it for creating workflow ID
		changesetID, _ := gocql.RandomUUID()

		// Set child workflow options
		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowElement("get_assets"),
				shared.WithWorkflowMod("trigger"),
				shared.WithWorkflowModID(strconv.FormatInt(payload.TriggerID, 10)),
			)

		getAssetsPayload := &GetAssetsPayload{
			StackID:       stackID,
			RepoID:        payload.RepoID,
			ChangeSetID:   changesetID,
			Digest:        payload.Digest,
			ImageRegistry: payload.ImageRegistry,
			Image:         payload.Image,
		}

		// execute GetAssets and wait until spawned
		var execution workflow.Execution

		cctx := workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(cctx, w.GetAssets, getAssetsPayload).
			GetChildWorkflowExecution().
			Get(cctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing GetAssets", "Error", err)
		}

		// create and save deployment data against a changeset
		deployment := NewDeployment()
		deployments[changesetID] = deployment
		deployment.state = GettingAssets
		deployment.workflows.GetAssets = execution.ID
	}
}

// onAssetsRetrievedSignal will receive assets sent by GetAssets, update deployment state and execute ProvisionInfra.
func onAssetsRetrievedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("received Assets", "changeset", assets.ChangesetID)

		// update deployment state
		deployment := deployments[assets.ChangesetID]
		deployment.state = GotAssets

		// execute provision infra workflow
		logger.Info("Executing provision Infra workflow")

		var execution workflow.Execution

		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("provision_infra"),
			)

		ctx = workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(ctx, w.ProvisionInfra, assets).
			GetChildWorkflowExecution().Get(ctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing ProvisionInfra", "Error", err)
		}

		logger.Info("Executed provision Infra workflow")

		deployment.state = ProvisioningInfra
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onInfraProvisionedSignal will receive assets by ProvisionInfra, update deployment state and execute Deploy.
func onInfraProvisionedSignal(
	ctx workflow.Context, stackID string, lock mutex.Mutex, deployments Deployments, activeinfra Infra,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Infra provisioned", "changeset", assets.ChangesetID)

		deployment := deployments[assets.ChangesetID]
		deployment.state = InfraProvisioned

		// deployment.OldInfra = activeinfra  // All traffic is currently being routed to this infra
		deployment.NewInfra = assets.Infra // handling zero traffic, no workload is deployed

		var execution workflow.Execution

		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("deploy"),
			)
		ctx = workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(ctx, w.Deploy, stackID, lock.(*mutex.Lock), assets).
			GetChildWorkflowExecution().Get(ctx, &execution)
		if err != nil {
			logger.Error("Error in Executing deployment workflow", "Error", err)
		}

		deployment.state = CreatingDeployment
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onDeploymentCompletedSignal will conclude the deployment.
func onDeploymentCompletedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Deployment complete", "changeset", assets.ChangesetID)
		delete(deployments, assets.ChangesetID)

		logger.Info("Deleted deployment data", "changeset", assets.ChangesetID)
	}
}
