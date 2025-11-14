-- name: CreatePullRequest :one
INSERT INTO pull_request (id, title, author_id, status, merged_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, author_id, status, merged_at;

-- name: GetPullRequests :many
SELECT id, title, author_id, status, merged_at FROM pull_request;

-- name: GetPullRequest :one
SELECT id, title, author_id, status, merged_at FROM pull_request WHERE id = $1;

-- name: UpdatePullRequestStatus :exec
UPDATE pull_request
SET status = $2, merged_at = $3
WHERE id = $1;
