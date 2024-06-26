package order

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/pkg/errors"
)

func (c *ControllerOrder) AcceptOrder(ctx context.Context, data model.OrderInitData) (string, error) {
	if _, ok := c.changers[data.Type]; !ok {
		return "", errors.New("Incorrect order type")
	}
	data, err := c.changers[data.Type].Change(data)
	if err != nil {
		return "", err
	}
	var orderID string
	err = c.tm.Do(ctx, func(ctx context.Context) error {
		dto, err := data.MapToDTO()
		if err != nil {
			return fmt.Errorf("error load dto: %s", err.Error())
		}
		orderUUID, err := c.storage.Create(ctx, dto)
		if err != nil {
			return fmt.Errorf("error create order: %s", err.Error())
		}
		orderID = string(orderUUID.Bytes[:])
		return nil
	})
	return orderID, err
}
