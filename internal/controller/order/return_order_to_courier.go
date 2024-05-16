package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

const timeReturn time.Duration = time.Duration(time.Hour * 48)

func (c *ControllerOrder) ReturnOrderToCourier(ctx context.Context, orderID string) error {
	//, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}
	err := c.tm.Do(ctx, func(ctxTX context.Context) error {
		uuidOrder := pgtype.UUID{}
		err := uuidOrder.Set(orderID)
		if err != nil {
			return fmt.Errorf("incorrect format uuid: %s", err)
		}
		orderdto, err := c.storage.GetByID(ctxTX, uuidOrder)
		if err != nil {
			return fmt.Errorf("order with id = %s not exists: %s", orderID, err)
		}
		if orderdto.State != model.EReceived {
			return fmt.Errorf("order was been not updated to return to courier: %s", err)
		}
		if time.Now().After(orderdto.DateReceipt.Time.Add(timeReturn)) {
			return fmt.Errorf("%s: %s", repository.ErrorTimeOutForReturnOrder, err)
		}
		err = c.storage.UpdateToReturned(ctxTX, uuidOrder)
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
