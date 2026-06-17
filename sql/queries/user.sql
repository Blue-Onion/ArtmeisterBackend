-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password
)
VALUES (
    $1, $2, $3
)
RETURNING 
    id,
    name,
    username,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at;


-- name: GetUser :one
SELECT 
    id,
    name,
    username,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at
FROM users
WHERE id = $1;


-- name: GetAllUser :many
SELECT 
    id,
    name,
    username,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at
FROM users;


-- name: GetUserByEmail :one
SELECT 
    id,
    name,
    username,
    email,
    password,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at
FROM users
WHERE email = $1;


-- name: GetUserByUsername :one
SELECT 
    id,
    name,
    username,
    email,
    password,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at
FROM users
WHERE username = $1;


-- name: PatchUserProfile :one
UPDATE users
SET
    name = COALESCE(sqlc.narg('name')::text, name),
    username = COALESCE(sqlc.narg('username')::text, username),
    email = COALESCE(sqlc.narg('email')::text, email),
    batch = COALESCE(sqlc.narg('batch')::text, batch),
    description = COALESCE(sqlc.narg('description')::text, description),
    image = COALESCE(sqlc.narg('image')::text, image),
    banner_image = COALESCE(sqlc.narg('banner_image')::text, banner_image),
    social_links = COALESCE(sqlc.narg('social_links')::jsonb, social_links),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING
    id,
    name,
    username,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at;



-- name: PatchUserPassword :one

UPDATE users
SET
    password = sqlc.arg('password'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING
    id,
    name,
    username,
    email,
    updated_at;




-- name: PatchUserAdmin :one
UPDATE users
SET
    status = COALESCE(sqlc.narg('status')::account_status, status),
    role = COALESCE(sqlc.narg('role')::user_role, role),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING
    id,
    name,
    username,
    email,
    batch,
    status,
    role,
    image,
    banner_image,
    description,
    social_links,
    created_at,
    updated_at;
