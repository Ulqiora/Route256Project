-- +goose Up
-- +goose StatementBegin
select * from pg_extension;
CREATE EXTENSION "uuid-ossp";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION "uuid-ossp" CASCADE;
-- +goose StatementEnd
