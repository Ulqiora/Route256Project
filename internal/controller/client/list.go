package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *Controller) List(ctx context.Context) ([]model.Client, error) {
	dtos, err := c.storage.List(ctx)
	if err != nil {
		return nil, err
	}
	clients, err := model.LoadClientsFromDTO(dtos)
	return clients, err
}
