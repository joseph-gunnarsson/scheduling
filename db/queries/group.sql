-- name: CreateGroup :one
INSERT INTO groups (name, description, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, name, description, owner_id, created_at, updated_at;

-- name: UpdateGroup :one
UPDATE groups
SET name = $2,
    description = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: GetGroupByID :one
SELECT *
FROM groups
WHERE id = $1;

-- name: GetGroupsByOwner :many
SELECT * FROM groups
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: PatchGroup :one
UPDATE groups
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE($2, description),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: AddUserToGroup :one
INSERT INTO user_groups (user_id, group_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteUserFromGroup :exec
DELETE FROM user_groups
WHERE user_id = $1 AND group_id = $2;

-- name: GetUserGroups :many
SELECT 
    g.id AS group_id,
    g.name AS group_name,
    g.description AS group_description,
    g.owner_id AS group_owner_id,
    g.created_at AS group_created_at,
    g.updated_at AS group_updated_at,
    ug.joined_at AS user_joined_at
FROM user_groups ug
JOIN groups g ON ug.group_id = g.id
WHERE ug.user_id = $1;

-- name: GetGroupMembers :many
SELECT 
    u.id AS user_id,
    u.username,
    u.email,
    u.first_name,
    u.last_name,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    ug.joined_at
FROM user_groups ug
JOIN users u ON ug.user_id = u.id
WHERE ug.group_id = $1;
