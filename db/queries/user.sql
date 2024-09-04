-- Create a new user
-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, first_name, last_name, created_at, updated_at;

-- Get user by ID
-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- Get user by username
-- name: GetUserByUsername :one
SELECT id, username, email, first_name, last_name
FROM users
WHERE username = $1;

-- Update user password
-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Login user by username or email and password
-- name: LoginUser :one
SELECT id, username, email, first_name, last_name, password_hash
FROM users
WHERE (username = $1 OR email = $1);

