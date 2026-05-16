-- name: CreateArt :one
INSERT INTO art (name, description, image, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetArtByID :one
SELECT * FROM art
WHERE id = $1;

-- name: GetArtByUser :many
SELECT * FROM art
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListArt :many
SELECT * FROM art
ORDER BY created_at DESC;

-- name: UpdateArt :one
UPDATE art
SET
    name        = COALESCE($2, name),
    description = COALESCE($3, description),
    image       = COALESCE($3, image),
    updated_at  = NOW()
WHERE id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteArt :exec
DELETE FROM art
WHERE id = $1 AND user_id = $2;
