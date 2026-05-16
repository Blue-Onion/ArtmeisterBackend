-- +goose Up
CREATE TABLE art(
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        DEFAULT '',
    image       TEXT        DEFAULT '',
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS art;
