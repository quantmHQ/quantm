// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: github_installation.sql

package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createInstallation = `-- name: CreateInstallation :one
INSERT INTO github_installations (org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status
`

type CreateInstallationParams struct {
	OrgID               uuid.UUID   `json:"org_id"`
	InstallationID      int64       `json:"installation_id"`
	InstallationLogin   string      `json:"installation_login"`
	InstallationLoginID int64       `json:"installation_login_id"`
	InstallationType    pgtype.Text `json:"installation_type"`
	SenderID            int64       `json:"sender_id"`
	SenderLogin         string      `json:"sender_login"`
	Status              pgtype.Text `json:"status"`
}

func (q *Queries) CreateInstallation(ctx context.Context, arg CreateInstallationParams) (GithubInstallation, error) {
	row := q.db.QueryRow(ctx, createInstallation,
		arg.OrgID,
		arg.InstallationID,
		arg.InstallationLogin,
		arg.InstallationLoginID,
		arg.InstallationType,
		arg.SenderID,
		arg.SenderLogin,
		arg.Status,
	)
	var i GithubInstallation
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.InstallationID,
		&i.InstallationLogin,
		&i.InstallationLoginID,
		&i.InstallationType,
		&i.SenderID,
		&i.SenderLogin,
		&i.Status,
	)
	return i, err
}

const deleteInstallation = `-- name: DeleteInstallation :one
DELETE FROM github_installations
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteInstallation(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, deleteInstallation, id)
	err := row.Scan(&id)
	return id, err
}

const getInstallation = `-- name: GetInstallation :one
SELECT id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status
FROM github_installations
WHERE id = $1
`

func (q *Queries) GetInstallation(ctx context.Context, id uuid.UUID) (GithubInstallation, error) {
	row := q.db.QueryRow(ctx, getInstallation, id)
	var i GithubInstallation
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.InstallationID,
		&i.InstallationLogin,
		&i.InstallationLoginID,
		&i.InstallationType,
		&i.SenderID,
		&i.SenderLogin,
		&i.Status,
	)
	return i, err
}

const getInstallationByInstallationIDAndInstallationLogin = `-- name: GetInstallationByInstallationIDAndInstallationLogin :one
SELECT id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status
FROM github_installations
WHERE installation_id = $1 AND installation_login = $2
`

type GetInstallationByInstallationIDAndInstallationLoginParams struct {
	InstallationID    int64  `json:"installation_id"`
	InstallationLogin string `json:"installation_login"`
}

func (q *Queries) GetInstallationByInstallationIDAndInstallationLogin(ctx context.Context, arg GetInstallationByInstallationIDAndInstallationLoginParams) (GithubInstallation, error) {
	row := q.db.QueryRow(ctx, getInstallationByInstallationIDAndInstallationLogin, arg.InstallationID, arg.InstallationLogin)
	var i GithubInstallation
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.InstallationID,
		&i.InstallationLogin,
		&i.InstallationLoginID,
		&i.InstallationType,
		&i.SenderID,
		&i.SenderLogin,
		&i.Status,
	)
	return i, err
}

const updateInstallation = `-- name: UpdateInstallation :one
UPDATE github_installations
SET 
    org_id = $2,
    installation_id = $3,
    installation_login = $4,
    installation_login_id = $5,
    installation_type = $6,
    sender_id = $7,
    sender_login = $8,
    status = $9
WHERE id = $1
RETURNING id, created_at, updated_at, org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login, status
`

type UpdateInstallationParams struct {
	ID                  uuid.UUID   `json:"id"`
	OrgID               uuid.UUID   `json:"org_id"`
	InstallationID      int64       `json:"installation_id"`
	InstallationLogin   string      `json:"installation_login"`
	InstallationLoginID int64       `json:"installation_login_id"`
	InstallationType    pgtype.Text `json:"installation_type"`
	SenderID            int64       `json:"sender_id"`
	SenderLogin         string      `json:"sender_login"`
	Status              pgtype.Text `json:"status"`
}

func (q *Queries) UpdateInstallation(ctx context.Context, arg UpdateInstallationParams) (GithubInstallation, error) {
	row := q.db.QueryRow(ctx, updateInstallation,
		arg.ID,
		arg.OrgID,
		arg.InstallationID,
		arg.InstallationLogin,
		arg.InstallationLoginID,
		arg.InstallationType,
		arg.SenderID,
		arg.SenderLogin,
		arg.Status,
	)
	var i GithubInstallation
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.InstallationID,
		&i.InstallationLogin,
		&i.InstallationLoginID,
		&i.InstallationType,
		&i.SenderID,
		&i.SenderLogin,
		&i.Status,
	)
	return i, err
}
