-- +goose Up
CREATE TYPE public.art_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE art(
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT NOT NULL,
    description  TEXT DEFAULT '',
    image        TEXT DEFAULT '',
    tags         TEXT[] NOT NULL DEFAULT '{}',
    status       public.art_status NOT NULL DEFAULT 'pending',
    user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS art;
DROP TYPE IF EXISTS public.art_status;
