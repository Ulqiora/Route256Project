package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
)

const sqlListPickPointQuery string = `
     SELECT
         id,name,address
     FROM pickpoint
     WHERE deleted=false
`

func (r *PickPointRepository) List(ctx context.Context) ([]repository.PickPointDTO, error) {
	queryEngine := r.manager.DefaultTrOrDB(ctx, r.db.GetPool(ctx))
	var result []repository.PickPointDTO
	rows, err := queryEngine.Query(ctx, sqlListPickPointQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pp repository.PickPointDTO
		err = rows.Scan(&pp.ID, &pp.Name, &pp.Address)
		if err != nil {
			return nil, err
		}
		result = append(result, pp)
	}
	for i := range result {
		rowsContacts, err := queryEngine.Query(ctx, sqlGetContactDetailsQuery, result[i].ID)
		if err != nil {
			return nil, err
		}
		rowsContacts.Close()
		for rowsContacts.Next() {
			var contact repository.ContactDetailDTO
			err = rowsContacts.Scan(contact)
			if err != nil {
				return nil, err
			}
			result[i].ContactDetails = append(result[i].ContactDetails, contact)
		}
		rowsContacts.Close()
	}

	return result, nil
}
