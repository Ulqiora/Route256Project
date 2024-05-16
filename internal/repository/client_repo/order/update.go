package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

const sqlOrderSetToReceivedQuery string = `
	UPDATE "order" SET id_state=(SELECT id FROM state_order WHERE type = 'Received'), date_receipt=CURRENT_TIMESTAMP WHERE id=$1
`
const sqlOrderSetToReturnedQuery string = `
	UPDATE "order" SET id_state=(SELECT id FROM state_order WHERE type = 'Returned') WHERE id=$1
`

func (repo *Repository) UpdateToReceived(ctx context.Context, orderID pgtype.UUID) error {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	_, err := queryEngine.Exec(ctx, sqlOrderSetToReceivedQuery, orderID)
	if err != nil {
		return repository.ErrorDataBase
	}
	return nil
}

func (repo *Repository) UpdateToReturned(ctx context.Context, orderID pgtype.UUID) error {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	_, err := queryEngine.Exec(ctx, sqlOrderSetToReturnedQuery, orderID)
	if err != nil {
		return repository.ErrorDataBase
	}
	return nil
}
