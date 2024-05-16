package cli

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/repository"
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

func (repo *Repository) ListReadyToIssued(ctx context.Context) ([]repository.OrderDTO, error) {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	var dtos []repository.OrderDTO
	rows, err := queryEngine.Query(ctx, sqlGetReadyToIssuedQuery)
	rows.Close()
	for rows.Next() {
		var dto repository.OrderDTO
		err := dto.LoadFromRow(rows)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}

	if err != nil {
		return nil, fmt.Errorf("%s, %s", repository.ErrorDataBase, err)
	}
	return dtos, nil
}
