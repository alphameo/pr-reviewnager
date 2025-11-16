-- name: CreateUser :exec
INSERT INTO "user" (id, name, active)
VALUES ($1, $2, $3);

-- name: GetUsers :many
SELECT id, name, active FROM "user";

-- name: GetUser :one
SELECT id, name, active FROM "user" WHERE id = $1;

-- name: GetUserByName :one
SELECT id, name, active FROM "user" WHERE name = $1;

-- name: UpdateUser :exec
UPDATE "user"
SET name = $2, active = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

-- name: UpsetUser :exec
INSERT INTO "user" (id, name, active) 
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET 
    name = EXCLUDED.name,
    acive = EXCLUDED.active;
