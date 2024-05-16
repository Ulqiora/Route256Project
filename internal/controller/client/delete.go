package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
)

func (c *Controller) Delete(ctx context.Context, ObjId string) error {
	var id pgtype.UUID
	err := id.Set(ObjId)
	if err != nil {
		return fmt.Errorf("controller: failed to set UUID: %w", err)
	}
	return c.storage.Delete(ctx, id)
}
