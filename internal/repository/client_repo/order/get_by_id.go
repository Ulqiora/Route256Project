package cli

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlGetOrderQuery string = `
	SELECT
	    "order".id,
	    id_customer,
	    id_pickpoint,
	    shelf_life,
	    time_created,
	    date_receipt,
	    penny AS price,
	    weight,
	    so.type AS id_state
	FROM "order" LEFT JOIN state_order so on so.id = "order".id_state
	WHERE "order".id = $1
`

func (repo *Repository) GetByID(ctx context.Context, id pgtype.UUID) (repository.OrderDTO, error) {
	var orderDto repository.OrderDTO
	if repo.cache != nil {
		bytesObj, err := repo.cache.Get(ctx, hashOrder(string(id.Bytes[:])))
		if err != nil {
			slog.Info(err.Error())
		} else {
			err = json.Unmarshal(bytesObj, &orderDto)
			if err != nil {
				return orderDto, err
			}
		}
	}
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	err := queryEngine.QueryRow(ctx, sqlGetOrderQuery, id).Scan(orderDto)
	if err != nil {
		return repository.OrderDTO{}, err
	}
	if repo.cache != nil {
		if err = repo.cache.Set(ctx, hashOrder(string(orderDto.ID.Bytes[:])), orderDto); err != nil {
			slog.Warn(err.Error())
		}
	}
	return orderDto, nil
}
