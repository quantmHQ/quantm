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

	"github.com/gocql/gocql"
	gh "github.com/google/go-github/v53/github"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Activities groups all the activities for the github provider.
	Activities struct{}
)

// CreateOrUpdateInstallation creates or update the Installation.
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *Installation) (*Installation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status ...")

		installation.Status = payload.Status
	} else {
		log.Info("installation not found, creating ...")
		log.Debug("payload", "payload", payload)

		installation = payload
	}

	log.Info("saving installation ...")

	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", "error", err)
		return installation, err
	}

	return installation, nil
}

// CreateOrUpdateGithubRepo creates a single row for Repo.
func (a *Activities) CreateOrUpdateGithubRepo(ctx context.Context, payload *Repo) error {
	log := activity.GetLogger(ctx)
	repo, err := a.GetGithubRepo(ctx, payload)

	// if we get the repo, the error will be nil
	if err == nil {
		log.Info("repository found, updating ...")
	} else {
		log.Info("repository not found, creating ...")
		log.Debug("payload", "payload", payload)
	}

	if err := db.Save(repo); err != nil {
		log.Error("error saving repository ...", "error", err)
		return err
	}

	return nil
}

// GetGithubRepo gets Repo against given Repo.
func (a *Activities) GetGithubRepo(ctx context.Context, payload *Repo) (*Repo, error) {
	repo := &Repo{}
	params := db.QueryParams{
		"name":      "'" + payload.Name + "'",
		"full_name": "'" + payload.FullName + "'",
		"github_id": strconv.FormatInt(payload.GithubID, 10),
		"team_id":   payload.TeamID.String(),
	}

	if err := db.Get(repo, params); err != nil {
		return payload, err
	}

	return repo, nil
}

// GetInstallation gets Installation against given installation_id.
func (a *Activities) GetInstallation(ctx context.Context, id int64) (*Installation, error) {
	installation := &Installation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": strconv.FormatInt(id, 10)}); err != nil {
		return installation, err
	}

	return installation, nil
}

// GetCoreRepo gets entity.Repo against given Repo.
func (a *Activities) GetCoreRepo(ctx context.Context, repo *Repo) (*core.Repo, error) {
	r := &core.Repo{}

	// TODO: add provider name in query
	params := db.QueryParams{
		"provider_id": "'" + strconv.FormatInt(repo.GithubID, 10) + "'",
		"provider":    "'github'",
	}

	// shared.Logger().Debug("GetCoreRepo", "params", params)

	if err := db.Get(r, params); err != nil {
		return r, err
	}

	return r, nil
}

// GetCoreRepo gets entity.Stack against given core Repo.
func (a *Activities) GetStack(ctx context.Context, repo *core.Repo) (*core.Stack, error) {
	s := &core.Stack{}

	params := db.QueryParams{
		"id": repo.StackID.String(),
	}

	if err := db.Get(s, params); err != nil {
		return s, err
	}

	return s, nil
}

// GetLatestCommit gets latest commit for default branch of the provided repo.
func (a *Activities) GetLatestCommit(ctx context.Context, providerID string, branch string) (string, error) {
	logger := activity.GetLogger(ctx)
	prepo := &Repo{}

	if err := db.Get(prepo, db.QueryParams{"github_id": providerID}); err != nil {
		return "", err
	}

	client, err := Instance().GetClientFromInstallation(prepo.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return "", err
	}

	gb, _, err := client.Repositories.GetBranch(context.Background(), strings.Split(prepo.FullName, "/")[0], prepo.Name, branch, false)
	if err != nil {
		logger.Error("GetBranch for Github Repo failed", "Error", err)
		return "", err
	}

	logger.Debug("Repo", "Name", prepo.FullName, "Branch name", gb.Name, "Last commit", gb.Commit.SHA)

	return *gb.Commit.SHA, nil
}

