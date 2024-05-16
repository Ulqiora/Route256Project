package cli

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/jackc/pgtype"
)

const sqlOrderSetToDeletedQuery string = `
	UPDATE "order" SET id_state=$2 WHERE id=$1
`

func (repo *Repository) Delete(ctx context.Context, id pgtype.UUID) error {
	queryEngine := repo.manager.DefaultTrOrDB(ctx, repo.db.GetPool(ctx))
	stateID, err := repo.getState(ctx, model.EDeleted)
	if err != nil {
		return err
	}
	_, err = queryEngine.Exec(ctx, sqlOrderSetToDeletedQuery, id, stateID)
	if err != nil {
		return err
	}
	return nil
}
