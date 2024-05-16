package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlDeleteClientQuery string = `
      UPDATE pickpoint SET deleted=TRUE WHERE id=$1
`

func (repo *ClientRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	err := connection.QueryRow(ctx, sqlDeleteClientQuery, id)
	if err != nil {
		return repository.ErrorObjectNotCreated
	}
	return nil
}
