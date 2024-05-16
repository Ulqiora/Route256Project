package order_changers

import "github.com/Ulqiora/Route256Project/internal/model"

type ChangerOrderTape struct {
}

func (c ChangerOrderTape) Change(dto model.OrderInitData) (model.OrderInitData, error) {
	dto.Penny += 100
	return dto, nil
}
