// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: oauth_accounts.sql

package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createOAuthAccount = `-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_account_id, expires_at, type)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, provider, provider_account_id, expires_at, type
`

type CreateOAuthAccountParams struct {
	UserID            uuid.UUID `json:"user_id"`
	Provider          string    `json:"provider"`
	ProviderAccountID string    `json:"provider_account_id"`
	ExpiresAt         time.Time `json:"expires_at"`
	Type              string    `json:"type"`
}

func (q *Queries) CreateOAuthAccount(ctx context.Context, arg CreateOAuthAccountParams) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, createOAuthAccount,
		arg.UserID,
		arg.Provider,
		arg.ProviderAccountID,
		arg.ExpiresAt,
		arg.Type,
	)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Provider,
		&i.ProviderAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountByID = `-- name: GetOAuthAccountByID :one
SELECT id, created_at, updated_at, user_id, provider, provider_account_id, expires_at, type
FROM oauth_accounts
WHERE id = $1
`

func (q *Queries) GetOAuthAccountByID(ctx context.Context, id uuid.UUID) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, getOAuthAccountByID, id)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Provider,
		&i.ProviderAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountByProviderAccountID = `-- name: GetOAuthAccountByProviderAccountID :one
SELECT id, created_at, updated_at, user_id, provider, provider_account_id, expires_at, type
FROM oauth_accounts
WHERE provider_account_id = $1 and provider = $2
`

type GetOAuthAccountByProviderAccountIDParams struct {
	ProviderAccountID string `json:"provider_account_id"`
	Provider          string `json:"provider"`
}

func (q *Queries) GetOAuthAccountByProviderAccountID(ctx context.Context, arg GetOAuthAccountByProviderAccountIDParams) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, getOAuthAccountByProviderAccountID, arg.ProviderAccountID, arg.Provider)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Provider,
		&i.ProviderAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountsByUserID = `-- name: GetOAuthAccountsByUserID :many
SELECT id, created_at, updated_at, user_id, provider, provider_account_id, expires_at, type
FROM oauth_accounts
WHERE user_id = $1
`

func (q *Queries) GetOAuthAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]OauthAccount, error) {
	rows, err := q.db.Query(ctx, getOAuthAccountsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OauthAccount
	for rows.Next() {
		var i OauthAccount
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Provider,
			&i.ProviderAccountID,
			&i.ExpiresAt,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
