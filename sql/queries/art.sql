-- name: CreateArt :one
INSERT INTO art (id, name, description, image, tags, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
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
WHERE status = 'approved'
ORDER BY created_at DESC;


-- name: ListArtByTag :many
SELECT * FROM art
WHERE status = 'approved'
  AND $1 = ANY(tags)
ORDER BY created_at DESC;


-- name: ListArtByTags :many
SELECT * FROM art
WHERE status = 'approved'
  AND tags && $1::text[]
ORDER BY created_at DESC;


-- name: UpdateArt :one
UPDATE art
SET
    name        = COALESCE($2, name),
    description = COALESCE($3, description),
    image       = COALESCE($4, image),
    tags        = COALESCE($5, tags),
    updated_at  = NOW()
WHERE id = $1 AND user_id = $6
RETURNING *;


-- name: DeleteArt :exec
DELETE FROM art
WHERE id = $1 AND user_id = $2;
