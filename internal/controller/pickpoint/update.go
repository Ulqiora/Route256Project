package pickpoint

import (
	"context"

	"homework/internal/model"
)

func (c *Controller) Update(ctx context.Context, object model.PickPoint) (uint64, error) {
	id, err := c.storage.Update(ctx, object.MapToDTO())
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
