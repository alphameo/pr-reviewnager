-- name: AddReviewerToPullRequest :exec
INSERT INTO pull_request_reviewer (pull_request_id, reviewer_id)
VALUES ($1, $2);

-- name: GetPullRequestReviewers :many
SELECT reviewer_id FROM pull_request_reviewer WHERE pull_request_id = $1;

-- name: RemoveReviewerFromPullRequest :exec
DELETE FROM pull_request_reviewer
WHERE pull_request_id = $1 AND reviewer_id = $2;
