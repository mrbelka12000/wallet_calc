-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR NOT NULL,
    password VARCHAR,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP default now()
);


CREATE TABLE IF NOT EXISTS accounts(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR NOT NULL,
    balance NUMERIC,
    parser_id UUID NOT NULL,
    created_at TIMESTAMP default now(),
    CONSTRAINT fk_accounts_users foreign key (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS categories(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
