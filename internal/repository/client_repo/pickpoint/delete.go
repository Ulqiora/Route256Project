package cli

import (
	"context"

	"github.com/jackc/pgtype"
)

const sqlDeletePickPointQuery string = `
      UPDATE pickpoint SET deleted=TRUE WHERE id=$1
`

func (r *PickPointRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	queryEngine := r.manager.DefaultTrOrDB(ctx, r.db.GetPool(ctx))
	_, err := queryEngine.Exec(ctx, sqlDeletePickPointQuery, id)
	return err
}
