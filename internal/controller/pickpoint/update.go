package pickpoint

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Update(ctx context.Context, object model.PickPoint) (string, error) {
	id, err := c.storage.Update(ctx, object.MapToDTO())
	if err != nil {
		return "", err
	}
	return string(id.Bytes[:]), nil
}
