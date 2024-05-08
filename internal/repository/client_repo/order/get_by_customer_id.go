package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"homework/internal/repository"
)

const sqlGetOrderByCustomerQuery string = `
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
	FROM "order" LEFT JOIN state_order so on so.id = "order".id_state
	WHERE "order".id_customer = $1
`

func (r *Repository) GetByCustomerID(ctx context.Context, id uint64) ([]repository.OrderDTO, error) {
	queryEngine := r.manager.GetQueryEngine(ctx)
	var result []repository.OrderDTO
	err := queryEngine.Select(ctx, &result, sqlGetOrderByCustomerQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s, %s", repository.ErrorOrderNotFounded, err)
		}
		return nil, fmt.Errorf("%s, %s", repository.ErrorDataBase, err)
	}
	if r.cache != nil {
		if err = r.cacheMulti(ctx, result); err != nil {
			slog.Warn(err.Error())
			return result, nil
		}
	}
	return result, nil
}

func (r *Repository) cacheMulti(ctx context.Context, orders []repository.OrderDTO) error {
	ordersMap := make(map[string]any, len(orders))
	for _, orderObj := range orders {
		ordersMap[hashOrder(orderObj.ID)] = orderObj
	}
	return r.cache.SetMulti(ctx, ordersMap)
}
