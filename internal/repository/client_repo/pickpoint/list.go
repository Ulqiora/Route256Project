package cli

import (
	"context"

	"homework/internal/repository"
)

const sqlListPickPointQuery string = `
     SELECT
         id,name,address
     FROM pickpoint
     WHERE deleted=false
`

func (r *PickPointRepository) List(ctx context.Context) ([]repository.PickPointDTO, error) {
	queryEngine := r.db.GetQueryEngine(ctx)
	var result []repository.PickPointDTO
	err := queryEngine.Select(ctx, &result, sqlListPickPointQuery)
	if err != nil {
		return nil, err
	}
	for i := range result {
		err = queryEngine.Select(ctx, &result[i].ContactDetails, sqlGetContactDetailsQuery, result[i].ID)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
