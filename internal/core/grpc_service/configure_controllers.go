package grpc_service

import (
	"homework/internal/controller/order"
	"homework/internal/controller/pickpoint"
	"homework/internal/database/transaction_manager"
	"homework/internal/model"
	"homework/internal/model/order_changers"
)

type Controllers struct {
	orderCtrl     *order.ControllerOrder
	pickpointCtrl *pickpoint.Controller
}

func ConfigureControllers(repos Repositories, txManager *transaction_manager.TransactionManager) Controllers {
	changerMap := make(map[model.TypePacking]order_changers.ChangerOrder)
	changerMap[model.TypeBox] = &order_changers.ChangerOrderBox{}
	changerMap[model.TypeTape] = &order_changers.ChangerOrderTape{}
	changerMap[model.TypePackage] = &order_changers.ChangerOrderPackage{}
	return Controllers{
		orderCtrl:     order.New(repos.Order, changerMap, txManager),
		pickpointCtrl: pickpoint.New(repos.Pickpoint, txManager),
	}
}
