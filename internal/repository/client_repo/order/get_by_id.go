package cli

import (
	"context"
	"encoding/json"
	"log/slog"

	"homework/internal/repository"
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

func (r *Repository) GetByID(ctx context.Context, id uint64) (repository.OrderDTO, error) {
	var orderDto repository.OrderDTO
	if r.cache != nil {
		bytesObj, err := r.cache.Get(ctx, hashOrder(id))
		if err != nil {
			slog.Info(err.Error())
		} else {
			err = json.Unmarshal(bytesObj, &orderDto)
			if err != nil {
				return orderDto, err
			}
		}
	}
	queryEngine := r.manager.GetQueryEngine(ctx)
	err := queryEngine.Get(ctx, &orderDto, sqlGetOrderQuery, id)
	if err != nil {
		return orderDto, err
	}
	if r.cache != nil {
		if err = r.cache.Set(ctx, hashOrder(orderDto.ID), orderDto); err != nil {
			slog.Warn(err.Error())
		}
	}
	return orderDto, nil
}
