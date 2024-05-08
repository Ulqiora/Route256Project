package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"homework/internal/repository"
)

const sqlGetContactDetailsQuery = `
    SELECT
        tc.type AS type,
        contact_detail.detail as detail
    FROM contact_detail
    LEFT OUTER JOIN type_contact tc on tc.id = contact_detail.id_type_contact
    WHERE contact_detail.id_pickpoint = $1 
`

const sqlGetPickPointQuery = `
    SELECT
        id,name,address
    FROM pickpoint
    WHERE id = $1 AND deleted=false
`

func (r *PickPointRepository) GetByID(ctx context.Context, id int) (repository.PickPointDTO, error) {
	var result repository.PickPointDTO
	if r.cache != nil {
		bytesData, err := r.cache.Get(ctx, hashFunction(id))
		if err == nil {
			err = json.Unmarshal(bytesData, &result)
			if err != nil {
				return result, err
			}
			return result, err
		}
	}
	queryEngine := r.db.GetQueryEngine(ctx)
	err := queryEngine.Get(ctx, &result, sqlGetPickPointQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.PickPointDTO{}, ErrorObjectNotFounded
		}
		return repository.PickPointDTO{}, err
	}
	err = queryEngine.Select(ctx, &result.ContactDetails, sqlGetContactDetailsQuery, id)
	if err != nil {
		return repository.PickPointDTO{}, err
	}
	if err = r.cache.Set(ctx, hashFunction(id), result); err != nil {
		return result, err
	}
	return result, nil
}
