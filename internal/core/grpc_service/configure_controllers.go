package grpc_service

import (
	"github.com/Ulqiora/Route256Project/internal/controller/client"
	"github.com/Ulqiora/Route256Project/internal/controller/order"
	"github.com/Ulqiora/Route256Project/internal/controller/pickpoint"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/model/order_changers"
	trmpgx "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Controllers struct {
	orderCtrl     *order.ControllerOrder
	pickpointCtrl *pickpoint.Controller
	clientCtrl    *client.Controller
}

func ConfigureControllers(repos Repositories, txManager *trmpgx.Manager) Controllers {
	changerMap := make(map[model.TypePacking]order_changers.ChangerOrder)
	changerMap[model.TypeBox] = &order_changers.ChangerOrderBox{}
	changerMap[model.TypeTape] = &order_changers.ChangerOrderTape{}
	changerMap[model.TypePackage] = &order_changers.ChangerOrderPackage{}
	return Controllers{
		orderCtrl:     order.New(repos.Order, changerMap, txManager),
		pickpointCtrl: pickpoint.New(repos.Pickpoint, txManager),
		clientCtrl:    client.New(repos.Client, txManager),
	}
}
