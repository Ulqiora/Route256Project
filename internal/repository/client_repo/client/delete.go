package cli

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlDeleteClientQuery string = `
      UPDATE pickpoint SET deleted=TRUE WHERE id=$1
`

func (repo *ClientRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	connection := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	_, err := connection.Exec(ctx, sqlDeleteClientQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %s", repository.ErrorObjectNotCreated, err)
	}
	return nil
}
