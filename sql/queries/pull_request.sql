-- name: CreatePullRequest :exec
INSERT INTO pull_request (id, title, author_id, created_at, status, merged_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetPullRequests :many
SELECT id, title, author_id, created_at, status, merged_at FROM pull_request;

-- name: GetPullRequest :one
SELECT id, title, author_id, created_at, status, merged_at FROM pull_request WHERE id = $1;

-- name: UpdatePullRequest :exec
UPDATE pull_request
SET title = $2, author_id = $3, created_at = $4, status = $5, merged_at = $6
WHERE id = $1;

-- name: UpdatePullRequestStatus :exec
UPDATE pull_request
SET status = $2, merged_at = $3
WHERE id = $1;

-- name: DeletePullRequest :exec
DELETE FROM pull_request
WHERE id = $1;
