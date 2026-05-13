-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password,
    batch,
    status,
    role
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING 
    id,
    name,
    email,
    batch,
    status,
    role,
    created_at,
    updated_at;


-- name: GetUser :one
SELECT 
    id,
    name,
    email,
    batch,
    status,
    role,
    created_at,
    updated_at
FROM users
WHERE id = $1;


-- name: GetUserByEmail :one
SELECT 
    id,
    name,
    email,
    password,
    batch,
    status,
    role,
    created_at,
    updated_at
FROM users
WHERE email = $1;


-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    email = $3,
    password = $4,
    batch = $5,
    status = $6,
    role = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    batch,
    status,
    role,
    created_at,
    updated_at;


-- name: UpdateUserProfile :one
UPDATE users
SET
    name = $2,
    email = $3,
    batch = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    batch,
    status,
    role,
    created_at,
    updated_at;
