package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
)

const sqlListOrderQuery string = `
	SELECT
	    "order".id,
	    id_customer,
	    id_pickpoint,
	    shelf_life,
	    time_created,
	    date_receipt,
	    so.type AS id_state
	FROM "order" LEFT JOIN state_order so on so.id = "order".id_state
`

func (repo *Repository) List(ctx context.Context) ([]repository.OrderDTO, error) {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	var result []repository.OrderDTO
	rows, err := queryEngine.Query(ctx, sqlListOrderQuery)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var client repository.OrderDTO
		err := client.LoadFromRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, client)
	}
	if repo.cache != nil {
		for _, client := range result {
			err = repo.cache.Set(ctx, string(client.ID.Bytes[:]), client)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