// TODO - break it to smalller activities (create, delete and merge).
func (a *Activities) RebaseAndMerge(ctx context.Context, repoOwner string, repoName string,
	targetBranchName string, installationID int64) (string, error) {
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err.Error(), err
	}

	// Get the default branch (e.g., "main")
	repo, _, err := client.Repositories.Get(ctx, repoOwner, repoName)
	if err != nil {
		shared.Logger().Error("RebaseAndMerge Activity", "Error getting repository: ", err)
		return err.Error(), err
	}

	defaultBranch := *repo.DefaultBranch
	newBranchName := defaultBranch + "-tempcopy-for-target-" + targetBranchName

	// Get the latest commit SHA of the default branch
	commits, _, err := client.Repositories.ListCommits(ctx, repoOwner, repoName, &gh.CommitsListOptions{
		SHA: defaultBranch,
	})
	if err != nil {
		shared.Logger().Error("RebaseAndMerge Activity", "Error getting commits: ", err)
		return err.Error(), err
	}

	// // Use the latest commit SHA
	// if len(commits) == 0 {
	// 	shared.Logger().Error("RebaseAndMerge Activity", "No commits found in the default branch.", nil)
	// 	return err.Error(), err
	// }

	latestCommitSHA := *commits[0].SHA

	// Create a new branch based on the latest commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + newBranchName),
		Object: &gh.GitObject{
			SHA: &latestCommitSHA,
		},
	}

	_, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, ref)
	if err != nil {
		shared.Logger().Error("RebaseAndMerge Activity", "Error creating branch: ", err)
		return err.Error(), err
	}

	shared.Logger().Info("RebaseAndMerge Activity", "Branch created successfully: ", newBranchName)

	// Perform rebase of the target branch with the new branch
	rebaseRequest := &gh.RepositoryMergeRequest{
		Base:          &newBranchName,
		Head:          &targetBranchName,
		CommitMessage: gh.String("Rebasing " + targetBranchName + " with " + newBranchName),
	}

	_, _, err = client.Repositories.Merge(ctx, repoOwner, repoName, rebaseRequest)
	if err != nil {
		shared.Logger().Error("RebaseAndMerge Activity", "Error rebasing branches: ", err)
		return err.Error(), err
	}

	shared.Logger().Info("RebaseAndMerge Activity", "status",
		fmt.Sprintf("Branch %s rebased with %s successfully.\n", targetBranchName, newBranchName))

	// Perform rebase of the new branch with the main branch
	rebaseRequest = &gh.RepositoryMergeRequest{
		Base:          &defaultBranch,
		Head:          &newBranchName,
		CommitMessage: gh.String("Rebasing " + newBranchName + " with " + defaultBranch),
	}

	repoCommit, _, err := client.Repositories.Merge(ctx, repoOwner, repoName, rebaseRequest)
	if err != nil {
		shared.Logger().Error("RebaseAndMerge Activity", "Error rebasing branches: ", err)
		return err.Error(), err
	}

	shared.Logger().Info("RebaseAndMerge Activity", "status",
		fmt.Sprintf("Branch %s rebased with %s successfully.\n", newBranchName, defaultBranch))

	return *repoCommit.SHA, nil
}

func (a *Activities) TriggerCIAction(ctx context.Context, installationID int64, repoOwner string,
	repoName string, targetBranch string) error {
	shared.Logger().Debug("activity TriggerGithubAction started")

	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	workflowName := "cicd_quantm.yaml" //TODO: either fix this or obtain it somehow

	paylod := gh.CreateWorkflowDispatchEventRequest{
		Ref: targetBranch,
		Inputs: map[string]any{
			"target-branch": targetBranch,
		},
	}

	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, repoOwner, repoName, workflowName, paylod)
	if err != nil {
		shared.Logger().Error("TriggerGithubAction", "Error triggering workflow:", err)
		return err
	}

	shared.Logger().Debug("TriggerGithubAction", "response", res)

	return nil
}

func (a *Activities) DeployChangeset(ctx context.Context, repoID string, changesetID *gocql.UUID) error {
	shared.Logger().Debug("DeployChangeset", "github activity DeployChangeset started for changeset", changesetID)

	gh_action_name := "deploy_quantm.yaml" //TODO: fixed it for now

	// get installationID, repoName, repoOwner from github_repos table
	githubRepo := &Repo{}
	params := db.QueryParams{
		"github_id": repoID,
	}

	if err := db.Get(githubRepo, params); err != nil {
		return err
	}

	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	paylod := gh.CreateWorkflowDispatchEventRequest{
		Ref: "main",
		Inputs: map[string]any{
			"changesetId": changesetID,
		},
	}

	var repoOwner, repoName string

	parts := strings.Split(githubRepo.FullName, "/")

	if len(parts) == 2 {
		repoOwner = parts[0]
		repoName = parts[1]
	}

	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, repoOwner, repoName, gh_action_name, paylod)
	if err != nil {
		shared.Logger().Error("DeployChangeset", "Error triggering workflow:", err)
		return err
	}

	shared.Logger().Debug("DeployChangeset", "response", res)

	return nil
}

func (a *Activities) TagCommit(ctx context.Context, repoID string, commitSHA string, tagName string, tagMessage string) error {
	// get installationID, repoName, repoOwner from github_repos table
	githubRepo := &Repo{}
	params := db.QueryParams{
		"github_id": repoID,
	}

	if err := db.Get(githubRepo, params); err != nil {
		return err
	}

	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	tag := &gh.Tag{
		Tag:     &tagName,
		Message: &tagMessage,
		Object: &gh.GitObject{
			SHA:  &commitSHA,
			Type: gh.String("commit"), // Specify the type of object being tagged
		},
	}

	var repoOwner, repoName string

	parts := strings.Split(githubRepo.FullName, "/")

	if len(parts) == 2 {
		repoOwner = parts[0]
		repoName = parts[1]
	}

	_, _, err = client.Git.CreateTag(ctx, repoOwner, repoName, tag)
	if err != nil {
		shared.Logger().Error("TagCommit", "Error creating tag", err)
	}

	// Push the tag to the remote repository
	ref := "refs/tags/" + tagName
	if _, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, &gh.Reference{
		Ref:    &ref,
		Object: &gh.GitObject{SHA: &commitSHA},
	}); err != nil {
		shared.Logger().Error("TagCommit", "Error pushing tag to remote repository", err)
	}

	return nil
}

