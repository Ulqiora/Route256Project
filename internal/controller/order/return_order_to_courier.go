package order

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"homework/internal/model"
	"homework/internal/repository"
)

const timeReturn time.Duration = time.Duration(time.Hour * 48)

func (c *ControllerOrder) ReturnOrderToCourier(ctx context.Context, orderID uint64) error {
	err := c.tm.Run(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}, func(ctxTX context.Context) error {
		orderdto, err := c.storage.GetByID(ctxTX, orderID)
		if err != nil {
			return fmt.Errorf("order with id = %d not exists: %s", orderID, err)
		}
		if orderdto.State != model.EReceived {
			return fmt.Errorf("order was been not updated to return to courier: %s", err)
		}
		if time.Now().After(orderdto.DateReceipt.Time.Add(timeReturn)) {
			return fmt.Errorf("%s: %s", repository.ErrorTimeOutForReturnOrder, err)
		}
		err = c.storage.UpdateToReturned(ctxTX, orderID)
		if err != nil {
			return fmt.Errorf("order was not updated to return to courier: %s", err)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Error return order to order_courier")
	}
	return nil
}
