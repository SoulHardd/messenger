-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    phone         TEXT NOT NULL UNIQUE,
    login         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL,
    expires_at         TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id         BIGSERIAL PRIMARY KEY,
    type       SMALLINT NOT NULL, -- 1=private; 2=group
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by BIGINT REFERENCES users(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_participants (
    chat_id BIGINT REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ,
    role    SMALLINT NOT NULL, -- 1=member; 2=admin
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT REFERENCES  chats(id) ON DELETE CASCADE NOT NULL,
    sender_id  BIGINT REFERENCES users(id) NOT NULL,
    content    TEXT,
    is_edited  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS messages_status (
    message_id BIGINT REFERENCES messages(id),
    user_id    BIGINT REFERENCES users(id),
    status     SMALLINT, -- 1=sent; 2=delivered; 3=read
    updated_at TIMESTAMPTZ default now(),
    PRIMARY KEY (message_id, user_id)
);

CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at DESC);
CREATE INDEX idx_chat_participants_user ON chat_participants(user_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS chat_participants;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS messages_status;

DROP INDEX IF EXISTS idx_messages_chat_created;
DROP INDEX IF EXISTS idx_chat_participants_user;
DROP INDEX IF EXISTS idx_sessions_user_id