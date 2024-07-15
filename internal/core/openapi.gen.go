// Package core provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/breuHQ/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen/v2.
//
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	"github.com/scylladb/gocqlx/v2/table"
	"go.breu.io/quantm/internal/shared"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidMessageProvider = errors.New("invalid MessageProvider value")
	ErrInvalidRepoProvider    = errors.New("invalid RepoProvider value")
)

type (
	MessageProviderMapType map[string]MessageProvider // MessageProviderMapType is a quick lookup map for MessageProvider.
)

// Defines values for MessageProvider.
const (
	MessageProviderNone  MessageProvider = "none"
	MessageProviderSlack MessageProvider = "slack"
)

// MessageProviderMap returns all known values for MessageProvider.
var (
	MessageProviderMap = MessageProviderMapType{
		MessageProviderNone.String():  MessageProviderNone,
		MessageProviderSlack.String(): MessageProviderSlack,
	}
)

/*
 * Helper methods for MessageProvider for easy marshalling and unmarshalling.
 */
func (v MessageProvider) String() string               { return string(v) }
func (v MessageProvider) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *MessageProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := MessageProviderMap[s]
	if !ok {
		return ErrInvalidMessageProvider
	}

	*v = val

	return nil
}

type (
	RepoProviderMapType map[string]RepoProvider // RepoProviderMapType is a quick lookup map for RepoProvider.
)

// Defines values for RepoProvider.
const (
	RepoProviderBitbucket RepoProvider = "bitbucket"
	RepoProviderGithub    RepoProvider = "github"
	RepoProviderGitlab    RepoProvider = "gitlab"
)

// RepoProviderMap returns all known values for RepoProvider.
var (
	RepoProviderMap = RepoProviderMapType{
		RepoProviderBitbucket.String(): RepoProviderBitbucket,
		RepoProviderGithub.String():    RepoProviderGithub,
		RepoProviderGitlab.String():    RepoProviderGitlab,
	}
)

/*
 * Helper methods for RepoProvider for easy marshalling and unmarshalling.
 */
func (v RepoProvider) String() string               { return string(v) }
func (v RepoProvider) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *RepoProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := RepoProviderMap[s]
	if !ok {
		return ErrInvalidRepoProvider
	}

	*v = val

	return nil
}

// BranchChanges Branch changes lines exceed.
type BranchChanges struct {
	// Additions Number of additions
	Additions int `json:"additions"`

	// Changes Total number of changes
	Changes int `json:"changes"`

	// CompareUrl Compare url
	CompareUrl string `json:"compare_url"`

	// Deletions Number of deletions
	Deletions int `json:"deletions"`

	// FileCount Number of changed files
	FileCount int `json:"file_count"`

	// Files List of changed files
	Files []string `json:"files"`

	// RepoUrl Repo url
	RepoUrl string `json:"repo_url"`
}

// LatestCommit get the latest commit information.
type LatestCommit struct {
	Branch    string `json:"branch"`
	CommitUrl string `json:"commit_url"`
	RepoName  string `json:"repo_name"`
	RepoUrl   string `json:"repo_url"`
	SHA       string `json:"sha"`
}

// MessageProvider defines model for MessageProvider.
type MessageProvider string

// MessageProviderData defines model for MessageProviderData.
type MessageProviderData struct {
	Slack *MessageProviderSlackData `json:"slack,omitempty"`
}

