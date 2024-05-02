// Package slack provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/breuHQ/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen/v2.
//
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	"github.com/scylladb/gocqlx/v2/table"
	"go.breu.io/quantm/internal/shared"
	externalRef0 "go.breu.io/quantm/internal/shared"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Slack defines model for Slack.
type Slack struct {
	ChannelID         string     `cql:"channel_id" json:"channel_id"`
	ChannelName       string     `cql:"channel_name" json:"channel_name"`
	CreatedAt         time.Time  `cql:"created_at" json:"created_at"`
	ID                gocql.UUID `cql:"id" json:"id"`
	TeamID            gocql.UUID `cql:"team_id" json:"team_id"`
	UpdatedAt         time.Time  `cql:"updated_at" json:"updated_at"`
	WorkspaceBotToken string     `cql:"workspace_bot_token" json:"workspace_bot_token"`
	WorkspaceID       string     `cql:"workspace_id" json:"workspace_id"`
	WorkspaceName     string     `cql:"workspace_name" json:"workspace_name"`
}

var (
	slackMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "slack",
			Columns: []string{"channel_id", "channel_name", "created_at", "id", "team_id", "updated_at", "workspace_bot_token", "workspace_id", "workspace_name"},
			PartKey: []string{"team_id"},
		},
	}

	slackTable = itable.New(*slackMeta.M)
)

func (slack *Slack) GetTable() itable.ITable {
	return slackTable
}

// SlackOauthParams defines parameters for SlackOauth.
type SlackOauthParams struct {
	Code string `form:"code" json:"code"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SlackOauth request
	SlackOauth(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SlackOauth(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSlackOauthRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSlackOauthRequest generates requests for SlackOauth
func NewSlackOauthRequest(server string, params *SlackOauthParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/slack")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "code", runtime.ParamLocationQuery, params.Code); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SlackOauthWithResponse request
	SlackOauthWithResponse(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*SlackOauthResponse, error)
}

type SlackOauthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Slack
	JSON400      *externalRef0.BadRequest
	JSON500      *externalRef0.InternalServerError
}

// Status returns HTTPResponse.Status
func (r SlackOauthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SlackOauthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SlackOauthWithResponse request returning *SlackOauthResponse
func (c *ClientWithResponses) SlackOauthWithResponse(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*SlackOauthResponse, error) {
	rsp, err := c.SlackOauth(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSlackOauthResponse(rsp)
}

// ParseSlackOauthResponse parses an HTTP response from a SlackOauthWithResponse call
func ParseSlackOauthResponse(rsp *http.Response) (*SlackOauthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SlackOauthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Slack
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef0.BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef0.InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Slack Oauth and save the info
	// (GET /v1/slack)
	SlackOauth(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SlackOauth converts echo context to params.

func (w *ServerInterfaceWrapper) SlackOauth(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params SlackOauthParams
	// ------------- Required query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, true, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter code: %s", err))
	}

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.SlackOauth
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

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

	router.GET(baseURL+"/v1/slack", wrapper.SlackOauth)

}
