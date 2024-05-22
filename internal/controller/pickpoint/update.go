package pickpoint

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Update(ctx context.Context, object model.PickPoint) (string, error) {
	dto, err := object.MapToDTO()
	if err != nil {
		return "", err
	}
	id, err := c.storage.Update(ctx, dto)
	if err != nil {
		return "", err
	}
	value, err := id.Value()
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
