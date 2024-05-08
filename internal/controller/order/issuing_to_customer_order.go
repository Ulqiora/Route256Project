package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"homework/internal/repository"
)

func (c *ControllerOrder) IssuingToCustomer(ctx context.Context, orderIDs []uint64) error {

	err := c.tm.Run(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}, func(ctxTX context.Context) error {
		orders, err := c.storage.ListReadyToIssued(ctxTX)
		if err != nil {
			return fmt.Errorf("error get list ready for issued: %s", err)
		}
		if isCorrectIDsForIssuing(orders, orderIDs) {
			for _, orderID := range orderIDs {
				err = c.storage.UpdateToReceived(ctxTX, orderID)
				if err != nil {
					return fmt.Errorf("error update to received orders: %s", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error issuing customer orders: %s", err)
	}
	return nil
}

func isCorrectIDsForIssuing(readyToIssuedOrders []repository.OrderDTO, idOrders []uint64) bool {
	var ordersMap = make(map[uint64]repository.OrderDTO)
	for _, orderObj := range readyToIssuedOrders {
		ordersMap[orderObj.ID] = orderObj
	}

	ordersFiltered := make([]repository.OrderDTO, 0, 8)
	for _, id := range idOrders {
		if _, ok := ordersMap[id]; !ok {
			return false
		} else {
			ordersFiltered = append(ordersFiltered, ordersMap[id])
		}
	}

	idCustomer := ordersFiltered[0].CustomerID
	for i := range ordersFiltered {
		if ordersFiltered[i].CustomerID != idCustomer {
			return false
		}
	}
	return true
}
