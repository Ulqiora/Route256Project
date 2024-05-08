package cli

import (
	"context"
	"fmt"

	"homework/internal/repository"
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

func (r *PickPointRepository) Create(ctx context.Context, dto repository.PickPointDTO) (uint64, error) {
	queryEngine := r.db.GetQueryEngine(ctx)
	var id uint64
	err := queryEngine.ExecQueryRow(ctx, sqlCreatePickPointQuery, dto.Name, dto.Address).Scan(&id)
	if err != nil {
		return 0, err
	}
	ids := make([]int, len(dto.ContactDetails))
	for i, contact := range dto.ContactDetails {
		var idType int
		err := queryEngine.ExecQueryRow(ctx, sqlCreateContactTypeQuery, contact.Type).Scan(&idType)
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
		err = queryEngine.ExecQueryRow(ctx, sqlCreateContactDetailsQuery, idType, contact.Detail, id).Scan(&ids[i])
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
	}
	return id, err
}
