package cli

import (
	"context"

	"homework/internal/model"
)

const sqlOrderSetToDeletedQuery string = `
	UPDATE "order" SET id_state=$2 WHERE id=$1
`

func (r *Repository) Delete(ctx context.Context, ID uint64) error {
	queryEngine := r.manager.GetQueryEngine(ctx)
	stateID, err := r.getState(ctx, model.EDeleted)
	if err != nil {
		return err
	}
	_, err = queryEngine.Exec(ctx, sqlOrderSetToDeletedQuery, ID, stateID)
	if err != nil {
		return err
	}
	return nil
}
