package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Update(ctx context.Context, obj model.Client) (string, error) {
	uuid, err := c.storage.Update(ctx, obj.MapToDTO())
	if err != nil {
		return "", err
	}
	return string(uuid.Bytes[:]), nil
}
