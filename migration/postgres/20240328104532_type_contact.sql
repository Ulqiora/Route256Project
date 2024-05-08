-- +goose Up
-- +goose StatementBegin
CREATE TABLE type_contact(
    id UUID DEFAULT uuid_generate_v4() UNIQUE,
    type TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE type_contact;
-- +goose StatementEnd
