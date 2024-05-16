package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlUpdatePickPointQuery string = `
    UPDATE pickpoint SET name=$2, address=$3 WHERE id=$1
`
const sqlDetelePickPointContactsQuery string = `
    DELETE FROM contact_detail WHERE id_pickpoint=$1
`

func (r *PickPointRepository) Update(ctx context.Context, dto repository.PickPointDTO) (pgtype.UUID, error) {
	queryEngine := r.manager.DefaultTrOrDB(ctx, r.db.GetPool(ctx))
	_, err := queryEngine.Exec(ctx, sqlUpdatePickPointQuery, dto.ID, dto.Name, dto.Address)
	if err != nil {
		return pgtype.UUID{}, err
	}
	_, err = queryEngine.Exec(ctx, sqlDetelePickPointContactsQuery, dto.ID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	for i := range dto.ContactDetails {
		var idType int
		err = queryEngine.QueryRow(ctx, sqlCreateContactTypeQuery, dto.ContactDetails[i].Type).Scan(&idType)
		if err != nil {
			return pgtype.UUID{}, err
		}
		_, err = queryEngine.Exec(ctx, sqlCreateContactDetailsQuery, idType, dto.ContactDetails[i].Detail, dto.ID)
		if err != nil {
			return pgtype.UUID{}, err
		}
	}
	if r.cache != nil {
		if err = r.cache.Set(ctx, hashFunction(string(dto.ID.Bytes[:])), dto); err != nil {
			return dto.ID, err
		}
	}
	return dto.ID, err
}
