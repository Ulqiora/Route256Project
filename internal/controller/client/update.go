package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Update(ctx context.Context, obj model.Client) (string, error) {
	dto, err := obj.MapToDTO()
	uuid, err := c.storage.Update(ctx, dto)
	if err != nil {
		return "", err
	}
	value, err := uuid.Value()
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
