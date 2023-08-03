// Package github provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen.
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package github

import (
	"encoding/json"
	"errors"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidSetupAction    = errors.New("invalid SetupAction value")
	ErrInvalidWorkflowStatus = errors.New("invalid WorkflowStatus value")
)

type (
	SetupActionMapType map[string]SetupAction // SetupActionMapType is a quick lookup map for SetupAction.
)

// Defines values for SetupAction.
const (
	SetupActionCreated SetupAction = "created"
	SetupActionDeleted SetupAction = "deleted"
	SetupActionUpdated SetupAction = "updated"
)

// SetupActionValues returns all known values for SetupAction.
var (
	SetupActionMap = SetupActionMapType{
		SetupActionCreated.String(): SetupActionCreated,
		SetupActionDeleted.String(): SetupActionDeleted,
		SetupActionUpdated.String(): SetupActionUpdated,
	}
)

/*
 * Helper methods for SetupAction for easy marshalling and unmarshalling.
 */
func (v SetupAction) String() string               { return string(v) }
func (v SetupAction) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *SetupAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := SetupActionMap[s]
	if !ok {
		return ErrInvalidSetupAction
	}

	*v = val

	return nil
}

type (
	WorkflowStatusMapType map[string]WorkflowStatus // WorkflowStatusMapType is a quick lookup map for WorkflowStatus.
)

// Defines values for WorkflowStatus.
const (
	WorkflowStatusFailure  WorkflowStatus = "failure"
	WorkflowStatusQueued   WorkflowStatus = "queued"
	WorkflowStatusSignaled WorkflowStatus = "signaled"
	WorkflowStatusSkipped  WorkflowStatus = "skipped"
	WorkflowStatusSuccess  WorkflowStatus = "success"
)

// WorkflowStatusValues returns all known values for WorkflowStatus.
var (
	WorkflowStatusMap = WorkflowStatusMapType{
		WorkflowStatusFailure.String():  WorkflowStatusFailure,
		WorkflowStatusQueued.String():   WorkflowStatusQueued,
		WorkflowStatusSignaled.String(): WorkflowStatusSignaled,
		WorkflowStatusSkipped.String():  WorkflowStatusSkipped,
		WorkflowStatusSuccess.String():  WorkflowStatusSuccess,
	}
)

/*
 * Helper methods for WorkflowStatus for easy marshalling and unmarshalling.
 */
func (v WorkflowStatus) String() string               { return string(v) }
func (v WorkflowStatus) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *WorkflowStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := WorkflowStatusMap[s]
	if !ok {
		return ErrInvalidWorkflowStatus
	}

	*v = val

	return nil
}

// ArtifactReadyRequest defines model for ArtifactReadyRequest.
type ArtifactReadyRequest struct {
	Image          string `cql:"image" json:"image"`
	InstallationID string `cql:"installation_id" json:"installation_id"`
	PullRequestID  string `cql:"pull_request_id" json:"pull_request_id"`
	RepoID         string `cql:"repo_id" json:"repo_id"`
}

var (
	artifactreadyrequestColumns = []string{"image", "installation_id", "pull_request_id", "repo_id"}

	artifactreadyrequestMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_artifact",
			Columns: artifactreadyrequestColumns,
		},
	}

	artifactreadyrequestTable = itable.New(*artifactreadyrequestMeta.M)
)

func (artifactreadyrequest *ArtifactReadyRequest) GetTable() itable.ITable {
	return artifactreadyrequestTable
}

// CompleteInstallationRequest complete the installation given the installation_id & setup_action
type CompleteInstallationRequest struct {
	InstallationID int64       `json:"installation_id"`
	SetupAction    SetupAction `json:"setup_action"`
}

// Installation defines model for GithubInstallation.
type Installation struct {
	CreatedAt         time.Time  `cql:"created_at" json:"created_at"`
	ID                gocql.UUID `cql:"id" json:"id"`
	InstallationID    int64      `cql:"installation_id" json:"installation_id" validate:"required,db_unique"`
	InstallationLogin string     `cql:"installation_login" json:"installation_login"`
	InstallationType  string     `cql:"installation_type" json:"installation_type"`
	SenderID          int64      `cql:"sender_id" json:"sender_id"`
	SenderLogin       string     `cql:"sender_login" json:"sender_login"`
	Status            string     `cql:"status" json:"status"`
	TeamID            gocql.UUID `cql:"team_id" json:"team_id"`
	UpdatedAt         time.Time  `cql:"updated_at" json:"updated_at"`
}

