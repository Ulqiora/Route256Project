-- +goose Up
-- +goose StatementBegin
CREATE TABLE client
(
    id UUID DEFAULT uuid_generate_v4() UNIQUE,
    name TEXT NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE client;
-- +goose StatementEnd
