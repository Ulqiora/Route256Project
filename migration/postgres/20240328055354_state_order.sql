-- +goose Up
-- +goose StatementBegin

CREATE TABLE state_order
(
    id UUID default uuid_generate_v1() UNIQUE,
    type TEXT NOT NULL
);
INSERT INTO state_order(type)VALUES ('ReadyToIssued');
INSERT INTO state_order(type)VALUES ('Received');
INSERT INTO state_order(type)VALUES ('Returned');
INSERT INTO state_order(type)VALUES ('Deleted');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE state_order;
-- +goose StatementEnd
