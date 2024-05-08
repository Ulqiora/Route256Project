package order_changers

import "homework/internal/model"

type ChangerOrderTape struct {
}

func (c ChangerOrderTape) Change(dto model.OrderInitData) (model.OrderInitData, error) {
	dto.Penny += 100
	return dto, nil
}
