package pickpoint

import (
	"context"

	"homework/internal/model"
)

func (c *Controller) List(ctx context.Context) ([]model.PickPoint, error) {
	dtos, err := c.storage.List(ctx)
	if err != nil {
		return nil, err
	}
	return model.LoadPickPointsFromDTO(dtos), nil
}
