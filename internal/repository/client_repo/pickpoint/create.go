package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlCreateContactDetailsQuery = `
    INSERT INTO contact_detail(id_type_contact,detail,id_pickpoint) VALUES ($1,$2,$3) RETURNING id
`

const sqlCreateContactTypeQuery = `
    INSERT INTO type_contact(type) VALUES ($1) ON CONFLICT (type) DO UPDATE SET type = EXCLUDED.type RETURNING id;
`

const sqlCreatePickPointQuery = `
    INSERT INTO pickpoint(name,address) VALUES ($1,$2) RETURNING id 
`

func (r *PickPointRepository) Create(ctx context.Context, dto repository.PickPointDTO) (pgtype.UUID, error) {
	queryEngine := r.manager.DefaultTrOrDB(ctx, r.db.GetPool(ctx))
	var id pgtype.UUID
	err := queryEngine.QueryRow(ctx, sqlCreatePickPointQuery, dto.Name, dto.Address).Scan(&id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	ids := make([]int, len(dto.ContactDetails))
	for i, contact := range dto.ContactDetails {
		var idType string
		err := queryEngine.QueryRow(ctx, sqlCreateContactTypeQuery, contact.Type).Scan(&idType)
		if err != nil {
			return pgtype.UUID{}, err
		}
		err = queryEngine.QueryRow(ctx, sqlCreateContactDetailsQuery, idType, contact.Detail, id).Scan(&ids[i])
		if err != nil {
			return pgtype.UUID{}, err
		}
	}
	return id, err
}
