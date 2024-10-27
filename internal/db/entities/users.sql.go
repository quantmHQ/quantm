// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, org_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, org_id, email, first_name, last_name, password, is_active, is_verified
`

type CreateUserParams struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	OrgID     uuid.UUID `json:"org_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.OrgID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}

const getFullUserByEmail = `-- name: GetFullUserByEmail :one
SELECT
  u.id,
  u.created_at,
  u.updated_at,
  u.first_name,
  u.last_name,
  u.email,
  u.org_id,
  ARRAY_AGG(
    DISTINCT ROW(
      team.id,
      team.created_at,
      team.updated_at,
      team.org_id,
      team.name,
      team.slug
    )
  ) AS teams,
  ARRAY_AGG(
    DISTINCT ROW(
      account.id,
      account.created_at,
      account.updated_at,
      account.user_id,
      account.provider,
      account.provider_account_id,
      account.expires_at,
      account.type
      )
  ) as accounts,
  ROW(org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug) AS org
FROM users AS u
LEFT JOIN team_users AS team_user
  ON u.id = team_user.user_id
LEFT JOIN teams AS team
  ON team_user.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON u.id = account.user_id
LEFT JOIN orgs AS org
  ON u.org_id = org.id
WHERE
  u.email = LOWER($1)
GROUP BY u.id, org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug
`

type GetFullUserByEmailRow struct {
	ID        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	OrgID     uuid.UUID   `json:"org_id"`
	Teams     interface{} `json:"teams"`
	Accounts  interface{} `json:"accounts"`
	Org       interface{} `json:"org"`
}

func (q *Queries) GetFullUserByEmail(ctx context.Context, lower string) (GetFullUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getFullUserByEmail, lower)
	var i GetFullUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.OrgID,
		&i.Teams,
		&i.Accounts,
		&i.Org,
	)
	return i, err
}

const getFullUserByID = `-- name: GetFullUserByID :one
SELECT
  u.id,
  u.created_at,
  u.updated_at,
  u.first_name,
  u.last_name,
  u.email,
  u.org_id,
  ARRAY_AGG(
    DISTINCT ROW(
      team.id,
      team.created_at,
      team.updated_at,
      team.org_id,
      team.name,
      team.slug
    )
  ) AS teams,
  ARRAY_AGG(
    DISTINCT ROW(
      account.id,
      account.created_at,
      account.updated_at,
      account.user_id,
      account.provider,
      account.provider_account_id,
      account.expires_at,
      account.type
      )
  ) as accounts,
  ROW(org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug) AS org
FROM users AS u
LEFT JOIN team_users AS team_user
  ON u.id = team_user.user_id
LEFT JOIN teams AS team
  ON team_user.team_id = team.id
LEFT JOIN oauth_accounts AS account
  ON u.id = account.user_id
LEFT JOIN orgs AS org
  ON u.org_id = org.id
WHERE
  u.id = $1
GROUP BY u.id, org.id, org.created_at, org.updated_at, org.name, org.domain, org.slug
`

type GetFullUserByIDRow struct {
	ID        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	OrgID     uuid.UUID   `json:"org_id"`
	Teams     interface{} `json:"teams"`
	Accounts  interface{} `json:"accounts"`
	Org       interface{} `json:"org"`
}

func (q *Queries) GetFullUserByID(ctx context.Context, id uuid.UUID) (GetFullUserByIDRow, error) {
	row := q.db.QueryRow(ctx, getFullUserByID, id)
	var i GetFullUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.OrgID,
		&i.Teams,
		&i.Accounts,
		&i.Org,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, org_id, email, first_name, last_name, password, is_active, is_verified
FROM users
WHERE email = LOWER($1)
`

func (q *Queries) GetUserByEmail(ctx context.Context, lower string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, lower)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}

const getUserByEmailFull = `-- name: GetUserByEmailFull :one
SELECT
  u.id, u.created_at, u.updated_at, u.first_name, u.last_name, u.email, u.org_id,
  array_agg(t.*) AS teams,
  array_agg(oa.*) AS oauth_accounts,
  array_agg(o.*) AS orgs
FROM users AS u
LEFT JOIN team_users AS tu
  ON u.id = tu.user_id
LEFT JOIN teams AS t
  ON tu.team_id = t.id
LEFT JOIN oauth_accounts AS oa
  ON u.id = oa.user_id
LEFT JOIN orgs AS o
  ON u.org_id = o.id
WHERE
  u.email = LOWER($1)
GROUP BY
  u.id
`

type GetUserByEmailFullRow struct {
	ID            uuid.UUID   `json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	FirstName     string      `json:"first_name"`
	LastName      string      `json:"last_name"`
	Email         string      `json:"email"`
	OrgID         uuid.UUID   `json:"org_id"`
	Teams         interface{} `json:"teams"`
	OauthAccounts interface{} `json:"oauth_accounts"`
	Orgs          interface{} `json:"orgs"`
}

func (q *Queries) GetUserByEmailFull(ctx context.Context, lower string) (GetUserByEmailFullRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmailFull, lower)
	var i GetUserByEmailFullRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.OrgID,
		&i.Teams,
		&i.OauthAccounts,
		&i.Orgs,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, created_at, updated_at, org_id, email, first_name, last_name, password, is_active, is_verified
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}

const getUserByProviderAccount = `-- name: GetUserByProviderAccount :one
SELECT
  u.id, u.created_at, u.updated_at, u.org_id, u.email, u.first_name, u.last_name, u.password, u.is_active, u.is_verified
FROM users as u
WHERE u.id IN (
  SELECT user_id
  FROM oauth_accounts
  WHERE provider = $1 AND provider_account_id = $2
)
`

type GetUserByProviderAccountParams struct {
	Provider          string `json:"provider"`
	ProviderAccountID string `json:"provider_account_id"`
}

func (q *Queries) GetUserByProviderAccount(ctx context.Context, arg GetUserByProviderAccountParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByProviderAccount, arg.Provider, arg.ProviderAccountID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = LOWER($4), org_id = $5
WHERE id = $1
RETURNING id, created_at, updated_at, org_id, email, first_name, last_name, password, is_active, is_verified
`

type UpdateUserParams struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Lower     string    `json:"lower"`
	OrgID     uuid.UUID `json:"org_id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Lower,
		arg.OrgID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}
