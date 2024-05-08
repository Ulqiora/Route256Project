package repository

import "github.com/jackc/pgtype"

type Client struct {
	ID   pgtype.UUID `bd:"id"`
	Name string      `bd:"name"`
}
