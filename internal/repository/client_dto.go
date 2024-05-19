package repository

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type ClientDTO struct {
	ID   pgtype.UUID `bd:"id"`
	Name string      `bd:"name"`
}

func (c *ClientDTO) LoadFromRow(row pgx.Row) error {
	err := row.Scan(&c.ID, &c.Name)
	return err
}
