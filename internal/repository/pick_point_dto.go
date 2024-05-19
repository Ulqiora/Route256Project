package repository

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type PickPointDTO struct {
	ID             pgtype.UUID        `json:"id"              db:"id"`
	Name           string             `json:"name"            db:"name"`
	Address        string             `json:"address"         db:"address"`
	ContactDetails []ContactDetailDTO `json:"contact_details" db:"contact_details"`
}

func (p *PickPointDTO) LoadFromRow(row pgx.Row) error {
	err := row.Scan(p.ID, &p.Name, &p.Address)
	return err
}
