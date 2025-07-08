-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    date VARCHAR NOT NULL,
    amount NUMERIC NOT NULL,
    category_id UUID,
    account_id UUID NOT NULL,
    type VARCHAR,
    created_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
