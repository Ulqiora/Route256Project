package order_changers

import "homework/internal/model"

type ChangerOrder interface {
	Change(dto model.OrderInitData) (model.OrderInitData, error)
}
