package cli

import (
	"context"
	"encoding/json"
	"log/slog"

	"homework/internal/repository"
)

const sqlGetStateQuery string = `
	SELECT
	    id
	FROM state_order
	WHERE type = $1
`

func (r *Repository) getState(ctx context.Context, state string) (uint64, error) {
	queryEngine := r.manager.GetQueryEngine(ctx)
	var idState uint64
	if r.cache != nil {
		bytesdata, err := r.cache.Get(ctx, hashStateOrder(state))
		if err != nil {
			slog.Info(err.Error())
		} else {
			err = json.Unmarshal(bytesdata, &idState)
			return idState, err
		}
	}
	err := queryEngine.ExecQueryRow(ctx, sqlGetStateQuery, state).Scan(&idState)
	if err != nil {
		slog.Warn(err.Error())
		return idState, repository.ErrorNotFoundedStateOrder
	}
	if r.cache != nil {
		err = r.cache.Set(ctx, hashStateOrder(state), idState)
		if err != nil {
			slog.Info(err.Error())
		}
	}
	return idState, nil
}
