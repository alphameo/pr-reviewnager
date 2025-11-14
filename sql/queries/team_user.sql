-- name: AddUserToTeam :exec
INSERT INTO team_user (team_id, user_id)
VALUES ($1, $2);

-- name: RemoveUserFromTeam :exec
DELETE FROM team_user
WHERE team_id = $1 AND user_id = $2;

-- name: GetUsersInTeam :many
SELECT u.id, u.name, u.active
FROM "user" u
JOIN team_user tu ON u.id = tu.user_id
WHERE tu.team_id = $1;

-- name: GetTeamsForUser :many
SELECT t.id, t.name
FROM team t
JOIN team_user tu ON t.id = tu.team_id
WHERE tu.user_id = $1;
