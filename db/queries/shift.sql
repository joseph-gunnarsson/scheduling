-- Create a new shift
-- name: CreateShift :one
INSERT INTO shifts (user_id, group_id, name, start_time, end_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, user_id, group_id, name, start_time, end_time, created_at, updated_at;

-- Update a shift
-- name: UpdateShift :exec
UPDATE shifts
SET name = $1,
    start_time = $2,
    end_time = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $4 AND user_id = $5 AND group_id = $6;

-- Delete a shift by ID
-- name: DeleteShift :exec
DELETE FROM shifts
WHERE id = $1 AND user_id = $2 AND group_id = $3;

-- Get shift by ID
-- name: GetShiftByID :one
SELECT id, user_id, group_id, name, start_time, end_time, created_at, updated_at
FROM shifts
WHERE id = $1;

-- List all shifts for a specific user in a group
-- name: ListShiftsByUserAndGroup :many
SELECT id, user_id, group_id, name, start_time, end_time, created_at, updated_at
FROM shifts
WHERE user_id = $1 AND group_id = $2
ORDER BY start_time ASC;

-- List all shifts in a specific group
-- name: ListShiftsByGroup :many
SELECT id, user_id, group_id, name, start_time, end_time, created_at, updated_at
FROM shifts
WHERE group_id = $1
ORDER BY start_time ASC;

-- List all shifts for a specific user across all groups
-- name: ListShiftsByUser :many
SELECT id, user_id, group_id, name, start_time, end_time, created_at, updated_at
FROM shifts
WHERE user_id = $1
ORDER BY start_time ASC;

-- List all shifts
-- name: ListAllShifts :many
SELECT id, user_id, group_id, name, start_time, end_time, created_at, updated_at
FROM shifts
ORDER BY start_time ASC;

-- Get all shifts by group ID including user names
-- name: ListShiftsByGroupWithNames :many
SELECT
    shifts.id AS shift_id,
    shifts.name AS shift_name,
    shifts.start_time,
    shifts.end_time,
    shifts.created_at AS shift_created_at,
    shifts.updated_at AS shift_updated_at,
    users.id AS user_id,
    users.first_name AS user_first_name,
    users.last_name AS user_last_name
FROM
    shifts
JOIN
    users ON shifts.user_id = users.id
WHERE
    shifts.group_id = $1
ORDER BY
    shifts.start_time ASC;



