-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens (
        id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        user_id BIGINT NOT NULL,
        token varchar(255) NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL,
        CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE

);

CREATE INDEX IF NOT EXISTS refresh_tokens_token_idx ON refresh_tokens USING btree (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
