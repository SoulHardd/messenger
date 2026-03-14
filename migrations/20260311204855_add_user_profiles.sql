-- +goose Up
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    nickname   TEXT,
    bio        TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_users_login ON users (login);

-- +goose Down
DROP TABLE IF EXISTS user_profiles;

DROP INDEX IF EXISTS idx_users_login
