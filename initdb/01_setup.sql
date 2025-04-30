CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id             UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    email          VARCHAR(50)     NOT NULL UNIQUE,
    phone          VARCHAR(15)     NOT NULL UNIQUE,
    password_hash  VARCHAR(100)    NOT NULL,
    is_deleted     BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ     NULL
);
