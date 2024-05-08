package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"homework/internal/model"
)

func (c *ControllerOrder) AcceptOrder(ctx context.Context, data model.OrderInitData) (uint64, error) {
	if _, ok := c.changers[data.Type]; !ok {
		return 0, errors.New("Incorrect order type")
	}
	data, err := c.changers[data.Type].Change(data)
	if err != nil {
		return 0, err
	}
	var orderID uint64
	err = c.tm.Run(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite}, func(ctxTX context.Context) error {
		fmt.Println("Start TRANSACTION")
		orderID, err = c.storage.Create(ctxTX, data.MapToDTO())
		if err != nil {
			return fmt.Errorf("error create order: %s", err)
		}
		return nil
	})
	return orderID, err
}
