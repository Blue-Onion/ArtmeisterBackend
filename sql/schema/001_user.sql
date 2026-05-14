-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE public.user_role AS ENUM ('admin', 'user', 'moderator');
CREATE TYPE public.account_status AS ENUM ('pending', 'approved', 'banned');

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,

    banner_image TEXT DEFAULT '',
    image TEXT DEFAULT '',
    batch TEXT DEFAULT '',

    status public.account_status NOT NULL DEFAULT 'pending',
    role public.user_role NOT NULL DEFAULT 'user',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS public.account_status;
DROP TYPE IF EXISTS public.user_role;
