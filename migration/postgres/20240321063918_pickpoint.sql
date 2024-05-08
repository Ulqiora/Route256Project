-- +goose Up
-- +goose StatementBegin
CREATE TABLE pickpoint(
    id UUID DEFAULT uuid_generate_v4() UNIQUE ,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted BOOL NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pickpoint;
-- +goose StatementEnd
