-- name: CreateTeamUser :exec
INSERT INTO team_user (team_id, user_id)
VALUES ($1, $2);

-- name: RemoveUserFromTeam :exec
DELETE FROM team_user
WHERE team_id = $1 AND user_id = $2;

-- name: DeleteTeamUsersByTeamID :exec
DELETE FROM team_user
WHERE team_id = $1;

-- name: GetUsersInTeam :many
SELECT u.id, u.name, u.active
FROM "user" u
JOIN team_user tu ON u.id = tu.user_id
WHERE tu.team_id = $1;

-- name: GetUserIDsInTeam :many
SELECT user_id
FROM team_user
WHERE team_id = $1;

-- name: GetTeamIDForUser :one
SELECT t.id
FROM team t
JOIN team_user tu ON t.id = tu.team_id
WHERE tu.user_id = $1;

-- name: GetTeamForUser :one
SELECT t.id, t.name
FROM team t
JOIN team_user tu ON t.id = tu.team_id
WHERE tu.user_id = $1;

-- name: GetActiveUsersInTeam :many
SELECT u.id, u.name, u.active
FROM "user" u
JOIN team_user tu ON u.id = tu.user_id
WHERE tu.team_id = $1 AND u.active = true;

-- name: GetTeamsWithUsers :many
SELECT 
    t.id as team_id,
    t.name as team_name,
    u.id as user_id,
    u.name as user_name,
    u.active as user_active
FROM team t
LEFT JOIN team_user tu ON t.id = tu.team_id
LEFT JOIN "user" u ON tu.user_id = u.id
ORDER BY t.id;
