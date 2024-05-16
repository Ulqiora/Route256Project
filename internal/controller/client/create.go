package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Create(ctx context.Context, obj model.Client) (string, error) {
	uuidObj, err := c.storage.Create(ctx, obj.MapToDTO())
	value, err := uuidObj.Value()
	return value.(string), err
}
