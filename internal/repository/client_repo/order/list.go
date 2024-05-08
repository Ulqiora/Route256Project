package cli

import (
	"context"
	"fmt"

	"homework/internal/repository"
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

func (r *Repository) List(ctx context.Context) ([]repository.OrderDTO, error) {
	queryEngine := r.manager.GetQueryEngine(ctx)
	var orders []repository.OrderDTO
	err := queryEngine.Select(ctx, &orders, sqlListOrderQuery)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", repository.ErrorDataBase, err)
	}
	return orders, nil
}
