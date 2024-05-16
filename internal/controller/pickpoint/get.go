package pickpoint

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/jackc/pgtype"
)

func (c *Controller) GetByID(ctx context.Context, id string) (model.PickPoint, error) {
	var obj model.PickPoint
	var pickpointUuid pgtype.UUID
	err := pickpointUuid.Scan(id)
	if err != nil {
		return model.PickPoint{}, err
	}
	dto, err := c.storage.GetByID(ctx, pickpointUuid)
	if err != nil {
		return model.PickPoint{}, err
	}
	return obj.LoadFromDTO(dto), nil
}
