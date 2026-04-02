-- +goose Up
UPDATE messages_status
SET status = 1
WHERE status IS NULL;

ALTER TABLE messages_status
    ALTER COLUMN status SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_messages_status_user_message
    ON messages_status(user_id, message_id);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_status_user_message;

ALTER TABLE messages_status
    ALTER COLUMN status DROP NOT NULL;
