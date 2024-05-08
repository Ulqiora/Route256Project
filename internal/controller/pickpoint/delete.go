package pickpoint

import (
	"context"
)

func (c *Controller) Delete(ctx context.Context, id uint64) error {
	err := c.storage.Delete(ctx, int(id))
	return err
}
