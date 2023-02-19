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
	"time"
)

// Github embedded types.
type (
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	InstallationID struct {
		ID int64 `json:"id"`
	}

	Commit struct {
		SHA       string      `json:"sha"`
		ID        string      `json:"id"`
		NodeID    string      `json:"node_id"`
		TreeID    string      `json:"tree_id"`
		Distinct  bool        `json:"distinct"`
		Message   string      `json:"message"`
		Timestamp string      `json:"timestamp"`
		URL       string      `json:"url"`
		Author    PartialUser `json:"author"`
		Committer PartialUser `json:"committer"`
		Added     []string    `json:"added"`
		Removed   []string    `json:"removed"`
		Modified  []string    `json:"modified"`
	}

	HeadCommit struct {
		ID        string      `json:"id"`
		NodeID    string      `json:"node_id"`
		TreeID    string      `json:"tree_id"`
		Distinct  bool        `json:"distinct"`
		Message   string      `json:"message"`
		Timestamp string      `json:"timestamp"`
		URL       string      `json:"url"`
		Author    PartialUser `json:"author"`
		Committer PartialUser `json:"committer"`
		Added     []string    `json:"added"`
		Removed   []string    `json:"removed"`
		Modified  []string    `json:"modified"`
	}

	PartialUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	Repository struct {
		ID               int64     `json:"id"`
		NodeID           string    `json:"node_id"`
		Name             string    `json:"name"`
		FullName         string    `json:"full_name"`
		Owner            User      `json:"owner"`
		Private          bool      `json:"private"`
		HTMLUrl          string    `json:"html_url"`
		Description      string    `json:"description"`
		Fork             bool      `json:"fork"`
		URL              string    `json:"url"`
		ForksURL         string    `json:"forks_url"`
		KeysURL          string    `json:"keys_url"`
		CollaboratorsURL string    `json:"collaborators_url"`
		TeamsURL         string    `json:"teams_url"`
		HooksURL         string    `json:"hooks_url"`
		IssueEventsURL   string    `json:"issue_events_url"`
		EventsURL        string    `json:"events_url"`
		AssigneesURL     string    `json:"assignees_url"`
		BranchesURL      string    `json:"branches_url"`
		TagsURL          string    `json:"tags_url"`
		BlobsURL         string    `json:"blobs_url"`
		GitTagsURL       string    `json:"git_tags_url"`
		GitRefsURL       string    `json:"git_refs_url"`
		TreesURL         string    `json:"trees_url"`
		StatusesURL      string    `json:"statuses_url"`
		LanguagesURL     string    `json:"languages_url"`
		StargazersURL    string    `json:"stargazers_url"`
		ContributorsURL  string    `json:"contributors_url"`
		SubscribersURL   string    `json:"subscribers_url"`
		SubscriptionURL  string    `json:"subscription_url"`
		CommitsURL       string    `json:"commits_url"`
		GitCommitsURL    string    `json:"git_commits_url"`
		CommentsURL      string    `json:"comments_url"`
		IssueCommentURL  string    `json:"issue_comment_url"`
		ContentsURL      string    `json:"contents_url"`
		CompareURL       string    `json:"compare_url"`
		MergesURL        string    `json:"merges_url"`
		ArchiveURL       string    `json:"archive_url"`
		DownloadsURL     string    `json:"downloads_url"`
		IssuesURL        string    `json:"issues_url"`
		PullsURL         string    `json:"pulls_url"`
		MilestonesURL    string    `json:"milestones_url"`
		NotificationsURL string    `json:"notifications_url"`
		LabelsURL        string    `json:"labels_url"`
		ReleasesURL      string    `json:"releases_url"`
		CreatedAt        int64     `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		PushedAt         int64     `json:"pushed_at"`
		GitURL           string    `json:"git_url"`
		SSHUrl           string    `json:"ssh_url"`
		CloneURL         string    `json:"clone_url"`
		SvnURL           string    `json:"svn_url"`
		Homepage         *string   `json:"homepage"`
		Size             int64     `json:"size"`
		StargazersCount  int64     `json:"stargazers_count"`
		WatchersCount    int64     `json:"watchers_count"`
		Language         *string   `json:"language"`
		HasIssues        bool      `json:"has_issues"`
		HasDownloads     bool      `json:"has_downloads"`
		HasWiki          bool      `json:"has_wiki"`
		HasPages         bool      `json:"has_pages"`
		ForksCount       int64     `json:"forks_count"`
		MirrorURL        *string   `json:"mirror_url"`
		OpenIssuesCount  int64     `json:"open_issues_count"`
		Forks            int64     `json:"forks"`
		OpenIssues       int64     `json:"open_issues"`
		Watchers         int64     `json:"watchers"`
		DefaultBranch    string    `json:"default_branch"`
		Stargazers       int64     `json:"stargazers"`
		MasterBranch     string    `json:"master_branch"`
	}

	PullRequestRepository struct {
		ID               int64     `json:"id"`
		NodeID           string    `json:"node_id"`
		Name             string    `json:"name"`
		FullName         string    `json:"full_name"`
		Owner            User      `json:"owner"`
		Private          bool      `json:"private"`
		HTMLUrl          string    `json:"html_url"`
		Description      string    `json:"description"`
		Fork             bool      `json:"fork"`
		URL              string    `json:"url"`
		ForksURL         string    `json:"forks_url"`
		KeysURL          string    `json:"keys_url"`
		CollaboratorsURL string    `json:"collaborators_url"`
		TeamsURL         string    `json:"teams_url"`
		HooksURL         string    `json:"hooks_url"`
		IssueEventsURL   string    `json:"issue_events_url"`
		EventsURL        string    `json:"events_url"`
		AssigneesURL     string    `json:"assignees_url"`
		BranchesURL      string    `json:"branches_url"`
		TagsURL          string    `json:"tags_url"`
		BlobsURL         string    `json:"blobs_url"`
		GitTagsURL       string    `json:"git_tags_url"`
		GitRefsURL       string    `json:"git_refs_url"`
		TreesURL         string    `json:"trees_url"`
		StatusesURL      string    `json:"statuses_url"`
		LanguagesURL     string    `json:"languages_url"`
		StargazersURL    string    `json:"stargazers_url"`
		ContributorsURL  string    `json:"contributors_url"`
		SubscribersURL   string    `json:"subscribers_url"`
		SubscriptionURL  string    `json:"subscription_url"`
		CommitsURL       string    `json:"commits_url"`
		GitCommitsURL    string    `json:"git_commits_url"`
		CommentsURL      string    `json:"comments_url"`
		IssueCommentURL  string    `json:"issue_comment_url"`
		ContentsURL      string    `json:"contents_url"`
		CompareURL       string    `json:"compare_url"`
		MergesURL        string    `json:"merges_url"`
		ArchiveURL       string    `json:"archive_url"`
		DownloadsURL     string    `json:"downloads_url"`
		IssuesURL        string    `json:"issues_url"`
		PullsURL         string    `json:"pulls_url"`
		MilestonesURL    string    `json:"milestones_url"`
		NotificationsURL string    `json:"notifications_url"`
		LabelsURL        string    `json:"labels_url"`
		ReleasesURL      string    `json:"releases_url"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		PushedAt         time.Time `json:"pushed_at"`
		GitURL           string    `json:"git_url"`
		SSHUrl           string    `json:"ssh_url"`
		CloneURL         string    `json:"clone_url"`
		SvnURL           string    `json:"svn_url"`
		Homepage         *string   `json:"homepage"`
		Size             int64     `json:"size"`
		StargazersCount  int64     `json:"stargazers_count"`
		WatchersCount    int64     `json:"watchers_count"`
		Language         *string   `json:"language"`
		HasIssues        bool      `json:"has_issues"`
		HasDownloads     bool      `json:"has_downloads"`
		HasWiki          bool      `json:"has_wiki"`
		HasPages         bool      `json:"has_pages"`
		ForksCount       int64     `json:"forks_count"`
		MirrorURL        *string   `json:"mirror_url"`
		OpenIssuesCount  int64     `json:"open_issues_count"`
		Forks            int64     `json:"forks"`
		OpenIssues       int64     `json:"open_issues"`
		Watchers         int64     `json:"watchers"`
		DefaultBranch    string    `json:"default_branch"`
		Stargazers       int64     `json:"stargazers"`
		MasterBranch     string    `json:"master_branch"`
	}

	User struct {
		Login             string `json:"login"`
		ID                int64  `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLUrl           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	}

	Permissions struct {
		Issues             string `json:"issues"`
		Metadata           string `json:"metadata"`
		PullRequests       string `json:"pull_requests"`
		RepositoryProjects string `json:"repository_projects"`
	}

	// Milestone contains GitHub's milestone information.
	Milestone struct {
		URL          string    `json:"url"`
		HTMLUrl      string    `json:"html_url"`
		LabelsURL    string    `json:"labels_url"`
		ID           int64     `json:"id"`
		NodeID       string    `json:"node_id"`
		Number       int64     `json:"number"`
		State        string    `json:"state"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Creator      User      `json:"creator"`
		OpenIssues   int64     `json:"open_issues"`
		ClosedIssues int64     `json:"closed_issues"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		ClosedAt     time.Time `json:"closed_at"`
		DueOn        time.Time `json:"due_on"`
	}

	FullInstallationPayload struct {
		ID                  int64       `json:"id"`
		NodeID              string      `json:"node_id"`
		Account             User        `json:"account"`
		RepositorySelection string      `json:"repository_selection"`
		AccessTokensURL     string      `json:"access_tokens_url"`
		RepositoriesURL     string      `json:"repositories_url"`
		HTMLUrl             string      `json:"html_url"`
		AppID               int         `json:"app_id"`
		TargetID            int         `json:"target_id"`
		TargetType          string      `json:"target_type"`
		Permissions         Permissions `json:"permissions"`
		Events              []string    `json:"events"`
		CreatedAt           time.Time   `json:"created_at"`
		UpdatedAt           time.Time   `json:"updated_at"`
		SingleFileName      *string     `json:"single_file_name"`
	}

	PartialRepository struct {
		ID       int64  `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
	}

	PullRequestHead struct {
		Label string                `json:"label"`
		Ref   string                `json:"ref"`
		SHA   string                `json:"sha"`
		User  User                  `json:"user"`
		Repo  PullRequestRepository `json:"repo"`
	}

	PullRequestBase struct {
		Label string                `json:"label"`
		Ref   string                `json:"ref"`
		SHA   string                `json:"sha"`
		User  User                  `json:"user"`
		Repo  PullRequestRepository `json:"repo"`
	}

	Href struct {
		Href string `json:"href"`
	}

	PullRequestLinks struct {
		Self           Href `json:"self"`
		HTML           Href `json:"html"`
		Issue          Href `json:"issue"`
		Comments       Href `json:"comments"`
		ReviewComments Href `json:"review_comments"`
		ReviewComment  Href `json:"review_comment"`
		Commits        Href `json:"commits"`
		Statuses       Href `json:"statuses"`
	}

	PullRequest struct {
		URL                string           `json:"url"`
		ID                 int64            `json:"id"`
		NodeID             string           `json:"node_id"`
		HTMLUrl            string           `json:"html_url"`
		DiffURL            string           `json:"diff_url"`
		PatchURL           string           `json:"patch_url"`
		IssueURL           string           `json:"issue_url"`
		Number             int64            `json:"number"`
		State              string           `json:"state"`
		Locked             bool             `json:"locked"`
		Title              string           `json:"title"`
		User               User             `json:"user"`
		Body               string           `json:"body"`
		CreatedAt          time.Time        `json:"created_at"`
		UpdatedAt          time.Time        `json:"updated_at"`
		ClosedAt           *time.Time       `json:"closed_at"`
		MergedAt           *time.Time       `json:"merged_at"`
		MergeCommitSha     *string          `json:"merge_commit_sha"`
		Assignee           *User            `json:"assignee"`
		Assignees          []*User          `json:"assignees"`
		Milestone          *Milestone       `json:"milestone"`
		Draft              bool             `json:"draft"`
		CommitsURL         string           `json:"commits_url"`
		ReviewCommentsURL  string           `json:"review_comments_url"`
		ReviewCommentURL   string           `json:"review_comment_url"`
		CommentsURL        string           `json:"comments_url"`
		StatusesURL        string           `json:"statuses_url"`
		RequestedReviewers []User           `json:"requested_reviewers,omitempty"`
		Labels             []Label          `json:"labels"`
		Head               PullRequestHead  `json:"head"`
		Base               PullRequestBase  `json:"base"`
		Links              PullRequestLinks `json:"_links"`
		Merged             bool             `json:"merged"`
		Mergeable          *bool            `json:"mergeable"`
		MergeableState     string           `json:"mergeable_state"`
		MergedBy           *User            `json:"merged_by"`
		Comments           int64            `json:"comments"`
		ReviewComments     int64            `json:"review_comments"`
		Commits            int64            `json:"commits"`
		Additions          int64            `json:"additions"`
		Deletions          int64            `json:"deletions"`
		ChangedFiles       int64            `json:"changed_files"`
	}

	Label struct {
		ID          int64  `json:"id"`
		NodeID      string `json:"node_id"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Name        string `json:"name"`
		Color       string `json:"color"`
		Default     bool   `json:"default"`
	}

	RequestedTeam struct {
		Name            string `json:"name"`
		ID              int64  `json:"id"`
		NodeID          string `json:"node_id"`
		Slug            string `json:"slug"`
		Description     string `json:"description"`
		Privacy         string `json:"privacy"`
		URL             string `json:"url"`
		HTMLUrl         string `json:"html_url"`
		MembersURL      string `json:"members_url"`
		RepositoriesURL string `json:"repositories_url"`
		Permission      string `json:"permission"`
	}

	Organization struct {
		Login            string `json:"login"`
		ID               int    `json:"id"`
		NodeID           string `json:"node_id"`
		URL              string `json:"url"`
		ReposURL         string `json:"repos_url"`
		EventsURL        string `json:"events_url"`
		HooksURL         string `json:"hooks_url"`
		IssuesURL        string `json:"issues_url"`
		MembersURL       string `json:"members_url"`
		PublicMembersURL string `json:"public_members_url"`
		AvatarURL        string `json:"avatar_url"`
		Description      string `json:"description"`
	}
)
