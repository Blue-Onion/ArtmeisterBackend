-- +goose Up
CREATE TYPE public.art_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE art (
    id           uuid PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT,
    image        TEXT NOT NULL,
    tags         TEXT[] NOT NULL DEFAULT '{}',
    status       public.art_status NOT NULL DEFAULT 'pending',
    user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS art;
DROP TYPE IF EXISTS public.art_status;
