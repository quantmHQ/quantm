-- name: CreateGithubInstallation :one
INSERT INTO github_installations (
  org_id, installation_id, installation_login, installation_login_id, installation_type, sender_id, sender_login
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetGithubInstallation :one
SELECT *
FROM github_installations
WHERE id = $1;

-- name: GetGithubInstallationByInstallationID :one
SELECT *
FROM github_installations
WHERE installation_id = $1;

-- name: GetGithubInstallationByInstallationIDAndInstallationLogin :one
SELECT *
FROM github_installations
WHERE installation_id = $1 AND installation_login = $2;

-- name: UpdateGithubInstallation :one
UPDATE github_installations
SET
    org_id = $2,
    installation_id = $3,
    installation_login = $4,
    installation_login_id = $5,
    installation_type = $6,
    sender_id = $7,
    sender_login = $8,
    is_active = $9
WHERE id = $1
RETURNING *;
