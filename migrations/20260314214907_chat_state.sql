-- +goose Up
CREATE TABLE IF NOT EXISTS chat_state (
    chat_id            BIGINT REFERENCES chats(id) ON DELETE CASCADE,
    user_id              BIGINT REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id BIGINT REFERENCES messages(id),
    unread_count         INT DEFAULT 0,
    last_message_id BIGINT,
    last_message_text TEXT,
    last_message_at TIMESTAMP,
    PRIMARY KEY(chat_id, user_id)
);

ALTER TABLE chats ADD COLUMN title TEXT;

CREATE INDEX IF NOT EXISTS idx_user_chats_user_time ON chat_state(user_id, last_message_at DESC);

-- +goose Down
DROP TABLE IF EXISTS chat_state;
ALTER TABLE chats DROP COLUMN title;
DROP INDEX IF EXISTS idx_user_chats_user_time;