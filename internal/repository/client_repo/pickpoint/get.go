package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
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

func (r *PickPointRepository) GetByID(ctx context.Context, id pgtype.UUID) (repository.PickPointDTO, error) {
	var result repository.PickPointDTO
	if r.cache != nil {
		bytesData, err := r.cache.Get(ctx, hashFunction(string(id.Bytes[:])))
		if err == nil {
			err = json.Unmarshal(bytesData, &result)
			if err != nil {
				return result, err
			}
			return result, err
		}
	}
	queryEngine := r.manager.DefaultTrOrDB(ctx, r.db.GetPool(ctx))
	row := queryEngine.QueryRow(ctx, sqlGetPickPointQuery, id)
	err := result.LoadFromRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.PickPointDTO{}, ErrorObjectNotFounded
		}
		return repository.PickPointDTO{}, err
	}
	rows, err := queryEngine.Query(ctx, sqlGetContactDetailsQuery, id)
	defer rows.Close()
	for rows.Next() {
		var contact repository.ContactDetailDTO
		err = rows.Scan(&contact.Type, &contact.Detail)
		if err != nil {
			return repository.PickPointDTO{}, err
		}
		result.ContactDetails = append(result.ContactDetails, contact)
	}
	if err != nil {
		return repository.PickPointDTO{}, err
	}
	if err = r.cache.Set(ctx, hashFunction(string(id.Bytes[:])), result); err != nil {
		return result, err
	}
	return result, nil
}
