package pickpoint

import (
	"context"

	"homework/internal/model"
)

func (c *Controller) Create(ctx context.Context, object model.PickPoint) (uint64, error) {
	id, err := c.storage.Create(ctx, object.MapToDTO())
	if err != nil {
		return 0, err
	}
	return id, nil
}
