-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION create_profile_after_user()
    RETURNS trigger
    LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO user_profiles (user_id)
    VALUES (NEW.id);
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER user_profile_trigger
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION create_profile_after_user();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS user_profile_trigger ON users;
DROP FUNCTION IF EXISTS create_profile_after_user();
