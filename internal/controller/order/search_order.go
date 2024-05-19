package order

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/controller"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

func (c *ControllerOrder) SearchOrders(ctx context.Context, customerID string, values controller.ValuesView) ([]model.Order, error) {
	uuidCustomer := pgtype.UUID{}
	err := uuidCustomer.Set(customerID)
	if err != nil {
		return nil, fmt.Errorf("incorrect format uuid: %s", err)
	}
	orderDTOs, err := c.storage.GetByCustomerID(ctx, uuidCustomer)
	if err != nil {
		return nil, errors.Wrap(err, "Error searching orders")
	}
	orders, err := model.LoadOrdersFromDTO(orderDTOs)
	if err != nil {
		return nil, fmt.Errorf("error loading orders: %s", err.Error())
	}
	if values.Has("last_n") && values.Get("last_n") != "0" {
		n, err := strconv.ParseUint(values.Get("last_n"), 10, 64)
		if err != nil {
			return nil, err
		}
		orders = filterByLastN(orders, n)
	} else if values.Has("pickpoint_id") && values.Get("pickpoint_id") != "0" {
		id := values.Get("pickpoint_id")
		orders = filterByPickPointID(orders, id)
	} else {
		return nil, errors.New("you need to set one of the additional parameters")
	}
	return orders, nil
}

func filterByLastN(inputOrders []model.Order, N uint64) []model.Order {
	outputOrders := make([]model.Order, len(inputOrders))
	copy(outputOrders, inputOrders)
	sort.Slice(outputOrders, func(i, j int) bool {
		return outputOrders[i].TimeCreated.Time().After(*outputOrders[j].TimeCreated.Time())
	})
	if N >= uint64(len(outputOrders)) {
		return outputOrders
	}
	return outputOrders[:N]
}

func filterByPickPointID(inputOrders []model.Order, pickpointID string) []model.Order {
	var outputOrders []model.Order
	for i := range inputOrders {
		if inputOrders[i].PickPointID == pickpointID {
			outputOrders = append(outputOrders, inputOrders[i])
		}
	}
	return outputOrders
}