func (a *Activities) DeleteBranch(ctx context.Context, installationID int64, repoName string, repoOwner string, branchName string) error {
	// Get github client
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Create a new branch based on the latest defaultBranch commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + branchName),
	}

	// delete the temp ref if its present before
	if _, err := client.Git.DeleteRef(context.Background(), repoOwner, repoName, *ref.Ref); err != nil {
		shared.Logger().Error("DeleteBranch", "Error deleting ref"+*ref.Ref, err)

		if strings.Contains(err.Error(), "422 Reference does not exist") {
			// if a ref doesnt exist already, dont return the error
			return nil
		}

		return err
	}

	return nil
}

func (a *Activities) CreateBranch(ctx context.Context, installationID int64, repoID int64, repoName string, repoOwner string,
	targetCommit string, newBranchName string) error {
	// Get github client
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Create a new branch based on the latest defaultBranch commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + newBranchName),
		Object: &gh.GitObject{
			SHA: &targetCommit,
		},
	}

	// create new ref
	if _, _, err = client.Git.CreateRef(context.Background(), repoOwner, repoName, ref); err != nil {
		shared.Logger().Error("CreateBranch activity", "Error creating branch: ", err)

		// dont want to retry this workflow so not returning error, just log and return
		return nil
	}

	return nil
}

func (a *Activities) MergeBranch(ctx context.Context, installationID int64, repoName string, repoOwner string, baseBranch string,
	targetBranch string) error {
	// Get github client for operations
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// targetBranch will be merged into the baseBranch
	rebaseReq := &gh.RepositoryMergeRequest{
		Base:          &baseBranch,
		Head:          &targetBranch,
		CommitMessage: gh.String("Rebasing " + targetBranch + " with " + baseBranch),
	}

	if _, _, err := client.Repositories.Merge(context.Background(), repoOwner, repoName, rebaseReq); err != nil {
		shared.Logger().Error("Merge failed", "Error", err)
		return err
	}

	return nil
}

func (a *Activities) CalculateChangesInBranch(ctx context.Context, installationID int64, repoName string, repoOwner string,
	defaultBranch string, targetBranch string) (int, error) {
	// Get github client for operations
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return -1, err
	}

	comparison, _, err := client.Repositories.CompareCommits(context.Background(), repoOwner, repoName, defaultBranch, targetBranch, nil)
	if err != nil {
		shared.Logger().Error("Error in CalculateChangesInBranch", "CompareCommits", err)
		return -1, err
	}

	var changes int
	for _, file := range comparison.Files {
		changes += file.GetChanges()
	}

	shared.Logger().Debug("CalculateChangesInBranch", "total changes in branch "+targetBranch, changes)

	return changes, nil
}

func (a *Activities) GetAllBranches(ctx context.Context, installationID int64, repoName string, repoOwner string) ([]string, error) {
	// get github client
	client, err := Instance().GetClientFromInstallation(installationID)
	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return nil, err
	}

	var branchNames []string

	page := 1

	for {
		branches, resp, err := client.Repositories.ListBranches(ctx, repoOwner, repoName, &gh.BranchListOptions{
			ListOptions: gh.ListOptions{
				Page:    page,
				PerPage: 30, // Adjust this value as needed
			},
		})
		if err != nil {
			shared.Logger().Error("GetAllBranches: could not get branches", "Error", err)
			return nil, err
		}

		for _, branch := range branches {
			branchNames = append(branchNames, *branch.Name)
		}

		// Check if there are more pages to fetch
		if resp.NextPage == 0 {
			break // No more pages
		}

		page = resp.NextPage
	}

	return branchNames, nil
}

func (a *Activities) GetRepoTeamID(ctx context.Context, repoID string) (string, error) {
	logger := activity.GetLogger(ctx)
	prepo := &Repo{}

	if err := db.Get(prepo, db.QueryParams{"github_id": repoID}); err != nil {
		logger.Error("GetRepoTeamID failed", "Error", err)
		return "", err
	}

	shared.Logger().Info("GetRepoTeamID Activity", "Get Repo Team ID successfully: ", prepo.TeamID)

	return prepo.TeamID.String(), nil
}
