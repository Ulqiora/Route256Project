-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order"
(
    id UUID DEFAULT uuid_generate_v4() UNIQUE,
    id_customer UUID NOT NULL,
    id_pickpoint UUID NOT NULL,
    shelf_life TIMESTAMP WITH TIME ZONE,
    time_created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    date_receipt TIMESTAMP WITH TIME ZONE,

    penny NUMERIC NOT NULL,
    weight NUMERIC NOT NULL,

    id_state UUID NOT NULL,
    FOREIGN KEY (id_pickpoint) REFERENCES pickpoint(id),
    FOREIGN KEY (id_state) REFERENCES state_order(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "order";
-- +goose StatementEnd
