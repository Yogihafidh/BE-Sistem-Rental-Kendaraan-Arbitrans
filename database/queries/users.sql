-- name: GetRoles :many
SELECT * FROM roles ORDER BY name ASC;

-- name: CreateRole :one
INSERT INTO roles(name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1
LIMIT 1;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserWithRoleByEmail :one
SELECT
    u.id,
    u.role_id,
    u.full_name,
    u.email,
    u.phone,
    u.password_hash,
    u.is_active,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    r.name AS role_name
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.email = $1 AND u.deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    role_id, full_name, email, phone, password_hash
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET role_id = $2,
    full_name = $3,
    email = $4,
    phone = $5,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ActivateUser :one
UPDATE users
SET is_active = true, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeactivateUser :one
UPDATE users
SET is_active = false, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now(), updated_at = now(), is_active = false
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = now(), updated_at = now()
WHERE id = $1;
