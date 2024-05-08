package cli

import (
	"context"
	"fmt"

	"homework/internal/repository"
)

const sqlGetReadyToIssuedQuery = `
	SELECT
		"order".id,
	    id_customer,
	    id_pickpoint,
	    shelf_life,
	    time_created,
	    date_receipt,
	    penny as price,
	    weight,
	    so.type AS id_state
	FROM "order"
	LEFT JOIN public.state_order so on so.id = "order".id_state
	WHERE shelf_life > CURRENT_TIMESTAMP AND id_state = (
				SELECT id
				FROM state_order
				WHERE type = 'ReadyToIssued') `

func (r *Repository) ListReadyToIssued(ctx context.Context) ([]repository.OrderDTO, error) {
	queryEngine := r.manager.GetQueryEngine(ctx)
	var dtos []repository.OrderDTO
	err := queryEngine.Select(ctx, &dtos, sqlGetReadyToIssuedQuery)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", repository.ErrorDataBase, err)
	}
	return dtos, nil
}
