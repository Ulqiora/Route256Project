package cli

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlGetStateQuery string = `
	SELECT
	    id
	FROM state_order
	WHERE type = $1
`

func (repo *Repository) getState(ctx context.Context, state string) (pgtype.UUID, error) {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	var idState pgtype.UUID
	if repo.cache != nil {
		bytesdata, err := repo.cache.Get(ctx, hashStateOrder(state))
		if err != nil {
			slog.Info(err.Error())
		} else {
			err = json.Unmarshal(bytesdata, &idState)
			return idState, err
		}
	}
	err := connection.QueryRow(ctx, sqlGetStateQuery, state).Scan(&idState)
	if err != nil {
		slog.Warn(err.Error())
		return idState, repository.ErrorNotFoundedStateOrder
	}
	if repo.cache != nil {
		err = repo.cache.Set(ctx, hashStateOrder(state), idState)
		if err != nil {
			slog.Info(err.Error())
		}
	}
	return idState, nil
}
