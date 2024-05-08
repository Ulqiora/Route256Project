-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_detail
(
    id UUID DEFAULT uuid_generate_v4() UNIQUE,
    id_type_contact UUID NOT NULL ,
    id_pickpoint UUID NOT NULL ,
    detail TEXT NOT NULL,
    FOREIGN KEY (id_type_contact) REFERENCES type_contact(id),
    FOREIGN KEY (id_pickpoint) REFERENCES pickpoint(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contact_detail;
-- +goose StatementEnd
