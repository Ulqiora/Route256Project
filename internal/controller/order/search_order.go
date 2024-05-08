package order

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"homework/internal/controller"
	"homework/internal/model"
)

func (c *ControllerOrder) SearchOrders(ctx context.Context, customerID uint64, values controller.ValuesView) ([]model.Order, error) {
	orderDTOs, err := c.storage.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, errors.Wrap(err, "Error searching orders")
	}
	orders := model.LoadOrdersFromDTO(orderDTOs)
	if values.Has("last_n") && values.Get("last_n") != "0" {
		n, err := strconv.ParseUint(values.Get("last_n"), 10, 64)
		if err != nil {
			return nil, err
		}
		orders = filterByLastN(orders, n)
	} else if values.Has("pickpoint_id") && values.Get("pickpoint_id") != "0" {
		id, err := strconv.ParseUint(values.Get("pickpoint_id"), 10, 64)
		if err != nil {
			return nil, err
		}
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
		fmt.Println(outputOrders)
		return outputOrders
	}
	fmt.Printf("last_n = %d", N)
	fmt.Println(outputOrders)
	return outputOrders[:N]
}

func filterByPickPointID(inputOrders []model.Order, pickpointID uint64) []model.Order {
	var outputOrders []model.Order
	for i := range inputOrders {
		if inputOrders[i].PickPointID == int64(pickpointID) {
			outputOrders = append(outputOrders, inputOrders[i])
		}
	}
	fmt.Printf("last_n = %d", pickpointID)
	return outputOrders
}
