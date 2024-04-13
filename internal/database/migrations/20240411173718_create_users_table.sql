-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
     id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
     username VARCHAR(50) UNIQUE NOT NULL,
     first_name VARCHAR(50) NOT NULL,
     last_name VARCHAR(50) NOT NULL,
     contact_number VARCHAR(20) NOT NULL,
     email varchar(255) UNIQUE NOT NULL,
     password varchar(255) NOT NULL,
     is_admin BOOLEAN DEFAULT '0' NOT NULL,
     created_at TIMESTAMP NOT NULL,
     updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
