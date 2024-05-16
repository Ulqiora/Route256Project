package pickpoint

import (
	"context"

	"github.com/jackc/pgtype"
)

func (c *Controller) Delete(ctx context.Context, id string) error {
	var pickpointUuid pgtype.UUID
	err := pickpointUuid.Scan(id)
	if err != nil {
		return err
	}
	err = c.storage.Delete(ctx, pickpointUuid)
	if err != nil {
		return err
	}
	return err
}