var (
	githubinstallationColumns = []string{"created_at", "id", "installation_id", "installation_login", "installation_type", "sender_id", "sender_login", "status", "team_id", "updated_at"}

	githubinstallationMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_installations",
			Columns: githubinstallationColumns,
		},
	}

	githubinstallationTable = itable.New(*githubinstallationMeta.M)
)

func (githubinstallation *Installation) GetTable() itable.ITable {
	return githubinstallationTable
}

// Repo defines model for GithubRepo.
type Repo struct {
	CreatedAt      time.Time  `cql:"created_at" json:"created_at"`
	FullName       string     `cql:"full_name" json:"full_name"`
	GithubID       int64      `cql:"github_id" json:"github_id"`
	ID             gocql.UUID `cql:"id" json:"id"`
	InstallationID int64      `cql:"installation_id" json:"installation_id"`
	Name           string     `cql:"name" json:"name"`
	TeamID         gocql.UUID `cql:"team_id" json:"team_id"`
	UpdatedAt      time.Time  `cql:"updated_at" json:"updated_at"`
}

var (
	githubrepoColumns = []string{"created_at", "full_name", "github_id", "id", "installation_id", "name", "team_id", "updated_at"}

	githubrepoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_repos",
			Columns: githubrepoColumns,
		},
	}

	githubrepoTable = itable.New(*githubrepoMeta.M)
)

func (githubrepo *Repo) GetTable() itable.ITable {
	return githubrepoTable
}

// SetupAction defines model for SetupAction.
type SetupAction string

// WorkflowResponse workflow status & run id
type WorkflowResponse struct {
	RunID string `json:"run_id"`

	// Status the workflow status
	Status WorkflowStatus `json:"status"`
}

// WorkflowStatus the workflow status
type WorkflowStatus string

// GithubArtifactReadyJSONRequestBody defines body for GithubArtifactReady for application/json ContentType.
type GithubArtifactReadyJSONRequestBody = ArtifactReadyRequest

// GithubCompleteInstallationJSONRequestBody defines body for GithubCompleteInstallation for application/json ContentType.
type GithubCompleteInstallationJSONRequestBody = CompleteInstallationRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// GitHub release artifact ready
	// (POST /providers/github/artifact-ready)
	GithubArtifactReady(ctx echo.Context) error

	// Complete GitHub App installation
	// (POST /providers/github/complete-installation)
	GithubCompleteInstallation(ctx echo.Context) error

	// Get GitHub installations
	// (GET /providers/github/installations)
	GithubGetInstallations(ctx echo.Context) error

	// Get GitHub repositories
	// (GET /providers/github/repos)
	GithubGetRepos(ctx echo.Context) error

	// Webhook reciever for github
	// (POST /providers/github/webhook)
	GithubWebhook(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GithubArtifactReady converts echo context to params.

func (w *ServerInterfaceWrapper) GithubArtifactReady(ctx echo.Context) error {
	var err error

	ctx.Set(APIKeyAuthScopes, []string{})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubArtifactReady
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubCompleteInstallation converts echo context to params.

func (w *ServerInterfaceWrapper) GithubCompleteInstallation(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubCompleteInstallation
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubGetInstallations converts echo context to params.

func (w *ServerInterfaceWrapper) GithubGetInstallations(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubGetInstallations
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubGetRepos converts echo context to params.

func (w *ServerInterfaceWrapper) GithubGetRepos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubGetRepos
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubWebhook converts echo context to params.

func (w *ServerInterfaceWrapper) GithubWebhook(ctx echo.Context) error {
	var err error

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubWebhook
	err = handler(ctx)

	return err
}

// EchoRouter is an interface that wraps the methods of echo.Echo & echo.Group to provide a common interface
// for registering routes.
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/providers/github/artifact-ready", wrapper.GithubArtifactReady)
	router.POST(baseURL+"/providers/github/complete-installation", wrapper.GithubCompleteInstallation)
	router.GET(baseURL+"/providers/github/installations", wrapper.GithubGetInstallations)
	router.GET(baseURL+"/providers/github/repos", wrapper.GithubGetRepos)
	router.POST(baseURL+"/providers/github/webhook", wrapper.GithubWebhook)

}
