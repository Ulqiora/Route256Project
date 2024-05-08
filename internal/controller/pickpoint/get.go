package pickpoint

import (
	"context"

	"homework/internal/model"
)

func (c *Controller) GetByID(ctx context.Context, id uint64) (model.PickPoint, error) {
	var obj model.PickPoint
	dto, err := c.storage.GetByID(ctx, int(id))
	if err != nil {
		return model.PickPoint{}, err
	}
	return obj.LoadFromDTO(dto), nil
}
