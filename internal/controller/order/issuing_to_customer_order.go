package order

import (
	"context"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/jackc/pgtype"
)

// IssuingToCustomer TODO: change transaction params!
func (c *ControllerOrder) IssuingToCustomer(ctx context.Context, orderIDs []string) error {
	//pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite}
	err := c.tm.Do(ctx, func(ctxTX context.Context) error {
		orders, err := c.storage.ListReadyToIssued(ctxTX)
		if err != nil {
			return fmt.Errorf("error get list ready for issued: %s", err)
		}
		if isCorrectIDsForIssuing(orders, orderIDs) {
			for _, orderID := range orderIDs {
				uuidOrder := pgtype.UUID{}
				err = uuidOrder.Set(orderID)
				if err != nil {
					return fmt.Errorf("incorrect format uuid: %s", err)
				}
				err = c.storage.UpdateToReceived(ctxTX, uuidOrder)
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

func isCorrectIDsForIssuing(readyToIssuedOrders []repository.OrderDTO, idOrders []string) bool {
	var ordersMap = make(map[string]repository.OrderDTO)
	for _, orderObj := range readyToIssuedOrders {
		ordersMap[string(orderObj.ID.Bytes[:])] = orderObj
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
