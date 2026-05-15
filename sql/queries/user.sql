-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password,
    batch,
    status,
    role,
    image,
    banner_image
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING 
    id,
    name,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
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
    image,
    banner_image,
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
    image,
    banner_image,
    created_at,
    updated_at
FROM users
WHERE email = $1;


-- name: PatchUserProfile :one
UPDATE users
SET
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    batch = COALESCE($4, batch),
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    created_at,
    updated_at;

-- name: PatchUserAdmin :one
UPDATE users
SET
    status = COALESCE($2, status),
    role = COALESCE($3, role),
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    created_at,
    updated_at;

-- name: PatchUserPassword :one
UPDATE users
SET
    password = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    updated_at;


-- name: PatchUserImages :one
UPDATE users
SET
    image = COALESCE($2, image),
    banner_image = COALESCE($3, banner_image),
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    image,
    banner_image,
    updated_at;
