CREATE TYPE "token_scope" AS ENUM ('access', 'refresh');

CREATE TYPE "token_status" AS ENUM ('active', 'revoked', 'used');

CREATE TABLE IF NOT EXISTS "tokens" (
    "id"            UUID                            PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id"       UUID                            NOT NULL,
    "hash"          BYTEA                           NOT NULL UNIQUE,
    "branch"        UUID                            NOT NULL,
    "status"        "token_status"                  NOT NULL DEFAULT 'active',
    "scope"         "token_scope"                   NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE        NOT NULL DEFAULT NOW(),
    "expires_at"    TIMESTAMP WITH TIME ZONE        NOT NULL,
    
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE INDEX "idx_tokens_user_id" ON tokens("user_id");