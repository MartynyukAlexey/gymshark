CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE "user_state" AS ENUM ('pending', 'active', 'deleted');

CREATE TABLE IF NOT EXISTS "users" (
    "id"                UUID                            PRIMARY KEY DEFAULT gen_random_uuid(),
    "email"             CITEXT                          NOT NULL UNIQUE,
    "password"          BYTEA                           NOT NULL,
    "state"             "user_state"                    NOT NULL DEFAULT 'pending',
    "avatar_id"         TEXT                            NOT NULL,
    "first_name"        TEXT                            NOT NULL,
    "last_name"         TEXT                            NOT NULL,
    "created_at"        TIMESTAMP WITH TIME ZONE        NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMP WITH TIME ZONE        NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_users_email" ON "users" ("email");