package cli

import (
	"context"
	"fmt"

	"homework/internal/model"
	"homework/internal/repository"
)

const sqlCreateOrderQuery string = `
	INSERT INTO "order"(id_customer,id_pickpoint,id_state,penny,weight,shelf_life) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
`

func (r *Repository) Create(ctx context.Context, dto repository.OrderDTO) (uint64, error) {
	queryEngine := r.manager.GetQueryEngine(ctx)

	idState, err := r.getState(ctx, model.EReadyToIssued)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", repository.ErrorOrderNotCreated, err)
	}
	var id uint64
	err = queryEngine.ExecQueryRow(ctx, sqlCreateOrderQuery, dto.CustomerID, dto.PickPointID, idState, dto.Penny, dto.Weight, dto.ShelfLife).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", repository.ErrorOrderNotCreated, err)
	}
	return id, nil
}
