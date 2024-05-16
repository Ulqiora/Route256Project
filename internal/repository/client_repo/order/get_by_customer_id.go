package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
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

func (repo *Repository) GetByCustomerID(ctx context.Context, id pgtype.UUID) ([]repository.OrderDTO, error) {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	var result []repository.OrderDTO
	rows, err := queryEngine.Query(ctx, sqlGetOrderByCustomerQuery, id)
	for rows.Next() {
		var item repository.OrderDTO
		err = rows.Scan(item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s, %s", repository.ErrorOrderNotFounded, err)
		}
		return nil, fmt.Errorf("%s, %s", repository.ErrorDataBase, err)
	}
	if repo.cache != nil {
		if err = repo.cacheMulti(ctx, result); err != nil {
			slog.Warn(err.Error())
			return result, nil
		}
	}
	return result, nil
}

func (repo *Repository) cacheMulti(ctx context.Context, orders []repository.OrderDTO) error {
	ordersMap := make(map[string]any, len(orders))
	for _, orderObj := range orders {
		ordersMap[hashOrder(string(orderObj.ID.Bytes[:]))] = orderObj
	}
	return repo.cache.SetMulti(ctx, ordersMap)
}
