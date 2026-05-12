-- name: CreateUser :one
INSERT INTO
    users (
        id,
        Name,
        Email,
        password,
        createdAt,
        updatedAt
    )
VALUES ($1, $2, $3, $4,$5,$6) rETURNING id,
    name,
    createdAt,
    updatedAt;
;

-- name: GetUser :one
SELECT name, email FROM users WHERE id = $1;
-- name: GetUserByEmail :one
SELECT id, name password FROM users WHERE email = $1;
-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    email = $3,
    password = $4,
    updatedAt = $5
WHERE id = $1
RETURNING id, name, email, createdAt, updatedAt;
-- name: UpdateUserProfile :one
UPDATE users
SET
    name = $2,
    email = $3,
    updatedAt = $4
WHERE id = $1
RETURNING id, name, email, createdAt, updatedAt;

