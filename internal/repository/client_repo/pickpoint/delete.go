package cli

import (
	"context"
)

const sqlDeletePickPointQuery string = `
      UPDATE pickpoint SET deleted=TRUE WHERE id=$1
`

func (r *PickPointRepository) Delete(ctx context.Context, id int) error {
	queryEngine := r.db.GetQueryEngine(ctx)
	_, err := queryEngine.Exec(ctx, sqlDeletePickPointQuery, id)
	return err
}
