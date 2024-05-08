package order

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"homework/internal/model"
)

func (c *ControllerOrder) ReturnOrder(ctx context.Context, orderID uint64, customerID uint64) error {
	err := c.tm.Run(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}, func(ctxTX context.Context) error {
		orderdto, err := c.storage.GetByID(ctxTX, orderID)
		if err != nil {
			return fmt.Errorf("order with id = %d not exists: %s", orderID, err)
		}
		if orderdto.CustomerID != int64(customerID) {
			return fmt.Errorf("this order belongs to another customer")
		}
		if (orderdto.State == model.EReadyToIssued && time.Now().Before(orderdto.ShelfLife.Time)) || orderdto.State == model.EReceived {
			err = c.storage.UpdateToReturned(ctxTX, orderID)
			if err != nil {
				return fmt.Errorf("order was not updated: %s", err)
			}
		} else {
			return fmt.Errorf("order can't be updated to returned: %s", err)
		}
		return nil
	})
	return err
}
