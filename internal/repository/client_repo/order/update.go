package cli

import (
	"context"

	"homework/internal/repository"
)

const sqlOrderSetToReceivedQuery string = `
	UPDATE "order" SET id_state=(SELECT id FROM state_order WHERE type = 'Received'), date_receipt=CURRENT_TIMESTAMP WHERE id=$1
`
const sqlOrderSetToReturnedQuery string = `
	UPDATE "order" SET id_state=(SELECT id FROM state_order WHERE type = 'Returned') WHERE id=$1
`

func (r *Repository) UpdateToReceived(ctx context.Context, orderID uint64) error {
	queryEngine := r.manager.GetQueryEngine(ctx)
	_, err := queryEngine.Exec(ctx, sqlOrderSetToReceivedQuery, orderID)
	if err != nil {
		return repository.ErrorDataBase
	}
	return nil
}

func (r *Repository) UpdateToReturned(ctx context.Context, orderID uint64) error {
	queryEngine := r.manager.GetQueryEngine(ctx)
	_, err := queryEngine.Exec(ctx, sqlOrderSetToReturnedQuery, orderID)
	if err != nil {
		return repository.ErrorDataBase
	}
	return nil
}
