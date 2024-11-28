CREATE TYPE "code_scope" AS ENUM ('reset', 'confirm');

CREATE TABLE IF NOT EXISTS "codes" (
    "id"            UUID                            PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id"       UUID                            NOT NULL,
    "hash"          BYTEA                           NOT NULL UNIQUE,
    "scope"         "code_scope"                    NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE        NOT NULL DEFAULT NOW(),
    "expires_at"    TIMESTAMP WITH TIME ZONE        NOT NULL,
    
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE INDEX "idx_codes_user_id" ON codes("user_id");