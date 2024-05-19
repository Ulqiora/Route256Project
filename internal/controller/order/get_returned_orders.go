package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/controller"
	"github.com/Ulqiora/Route256Project/internal/model"
)

func (c *ControllerOrder) GetReturnedOrders(ctx context.Context, values controller.ValuesView) ([]model.Order, error) {
	var filter paginationFilter
	err := filter.loadFromRequest(values)
	if err != nil {
		return nil, err
	}
	listDTO, err := c.storage.List(ctx)
	if err != nil {
		return nil, err
	}
	orders, err := model.LoadOrdersFromDTO(listDTO)
	list, err := filterByPageLimit(
		filterByStateOrder(orders, model.EReturned),
		filter,
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func filterByPageLimit(orders []model.Order, filter paginationFilter) ([]model.Order, error) {
	if filter.Page*filter.Limit > uint64(len(orders)) {
		return nil, errors.New("incorrect parameters")
	}

	if uint64(len(orders)) >= ((filter.Page + 1) * filter.Limit) {
		return orders[filter.Page*filter.Limit : ((filter.Page + 1) * filter.Limit)], nil
	} else {
		return orders[filter.Page*filter.Limit:], nil
	}
}

func filterByStateOrder(orders []model.Order, state string) []model.Order {
	ordersFiltered := make([]model.Order, 0)
	for _, orderObj := range orders {
		if orderObj.State == state {
			ordersFiltered = append(ordersFiltered, orderObj)
		}
	}
	fmt.Println(len(ordersFiltered))
	return ordersFiltered
}
