-- name: CreatePullRequestReviewer :exec
INSERT INTO pull_request_reviewer (pull_request_id, reviewer_id)
VALUES ($1, $2);

-- name: GetPullRequestReviewerReviewerIDs :many
SELECT reviewer_id FROM pull_request_reviewer WHERE pull_request_id = $1;

-- name: GetPullRequestsByReviewer :many
SELECT pr.id, pr.title, pr.author_id, pr.status, pr.merged_at
FROM pull_request pr
JOIN pull_request_reviewer prr ON pr.id = prr.pull_request_id
WHERE prr.reviewer_id = $1;

-- name: DeletePullRequestReviewer :exec
DELETE FROM pull_request_reviewer
WHERE pull_request_id = $1 AND reviewer_id = $2;

-- name: DeletePullRequestReviewersByReviewerID :exec
DELETE FROM pull_request_reviewer
WHERE reviewer_id = $1;

-- name: DeletePullRequestReviewersByPRID :exec
DELETE FROM pull_request_reviewer
WHERE pull_request_id = $1;
