package client

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/jackc/pgtype"
)

func (c *Controller) Get(ctx context.Context, ObjId string) (model.Client, error) {
	var result model.Client
	var id pgtype.UUID
	err := id.Set(ObjId)
	if err != nil {
		return result, fmt.Errorf("controller: failed to set UUID: %w", err)
	}
	clientDTO, err := c.storage.GetByID(ctx, id)
	if err != nil {
		return result, fmt.Errorf("controller: failed to get by ID: %w", err)
	}
	result.LoadFromDTO(&clientDTO)
	return result, nil
}
