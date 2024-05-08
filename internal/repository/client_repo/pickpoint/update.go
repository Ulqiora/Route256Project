package cli

import (
	"context"

	"homework/internal/repository"
)

const sqlUpdatePickPointQuery string = `
    UPDATE pickpoint SET name=$2, address=$3 WHERE id=$1
`
const sqlDetelePickPointContactsQuery string = `
    DELETE FROM contact_detail WHERE id_pickpoint=$1
`

func (r *PickPointRepository) Update(ctx context.Context, dto repository.PickPointDTO) (int, error) {
	queryEngine := r.db.GetQueryEngine(ctx)
	_, err := queryEngine.Exec(ctx, sqlUpdatePickPointQuery, dto.ID, dto.Name, dto.Address)
	if err != nil {
		return 0, err
	}
	_, err = queryEngine.Exec(ctx, sqlDetelePickPointContactsQuery, dto.ID)
	if err != nil {
		return 0, err
	}
	for i := range dto.ContactDetails {
		var idType int
		err = queryEngine.ExecQueryRow(ctx, sqlCreateContactTypeQuery, dto.ContactDetails[i].Type).Scan(&idType)
		if err != nil {
			return 0, err
		}
		_, err = queryEngine.Exec(ctx, sqlCreateContactDetailsQuery, idType, dto.ContactDetails[i].Detail, dto.ID)
		if err != nil {
			return 0, err
		}
	}
	if r.cache != nil {
		if err = r.cache.Set(ctx, hashFunction(dto.ID), dto); err != nil {
			return dto.ID, err
		}
	}
	return dto.ID, err
}
