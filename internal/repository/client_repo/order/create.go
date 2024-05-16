package cli

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlCreateOrderQuery string = `
	INSERT INTO "order"(id_customer,id_pickpoint,id_state,penny,weight,shelf_life) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
`

func (repo *Repository) Create(ctx context.Context, dto repository.OrderDTO) (pgtype.UUID, error) {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))

	idState, err := repo.getState(ctx, model.EReadyToIssued)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("%s: %s", repository.ErrorOrderNotCreated, err)
	}
	var id pgtype.UUID
	err = queryEngine.QueryRow(ctx, sqlCreateOrderQuery, dto.CustomerID, dto.PickPointID, idState, dto.Penny, dto.Weight, dto.ShelfLife).Scan(&id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("%s: %s", repository.ErrorOrderNotCreated, err)
	}
	return id, nil
}
