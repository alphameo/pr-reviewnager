-- name: CreateTeam :exec
INSERT INTO team (id, name)
VALUES ($1, $2);

-- name: GetTeams :many
SELECT
    id,
    name
FROM team;

-- name: GetTeam :one
SELECT
    id,
    name
FROM team
WHERE id = $1;

-- name: GetTeamByName :one
SELECT
    id,
    name
FROM team
WHERE name = $1;

-- name: UpdateTeam :exec
UPDATE team
SET name = $2
WHERE id = $1;

-- name: DeleteTeam :exec
DELETE FROM team
WHERE id = $1;