// MessageProviderSlackData defines model for MessageProviderSlackData.
type MessageProviderSlackData struct {
	BotToken      string `json:"bot_token"`
	ChannelID     string `json:"channel_id"`
	ChannelName   string `json:"channel_name"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

// Repo defines model for Repo.
type Repo struct {
	CreatedAt time.Time `cql:"created_at" json:"created_at"`

	// CtrlId references the id field of the repos tables against the provider. For us, this means, that it will be the id field for
	//   - github_repos
	//   - gitlab_repos
	// etc.
	CtrlID              gocql.UUID          `cql:"ctrl_id" json:"ctrl_id"`
	DefaultBranch       string              `cql:"default_branch" json:"default_branch"`
	ID                  gocql.UUID          `cql:"id" json:"id"`
	IsMonorepo          bool                `cql:"is_monorepo" json:"is_monorepo"`
	MessageProvider     MessageProvider     `cql:"message_provider" json:"message_provider"`
	MessageProviderData MessageProviderData `cql:"message_provider_data" json:"message_provider_data"`
	Name                string              `cql:"name" json:"name"`
	Provider            RepoProvider        `cql:"provider" json:"provider"`
	ProviderID          string              `cql:"provider_id" json:"provider_id"`
	StaleDuration       shared.Duration     `cql:"stale_duration" json:"stale_duration"`
	TeamID              gocql.UUID          `cql:"team_id" json:"team_id"`
	Threshold           shared.Int64        `cql:"threshold" json:"threshold"`
	UpdatedAt           time.Time           `cql:"updated_at" json:"updated_at"`
}

var (
	repoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "repos",
			Columns: []string{"created_at", "ctrl_id", "default_branch", "id", "is_monorepo", "message_provider", "message_provider_data", "name", "provider", "provider_id", "stale_duration", "team_id", "threshold", "updated_at"},
			PartKey: []string{"id", "team_id"},
		},
	}

	repoTable = itable.New(*repoMeta.M)
)

func (repo *Repo) GetTable() itable.ITable {
	return repoTable
}

// RepoCreateRequest defines model for RepoCreateRequest.
type RepoCreateRequest struct {
	CtrlID              gocql.UUID          `json:"ctrl_id"`
	IsMonorepo          bool                `json:"is_monorepo"`
	MessageProvider     MessageProvider     `json:"message_provider"`
	MessageProviderData MessageProviderData `json:"message_provider_data"`
	Provider            RepoProvider        `json:"provider"`
	StaleDuration       shared.Duration     `json:"stale_duration"`
	Threshold           shared.Int64        `json:"threshold"`
}

// RepoListResponse defines model for RepoListResponse.
type RepoListResponse = []Repo

// RepoProvider defines model for RepoProvider.
type RepoProvider string

// RepoProviderData defines model for RepoProviderData.
type RepoProviderData struct {
	DefaultBranch string `json:"default_branch"`
	Name          string `json:"name"`
}

// CreateRepoJSONRequestBody defines body for CreateRepo for application/json ContentType.
type CreateRepoJSONRequestBody = RepoCreateRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List Repos
	// (GET /core/repos)
	ListRepos(ctx echo.Context) error

	// Create repo
	// (POST /core/repos)
	CreateRepo(ctx echo.Context) error

	// Get repo
	// (GET /core/repos/{id})
	GetRepo(ctx echo.Context, id string) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(ctx echo.Context, handler echo.HandlerFunc) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListRepos converts echo context to params.

func (w *ServerInterfaceWrapper) ListRepos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	handler := func(ctx echo.Context) error {
		return w.Handler.ListRepos(ctx)
	}
	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SecureHandler(ctx, handler)

	return err
}

// CreateRepo converts echo context to params.

func (w *ServerInterfaceWrapper) CreateRepo(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	handler := func(ctx echo.Context) error {
		return w.Handler.CreateRepo(ctx)
	}
	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SecureHandler(ctx, handler)

	return err
}

// GetRepo converts echo context to params.

func (w *ServerInterfaceWrapper) GetRepo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	handler := func(ctx echo.Context) error {
		return w.Handler.GetRepo(ctx, id)
	}
	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SecureHandler(ctx, handler)

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

	router.GET(baseURL+"/core/repos", wrapper.ListRepos)
	router.POST(baseURL+"/core/repos", wrapper.CreateRepo)
	router.GET(baseURL+"/core/repos/:id", wrapper.GetRepo)

}
