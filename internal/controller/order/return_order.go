package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/jackc/pgtype"
)

// ReturnOrder TODO: change transaction options
func (c *ControllerOrder) ReturnOrder(ctx context.Context, orderID string, customerID string) error {
	// pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}
	err := c.tm.Do(ctx, func(ctxTX context.Context) error {
		uuidOrder := pgtype.UUID{}
		err := uuidOrder.Set(orderID)
		if err != nil {
			return fmt.Errorf("incorrect format uuid: %s", err)
		}
		order, err := c.storage.GetByID(ctxTX, uuidOrder)
		if err != nil {
			return fmt.Errorf("order with id = %d not exists: %s", orderID, err)
		}
		if string(order.CustomerID.Bytes[:]) != string(customerID) {
			return fmt.Errorf("this order belongs to another customer")
		}
		if (order.State == model.EReadyToIssued && time.Now().Before(order.ShelfLife.Time)) || order.State == model.EReceived {
			err = c.storage.UpdateToReturned(ctxTX, uuidOrder)
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
