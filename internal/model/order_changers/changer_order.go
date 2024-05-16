package order_changers

import "github.com/Ulqiora/Route256Project/internal/model"

type ChangerOrder interface {
	Change(dto model.OrderInitData) (model.OrderInitData, error)
}
