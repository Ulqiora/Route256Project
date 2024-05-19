package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) Create(ctx context.Context, obj model.Client) (string, error) {
	dto, err := obj.MapToDTO()
	if err != nil {
		return "", err
	}
	uuidObj, err := c.storage.Create(ctx, dto)
	if err != nil {
		return "", err
	}
	value, err := uuidObj.Value()
	if err != nil {
		return "", err
	}
	return value.(string), err
}
