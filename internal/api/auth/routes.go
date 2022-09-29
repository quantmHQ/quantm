// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package auth

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
)

// CreateRoutes is for creating auth related routes.
func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	r := &Routes{}
	g.POST("/register", r.register)
	g.POST("/login", r.login)

	aks := g.Group("/api-keys", Middleware)

	tr := &TeamAPIKeyRoutes{}
	aks.POST("/team", tr.create)

	ur := &UserAPIKeyRoutes{}
	aks.POST("/user", ur.create)
}

type (
	Routes           struct{}
	TeamAPIKeyRoutes struct{}
	UserAPIKeyRoutes struct{}
)

// register is a handler for /auth/register endpoint.
func (routes *Routes) register(ctx echo.Context) error {
	request := &RegistrationRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return err
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return err
	}

	// Validating team
	team := &entities.Team{Name: request.TeamName}
	if err := ctx.Validate(team); err != nil {
		return err
	}

	user := &entities.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}
	if err := ctx.Validate(user); err != nil {
		return err
	}

	if err := db.Save(team); err != nil {
		return err
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &RegistrationResponse{Team: team, User: user})
}

// login is a handler for /auth/login endpoint.
func (routes *Routes) login(ctx echo.Context) error {
	request := &LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &entities.User{}

	if err := db.Get(user, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if user.VerifyPassword(request.Password) {
		access, _ := GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: access, RefreshToken: refresh})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
}

// create a new API Key for team.
func (routes *TeamAPIKeyRoutes) create(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	id, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	guard := &entities.Guard{}
	key, err := guard.NewForTeam(id)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: key})
}

// create a new API Key for user.
func (routes *UserAPIKeyRoutes) create(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	id, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	guard := &entities.Guard{}
	key, err := guard.NewForUser(request.Name, id)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: key})
}
